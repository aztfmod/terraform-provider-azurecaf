// Package readmeconsistency validates that the README resource status table
// and the documented resource-type counts agree with resourceDefinition.json.
package readmeconsistency

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Status markers used in the README resource status table.
const (
	StatusImplemented      = "✔"
	StatusNotImplemented   = "❌"
	StatusWillNotImplement = "⚠"
)

// IssueKind classifies a discrepancy.
type IssueKind string

// Discrepancy kinds reported by Check and CheckCounts.
const (
	KindMissingRow     IssueKind = "missing_row"
	KindNotMarked      IssueKind = "not_marked_implemented"
	KindNotDefined     IssueKind = "marked_but_not_defined"
	KindDuplicateRow   IssueKind = "duplicate_row"
	KindInvalidStatus  IssueKind = "invalid_status"
	KindStaleExclusion IssueKind = "stale_exclusion"
	KindCountMismatch  IssueKind = "count_mismatch"
)

// ExclusionsFile is the source file holding the exclusion list, used in reports.
const ExclusionsFile = "internal/readmeconsistency/exclusions.go"

// Issue describes a single discrepancy.
type Issue struct {
	Kind     IssueKind
	File     string
	Line     int
	Resource string
	Detail   string
	Fixable  bool
}

func (i Issue) String() string {
	loc := i.File
	if i.Line > 0 {
		loc = fmt.Sprintf("%s:%d", i.File, i.Line)
	}
	subject := i.Resource
	if subject == "" {
		subject = "-"
	}
	return fmt.Sprintf("[%s] %s %s: %s", i.Kind, loc, subject, i.Detail)
}

// Row is a resource row in the README status table.
type Row struct {
	Name   string
	Status string
	// Line is the 1-based line number in the README.
	Line int
}

// Table is the parsed README resource status table.
type Table struct {
	Rows []Row
	// HeaderLine and LastLine are 0-based indexes into the README lines.
	HeaderLine int
	LastLine   int
}

var (
	tableHeaderPattern    = regexp.MustCompile(`^\|\s*resource\s*\|\s*status\s*\|$`)
	tableSeparatorPattern = regexp.MustCompile(`^\|\s*:?-+:?\s*\|\s*:?-+:?\s*\|$`)
	tableRowPattern       = regexp.MustCompile(`^\|\s*([^|]*?)\s*\|\s*([^|]*?)\s*\|$`)
	countPattern          = regexp.MustCompile(`\b(\d+)(?:\*\*)?\s+(?:Azure\s+|supported\s+)?resource types`)
)

// ParseDefinitions returns the resource names declared in resourceDefinition.json.
func ParseDefinitions(data []byte) ([]string, error) {
	var defs []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, fmt.Errorf("parse resource definitions: %w", err)
	}
	names := make([]string, 0, len(defs))
	for i, d := range defs {
		if strings.TrimSpace(d.Name) == "" {
			return nil, fmt.Errorf("resource definition at index %d has no name", i)
		}
		names = append(names, d.Name)
	}
	return names, nil
}

// normalizeStatus strips the emoji variation selector so "⚠️" equals "⚠".
func normalizeStatus(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\uFE0F", ""))
}

func validStatus(s string) bool {
	switch s {
	case StatusImplemented, StatusNotImplemented, StatusWillNotImplement:
		return true
	}
	return false
}

func trimLine(l string) string {
	return strings.TrimSpace(strings.TrimRight(l, "\r"))
}

// ParseStatusTable extracts the resource status table from README content.
func ParseStatusTable(readme string) (Table, error) {
	lines := strings.Split(readme, "\n")
	header := -1
	for i, l := range lines {
		if tableHeaderPattern.MatchString(trimLine(l)) {
			header = i
			break
		}
	}
	if header < 0 {
		return Table{}, fmt.Errorf("resource status table header '|resource | status |' not found")
	}
	if header+1 >= len(lines) || !tableSeparatorPattern.MatchString(trimLine(lines[header+1])) {
		return Table{}, fmt.Errorf("resource status table separator row missing after line %d", header+1)
	}

	t := Table{HeaderLine: header, LastLine: header + 1}
	for i := header + 2; i < len(lines); i++ {
		l := trimLine(lines[i])
		if !strings.HasPrefix(l, "|") {
			break
		}
		m := tableRowPattern.FindStringSubmatch(l)
		if m == nil {
			return Table{}, fmt.Errorf("malformed resource status row at line %d: %q", i+1, l)
		}
		t.Rows = append(t.Rows, Row{Name: m[1], Status: normalizeStatus(m[2]), Line: i + 1})
		t.LastLine = i
	}
	return t, nil
}

// Check compares the definitions with the README status table.
// excluded lists definition names (aliases/helpers) intentionally absent from the table.
func Check(file string, definitions []string, table Table, excluded map[string]string) []Issue {
	defined := toSet(definitions)

	var issues []Issue
	seen := make(map[string]Row, len(table.Rows))
	for _, r := range table.Rows {
		if prev, dup := seen[r.Name]; dup {
			issues = append(issues, Issue{Kind: KindDuplicateRow, File: file, Line: r.Line, Resource: r.Name,
				Detail: fmt.Sprintf("row duplicates line %d", prev.Line)})
			continue
		}
		seen[r.Name] = r

		_, isExcluded := excluded[r.Name]
		switch {
		case !validStatus(r.Status):
			issues = append(issues, Issue{Kind: KindInvalidStatus, File: file, Line: r.Line, Resource: r.Name,
				Detail: fmt.Sprintf("status %q must be one of %s, %s, %s", r.Status,
					StatusImplemented, StatusNotImplemented, StatusWillNotImplement)})
		case r.Status == StatusImplemented && !defined[r.Name]:
			issues = append(issues, Issue{Kind: KindNotDefined, File: file, Line: r.Line, Resource: r.Name,
				Detail: "marked ✔ in README but missing from resourceDefinition.json"})
		case r.Status != StatusImplemented && defined[r.Name] && !isExcluded:
			issues = append(issues, Issue{Kind: KindNotMarked, File: file, Line: r.Line, Resource: r.Name,
				Detail: fmt.Sprintf("defined in resourceDefinition.json but marked %s in README", r.Status), Fixable: true})
		}
	}

	for _, n := range sortedCopy(definitions) {
		if _, skip := excluded[n]; skip {
			continue
		}
		if _, ok := seen[n]; !ok {
			issues = append(issues, Issue{Kind: KindMissingRow, File: file, Resource: n,
				Detail: "defined in resourceDefinition.json but missing from the README status table", Fixable: true})
		}
	}

	for _, n := range sortedKeys(excluded) {
		if !defined[n] {
			issues = append(issues, Issue{Kind: KindStaleExclusion, File: ExclusionsFile, Resource: n,
				Detail: "excluded from the README table but no longer defined; remove it from the exclusion list"})
		}
	}
	return issues
}

// CheckCounts reports "<N> ... resource types" phrases that disagree with want.
func CheckCounts(file, content string, want int) []Issue {
	var issues []Issue
	for i, l := range strings.Split(content, "\n") {
		for _, m := range countPattern.FindAllStringSubmatch(l, -1) {
			got, err := strconv.Atoi(m[1])
			if err != nil || got == want {
				continue
			}
			issues = append(issues, Issue{Kind: KindCountMismatch, File: file, Line: i + 1,
				Detail: fmt.Sprintf("states %d resource types, resourceDefinition.json defines %d", got, want), Fixable: true})
		}
	}
	return issues
}

// FixCounts rewrites "<N> ... resource types" phrases to want.
func FixCounts(content string, want int) string {
	return countPattern.ReplaceAllStringFunc(content, func(s string) string {
		m := countPattern.FindStringSubmatchIndex(s)
		return s[:m[2]] + strconv.Itoa(want) + s[m[3]:]
	})
}

// FixTable marks defined resources as implemented and inserts missing rows in
// alphabetical position. Issues needing human judgement are left untouched.
func FixTable(readme string, definitions []string, excluded map[string]string) (string, error) {
	table, err := ParseStatusTable(readme)
	if err != nil {
		return "", err
	}
	lines := strings.Split(readme, "\n")
	defined := toSet(definitions)

	present := make(map[string]bool, len(table.Rows))
	for _, r := range table.Rows {
		present[r.Name] = true
		_, isExcluded := excluded[r.Name]
		if defined[r.Name] && !isExcluded && validStatus(r.Status) && r.Status != StatusImplemented {
			lines[r.Line-1] = formatRow(r.Name, StatusImplemented)
		}
	}

	rows := append([]string(nil), lines[table.HeaderLine+2:table.LastLine+1]...)
	added := 0
	for _, n := range sortedCopy(definitions) {
		if _, skip := excluded[n]; skip || present[n] {
			continue
		}
		pos := len(rows)
		for i, l := range rows {
			if m := tableRowPattern.FindStringSubmatch(trimLine(l)); m != nil && m[1] > n {
				pos = i
				break
			}
		}
		rows = append(rows[:pos], append([]string{formatRow(n, StatusImplemented)}, rows[pos:]...)...)
		present[n] = true
		added++
	}

	out := make([]string, 0, len(lines)+added)
	out = append(out, lines[:table.HeaderLine+2]...)
	out = append(out, rows...)
	out = append(out, lines[table.LastLine+1:]...)
	return strings.Join(out, "\n"), nil
}

func formatRow(name, status string) string {
	return fmt.Sprintf("|%s | %s |", name, status)
}

func toSet(in []string) map[string]bool {
	out := make(map[string]bool, len(in))
	for _, n := range in {
		out[n] = true
	}
	return out
}

func sortedCopy(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
