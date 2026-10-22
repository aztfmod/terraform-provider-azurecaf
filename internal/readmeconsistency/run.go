package readmeconsistency

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Repository-relative paths inspected by Run.
const (
	DefinitionsFile = "resourceDefinition.json"
	ReadmeFile      = "README.md"
)

// CountFiles are documents whose "<N> resource types" phrases must match the
// number of definitions.
var CountFiles = []string{
	ReadmeFile,
	"docs/index.md",
	"docs/resources/azurecaf_name.md",
	"docs/data-sources/azurecaf_name.md",
}

// Result is the outcome of Run.
type Result struct {
	Definitions int
	Rows        int
	Issues      []Issue
	// Fixed lists files rewritten when fix is enabled.
	Fixed []string
}

// Run validates the repository at root. When fix is true it first applies
// automatic fixes, then reports any remaining issues.
func Run(root string, fix bool) (Result, error) {
	defData, err := os.ReadFile(filepath.Join(root, DefinitionsFile))
	if err != nil {
		return Result{}, err
	}
	defs, err := ParseDefinitions(defData)
	if err != nil {
		return Result{}, err
	}

	var res Result
	if fix {
		if res.Fixed, err = applyFixes(root, defs); err != nil {
			return Result{}, err
		}
	}

	readme, err := os.ReadFile(filepath.Join(root, ReadmeFile))
	if err != nil {
		return Result{}, err
	}
	table, err := ParseStatusTable(string(readme))
	if err != nil {
		return Result{}, fmt.Errorf("%s: %w", ReadmeFile, err)
	}
	res.Definitions = len(defs)
	res.Rows = len(table.Rows)
	res.Issues = Check(ReadmeFile, defs, table, Excluded)

	for _, f := range CountFiles {
		content, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			return Result{}, err
		}
		res.Issues = append(res.Issues, CheckCounts(f, string(content), len(defs))...)
	}
	return res, nil
}

func applyFixes(root string, defs []string) ([]string, error) {
	var fixed []string
	for _, f := range CountFiles {
		path := filepath.Join(root, f)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		updated := FixCounts(string(data), len(defs))
		if f == ReadmeFile {
			if updated, err = FixTable(updated, defs, Excluded); err != nil {
				return nil, fmt.Errorf("%s: %w", f, err)
			}
		}
		if updated == string(data) {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(updated), info.Mode().Perm()); err != nil {
			return nil, err
		}
		fixed = append(fixed, f)
	}
	return fixed, nil
}

// FormatReport renders a human-readable report.
func FormatReport(r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "README consistency: %d definitions (%d excluded), %d status table rows\n",
		r.Definitions, len(Excluded), r.Rows)
	for _, f := range r.Fixed {
		fmt.Fprintf(&b, "fixed: %s\n", f)
	}
	if len(r.Issues) == 0 {
		b.WriteString("OK: README.md and docs are consistent with resourceDefinition.json\n")
		return b.String()
	}
	fixable := 0
	fmt.Fprintf(&b, "FAIL: %d discrepancies found\n", len(r.Issues))
	for _, i := range r.Issues {
		b.WriteString("  " + i.String() + "\n")
		if i.Fixable {
			fixable++
		}
	}
	if fixable > 0 {
		fmt.Fprintf(&b, "%d can be fixed automatically with: make fix_readme_consistency\n", fixable)
	}
	if fixable < len(r.Issues) {
		fmt.Fprintf(&b, "%d need a manual fix (README.md, resourceDefinition.json, or %s)\n",
			len(r.Issues)-fixable, ExclusionsFile)
	}
	return b.String()
}
