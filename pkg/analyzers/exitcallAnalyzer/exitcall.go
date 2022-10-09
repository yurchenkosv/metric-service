package exitcall

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const Doc = ` check that no os.Exit call made from main package`

var Analyzer = &analysis.Analyzer{
	Name: "noExit",
	Doc:  Doc,
	Run:  run,
}

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

func findExitFunctionCall(pass *analysis.Pass, x *ast.CallExpr) {
	switch exp := x.Fun.(type) {
	case *ast.SelectorExpr:
		ident := exp.X.(*ast.Ident)
		if ident.Name == "os" && exp.Sel.Name == "Exit" {
			pass.Reportf(ident.Pos(), "os.Exit in main package")
		}
	}
}
