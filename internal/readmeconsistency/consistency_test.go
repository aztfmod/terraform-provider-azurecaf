package readmeconsistency

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const sampleReadme = `# Title

Supports **3 Azure resource types**.

|resource | status |
|---|---|
|azurerm_a | ✔ |
|azurerm_c | ❌ |
|azurerm_e | ⚠️ |

## Next
`

func kinds(issues []Issue) []string {
	var out []string
	for _, i := range issues {
		out = append(out, string(i.Kind)+":"+i.Resource)
	}
	sort.Strings(out)
	return out
}

func TestParseDefinitions(t *testing.T) {
	names, err := ParseDefinitions([]byte(`[{"name":"a"},{"name":"b","slug":"x"}]`))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(names, []string{"a", "b"}) {
		t.Fatalf("got %v", names)
	}
	if _, err := ParseDefinitions([]byte(`{`)); err == nil {
		t.Fatal("expected JSON error")
	}
	if _, err := ParseDefinitions([]byte(`[{"name":" "}]`)); err == nil {
		t.Fatal("expected empty-name error")
	}
}

func TestParseStatusTable(t *testing.T) {
	table, err := ParseStatusTable(sampleReadme)
	if err != nil {
		t.Fatal(err)
	}
	want := []Row{
		{Name: "azurerm_a", Status: StatusImplemented, Line: 7},
		{Name: "azurerm_c", Status: StatusNotImplemented, Line: 8},
		{Name: "azurerm_e", Status: StatusWillNotImplement, Line: 9},
	}
	if !reflect.DeepEqual(table.Rows, want) {
		t.Fatalf("rows = %#v", table.Rows)
	}
	if table.HeaderLine != 4 || table.LastLine != 8 {
		t.Fatalf("header/last = %d/%d", table.HeaderLine, table.LastLine)
	}
}

func TestParseStatusTableCRLF(t *testing.T) {
	table, err := ParseStatusTable(strings.ReplaceAll(sampleReadme, "\n", "\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Rows) != 3 || table.Rows[0].Status != StatusImplemented {
		t.Fatalf("rows = %#v", table.Rows)
	}
}

func TestParseStatusTableErrors(t *testing.T) {
	cases := map[string]string{
		"no header":     "# nothing\n",
		"no separator":  "|resource | status |\n|azurerm_a | ✔ |\n",
		"header at EOF": "|resource | status |",
		"malformed row": "|resource | status |\n|---|---|\n|azurerm_a | ✔ | extra |\n",
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseStatusTable(in); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestCheckConsistent(t *testing.T) {
	table, _ := ParseStatusTable(sampleReadme)
	if issues := Check("README.md", []string{"azurerm_a", "alias"}, table, map[string]string{"alias": "x"}); len(issues) != 0 {
		t.Fatalf("unexpected issues: %v", issues)
	}
}

func TestCheckDetectsAllKinds(t *testing.T) {
	readme := `|resource | status |
|---|---|
|azurerm_a | ✔ |
|azurerm_b | ❌ |
|azurerm_c | ✔ |
|azurerm_c | ✔ |
|azurerm_d | maybe |
|azurerm_alias | ❌ |
`
	table, err := ParseStatusTable(readme)
	if err != nil {
		t.Fatal(err)
	}
	defs := []string{"azurerm_a", "azurerm_b", "azurerm_new", "azurerm_alias"}
	excluded := map[string]string{"azurerm_alias": "alias", "azurerm_gone": "stale"}
	got := kinds(Check("README.md", defs, table, excluded))
	want := []string{
		"duplicate_row:azurerm_c",
		"invalid_status:azurerm_d",
		"marked_but_not_defined:azurerm_c",
		"missing_row:azurerm_new",
		"not_marked_implemented:azurerm_b",
		"stale_exclusion:azurerm_gone",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %v\nwant %v", got, want)
	}
}

func TestCheckCounts(t *testing.T) {
	content := "Supports **3 Azure resource types**.\nAll 4 resource types\n5 supported resource types\nresource types without number\n"
	issues := CheckCounts("README.md", content, 3)
	if len(issues) != 2 || issues[0].Line != 2 || issues[1].Line != 3 {
		t.Fatalf("issues = %v", issues)
	}
	if !issues[0].Fixable || issues[0].Kind != KindCountMismatch {
		t.Fatalf("issue = %#v", issues[0])
	}
	fixed := FixCounts(content, 3)
	if len(CheckCounts("README.md", fixed, 3)) != 0 {
		t.Fatalf("FixCounts left mismatches:\n%s", fixed)
	}
	if !strings.Contains(fixed, "**3 Azure resource types**") || !strings.Contains(fixed, "3 supported resource types") {
		t.Fatalf("FixCounts damaged formatting:\n%s", fixed)
	}
}

func TestFixTable(t *testing.T) {
	defs := []string{"azurerm_a", "azurerm_c", "azurerm_b", "azurerm_z", "alias"}
	fixed, err := FixTable(sampleReadme, defs, map[string]string{"alias": "x"})
	if err != nil {
		t.Fatal(err)
	}
	table, err := ParseStatusTable(fixed)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, r := range table.Rows {
		names = append(names, r.Name+"="+r.Status)
	}
	want := []string{"azurerm_a=✔", "azurerm_b=✔", "azurerm_c=✔", "azurerm_e=⚠", "azurerm_z=✔"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("got %v", names)
	}
	if !strings.HasSuffix(fixed, "\n## Next\n") {
		t.Fatalf("content after table was altered:\n%s", fixed)
	}
	if issues := Check("README.md", defs, table, map[string]string{"alias": "x"}); len(issues) != 0 {
		t.Fatalf("issues after fix: %v", issues)
	}
}

func TestFixTableError(t *testing.T) {
	if _, err := FixTable("no table", nil, nil); err == nil {
		t.Fatal("expected error")
	}
}

func writeRepo(t *testing.T, defs, readme string) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{DefinitionsFile: defs, ReadmeFile: readme}
	for _, f := range CountFiles[1:] {
		files[f] = "Supports 3 resource types\n"
	}
	for name, content := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestRunAndFix(t *testing.T) {
	saved := Excluded
	Excluded = map[string]string{}
	t.Cleanup(func() { Excluded = saved })

	defs := `[{"name":"azurerm_a"},{"name":"azurerm_c"},{"name":"azurerm_d"},{"name":"azurerm_f"}]`
	root := writeRepo(t, defs, sampleReadme)

	res, err := Run(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if res.Definitions != 4 || res.Rows != 3 || len(res.Issues) != 7 {
		t.Fatalf("result = %+v", res)
	}
	report := FormatReport(res)
	if !strings.Contains(report, "FAIL: 7 discrepancies") || !strings.Contains(report, "make fix_readme_consistency") {
		t.Fatalf("report:\n%s", report)
	}

	res, err = Run(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues) != 0 || len(res.Fixed) != len(CountFiles) {
		t.Fatalf("after fix: %+v", res)
	}
	if report := FormatReport(res); !strings.Contains(report, "OK:") || !strings.Contains(report, "fixed: README.md") {
		t.Fatalf("report:\n%s", report)
	}

	res, err = Run(root, true)
	if err != nil || len(res.Fixed) != 0 {
		t.Fatalf("fix should be idempotent: %+v %v", res, err)
	}
}

func TestRunReportsManualIssues(t *testing.T) {
	saved := Excluded
	Excluded = map[string]string{}
	t.Cleanup(func() { Excluded = saved })

	readme := strings.Replace(sampleReadme, "**3 Azure", "**1 Azure", 1)
	root := writeRepo(t, `[{"name":"azurerm_x"}]`, readme)
	for _, f := range CountFiles[1:] {
		if err := os.WriteFile(filepath.Join(root, f), []byte("1 resource types\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	res, err := Run(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if got := kinds(res.Issues); !reflect.DeepEqual(got, []string{"marked_but_not_defined:azurerm_a"}) {
		t.Fatalf("got %v", got)
	}
	if report := FormatReport(res); !strings.Contains(report, "1 need a manual fix") || strings.Contains(report, "make fix_readme_consistency") {
		t.Fatalf("report:\n%s", report)
	}
}

func TestRunErrors(t *testing.T) {
	if _, err := Run(t.TempDir(), false); err == nil {
		t.Fatal("expected missing definitions error")
	}
	if _, err := Run(writeRepo(t, `{`, sampleReadme), false); err == nil {
		t.Fatal("expected JSON error")
	}
	if _, err := Run(writeRepo(t, `[]`, "no table"), false); err == nil {
		t.Fatal("expected table error")
	}
	if _, err := Run(writeRepo(t, `[]`, "no table"), true); err == nil {
		t.Fatal("expected fix table error")
	}
	root := writeRepo(t, `[]`, sampleReadme)
	if err := os.Remove(filepath.Join(root, CountFiles[1])); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root, false); err == nil {
		t.Fatal("expected missing count file error")
	}
	if _, err := Run(root, true); err == nil {
		t.Fatal("expected missing count file error during fix")
	}
}

func TestIssueString(t *testing.T) {
	if got := (Issue{Kind: KindMissingRow, File: "README.md", Resource: "r", Detail: "d"}).String(); got != "[missing_row] README.md r: d" {
		t.Fatalf("got %q", got)
	}
	if got := (Issue{Kind: KindCountMismatch, File: "README.md", Line: 3, Detail: "d"}).String(); got != "[count_mismatch] README.md:3 -: d" {
		t.Fatalf("got %q", got)
	}
}

// TestRepositoryIsConsistent enforces consistency of this repository's README
// and docs with resourceDefinition.json as part of `go test ./...`.
func TestRepositoryIsConsistent(t *testing.T) {
	res, err := Run(filepath.Join("..", ".."), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues) > 0 {
		t.Fatal(FormatReport(res))
	}
}
