package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {

	checks := map[string]bool{
		"SA5000": true,
		"SA6000": true,
		"SA9004": true,
	}

	var mychecks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	mychecks = append(mychecks, asmdecl.Analyzer)
	mychecks = append(mychecks, composite.Analyzer)
	mychecks = append(mychecks, ExitCheckAnalyzer)

	multichecker.Main(
		mychecks...,
	)
}
