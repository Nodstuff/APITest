package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"net/http"
)

type TestData struct {
	TestName             string
	Endpoint             ObjectPath
	RequestMethod        RequestMethod
	Payload              map[string]string
	ErrorExpected        bool
	ExpectedResponseCode int
	AssertionFunc        func(resp *http.Response, data TestData) error
}

type TestResult struct {
	TestName             string
	ResponseCode         int
	ExpectedResponseCode int
	Response             string
	Fail                 bool
}

const (
	colorReset = "\033[0m"

	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[34m"
)

func (tr *TestResult) Print() {
	var textColor = colorGreen
	if tr.Fail {
		textColor = colorRed
	}
	fmt.Println(textColor, fmt.Sprintf("--------------------------%s-------------------------------", tr.TestName))
	fmt.Println(textColor, fmt.Sprintf("ResponseCode: %d", tr.ResponseCode))
	fmt.Println(textColor, fmt.Sprintf("ExpectedResponseCode: %d", tr.ExpectedResponseCode))
	if tr.Fail {
		fmt.Println(textColor, fmt.Sprintf("Response: %s", tr.Response))
	}
	fmt.Println(textColor, fmt.Sprintf("Test Fail: %v", tr.Fail), colorBlue)
}

func RunTest(data TestData) TestResult {
	res := TestResult{
		TestName:             data.TestName,
		ExpectedResponseCode: data.ExpectedResponseCode,
	}

	jsonBytes, err := json.Marshal(data.Payload)
	if err != nil {
		res.Response = err.Error()
		res.Fail = true
		return res
	}
	req, err := http.NewRequest(
		string(data.RequestMethod),
		fmt.Sprintf(
			"http://localhost:3000%s%s",
			APIPrefixPathAndVersioningForV1,
			string(data.Endpoint)),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		res.Response = err.Error()
		res.Fail = true
		return res
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		res.Response = err.Error()
		res.Fail = true
		return res
	}

	if data.AssertionFunc != nil {
		errs := data.AssertionFunc(resp, data)
		if errs.(*multierror.Error).Errors != nil {
			res.Fail = true
			res.Response = errs.Error()
			return res
		}
	}

	defer resp.Body.Close()

	return checkResponseCode(resp, data, res)
}

func checkResponseCode(resp *http.Response, data TestData, res TestResult) TestResult {
	res.ResponseCode = resp.StatusCode
	res.ExpectedResponseCode = data.ExpectedResponseCode

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		res.Response = err.Error()
		res.Fail = true
		return res
	}
	res.Response = string(b)

	if resp.StatusCode == data.ExpectedResponseCode && !data.ErrorExpected {
		return res
	} else {
		res.Fail = true
		return res
	}
}
