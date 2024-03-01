package main

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
)

// ErrNoExitAnalyzer check for os.Exit()
var ErrNoExitAnalyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "check for direct usage of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr: // os.Exit
				switch ce := x.Fun.(type) {
				case *ast.SelectorExpr:
					if ce.Sel.Name == "Exit" && fmt.Sprintf("%s", ce.X) == "os" {
						pass.Reportf(ce.Pos(), "os.Exit() it's bad")
					}
				}
			}
			return true
		})
	}

	return nil, nil
}

func main() {
	var mychecks = []*analysis.Analyzer{
		ErrNoExitAnalyzer,
	}

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	multichecker.Main(mychecks...)
}
