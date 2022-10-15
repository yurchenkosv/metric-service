package exitcall

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const Doc = `check that no os.Exit call made from main package`

var Analyzer = &analysis.Analyzer{
	Name: "noExit", // Analyzer name
	Doc:  Doc,      // Documentation of analyzer
	Run:  run,      // Main function that analyzer executes
}

// run is main function analyzer executes
// This function searches for os.Exit in main package
func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		if pass.Pkg.Name() == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.CallExpr:
					findExitFunctionCall(pass, x)
				}
				return true
			})
		}
	}
	return nil, nil
}

// findExitFunctionCall is helper function that searches for Exit call in ast tree
func findExitFunctionCall(pass *analysis.Pass, x *ast.CallExpr) {
	switch exp := x.Fun.(type) {
	case *ast.SelectorExpr:
		if exp.Sel.Name == "Exit" {
			pass.Reportf(exp.Sel.Pos(), "os.Exit in main package")
		}
	}
}
