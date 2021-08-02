package main

import (
	run "github.com/iwita/simple-sw/pkg/runtime"
)

func main() {
	r := run.NewRuntime("testfiles/applicationrequest.yaml", run.WithInputFile("testfiles/applicant.json"))
	r.Start()
}
