package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"

	exitcall "github.com/yurchenkosv/metric-service/pkg/analyzers/exitcallAnalyzer"
)

func main() {
	var checks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		}
	}

	for _, v := range quickfix.Analyzers {
		if v.Analyzer.Name == "QF1003" {
			checks = append(checks, v.Analyzer)
			break
		}
	}

	checks = append(checks, exitcall.Analyzer)
	checks = append(checks, asmdecl.Analyzer)
	checks = append(checks, assign.Analyzer)
	checks = append(checks, atomic.Analyzer)
	checks = append(checks, atomicalign.Analyzer)
	checks = append(checks, bools.Analyzer)
	checks = append(checks, buildssa.Analyzer)
	checks = append(checks, buildtag.Analyzer)
	checks = append(checks, cgocall.Analyzer)
	checks = append(checks, copylock.Analyzer)
	checks = append(checks, ctrlflow.Analyzer)
	checks = append(checks, deepequalerrors.Analyzer)
	checks = append(checks, fieldalignment.Analyzer)
	checks = append(checks, findcall.Analyzer)
	checks = append(checks, framepointer.Analyzer)
	checks = append(checks, ifaceassert.Analyzer)
	checks = append(checks, inspect.Analyzer)
	checks = append(checks, tests.Analyzer)
	checks = append(checks, loopclosure.Analyzer)
	checks = append(checks, lostcancel.Analyzer)
	checks = append(checks, printf.Analyzer)
	checks = append(checks, reflectvaluecompare.Analyzer)
	checks = append(checks, shadow.Analyzer)
	checks = append(checks, shift.Analyzer)
	checks = append(checks, sigchanyzer.Analyzer)
	checks = append(checks, sortslice.Analyzer)
	checks = append(checks, stdmethods.Analyzer)
	checks = append(checks, stringintconv.Analyzer)
	checks = append(checks, testinggoroutine.Analyzer)
	checks = append(checks, unmarshal.Analyzer)
	checks = append(checks, unreachable.Analyzer)
	checks = append(checks, unsafeptr.Analyzer)
	checks = append(checks, unusedresult.Analyzer)
	checks = append(checks, unusedwrite.Analyzer)

	multichecker.Main(checks...)
}
