// Command readme-consistency checks that README.md and the provider docs agree
// with resourceDefinition.json. Use -fix to apply automatic fixes.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aztfmod/terraform-provider-azurecaf/internal/readmeconsistency"
)

func main() {
	root := flag.String("root", ".", "repository root")
	fix := flag.Bool("fix", false, "apply automatic fixes before checking")
	flag.Parse()

	res, err := readmeconsistency.Run(*root, *fix)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	fmt.Print(readmeconsistency.FormatReport(res))
	if len(res.Issues) > 0 {
		os.Exit(1)
	}
}
