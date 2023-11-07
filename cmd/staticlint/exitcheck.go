package main

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for no using Exit() func in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	fset := token.NewFileSet()
	for _, f := range pass.Files {
		if f.Name.Name == "main" {
			ast.Inspect(f, func(n ast.Node) bool {
				if c, ok := n.(*ast.CallExpr); ok {
					if s, ok := c.Fun.(*ast.SelectorExpr); ok {
						if s.Sel.Name == "Exit" {
							fmt.Printf("find os.Exit func in package main: %+v \n", s)
							printer.Fprint(os.Stdout, fset, f)
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
