package main

import (
	"APITest/internal"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/assert"
	"net/http"
	"os"
	"time"
)

func main() {
	now := time.Now()
	var failed bool
	var results []internal.TestResult

	tests := []internal.TestData{
		{
			TestName:             "Sample test",
			Endpoint:             internal.InfraWorkspace,
			RequestMethod:        internal.MethodGet,
			ErrorExpected:        false,
			ExpectedResponseCode: 200,
			Payload:              map[string]string{},
			AssertionFunc: func(resp *http.Response, data internal.TestData) error {
				var errs error
				errs = multierror.Append(errs, assert.So(resp.StatusCode, assertions.ShouldEqual, data.ExpectedResponseCode).Error())
				errs = multierror.Append(errs, assert.So("Sample test", assertions.ShouldEqual, data.TestName).Error())
				return errs
			},
		},
		{
			TestName:             "Sample Test 2",
			Endpoint:             internal.ServiceWorkspace,
			RequestMethod:        internal.MethodPost,
			ErrorExpected:        false,
			ExpectedResponseCode: 201,
			Payload:              map[string]string{"jim": "bah"},
		},
		{
			TestName:             "Sample test 3",
			Endpoint:             internal.Environment,
			RequestMethod:        internal.MethodPut,
			ErrorExpected:        false,
			ExpectedResponseCode: 202,
			Payload:              map[string]string{},
		},
		{
			TestName:             "Sample test 4",
			Endpoint:             internal.Proxy,
			RequestMethod:        internal.MethodDelete,
			ErrorExpected:        true,
			ExpectedResponseCode: 500,
			Payload:              map[string]string{},
		},
	}

	for _, t := range tests {
		results = append(results, internal.RunTest(t))
	}

	for _, r := range results {
		if r.Fail {
			failed = true
		}
		r.Print()
	}

	elapsed := time.Since(now)
	fmt.Printf("Tests took %s to run", elapsed)

	if failed {
		os.Exit(1)
	}
}
