package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var testOk = `1
2
3
4
4
4
5`

var testOkResult = `1
2
3
4
5
`

func TestOk(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testOk))
	out := new(bytes.Buffer)
	err := uniq(in, out)
	if err != nil {
		t.Errorf("test for OK failed - error")
	}
	if out.String() != testOkResult {
		t.Errorf("test for OK failed - results not match")
	}
}

var testFail = `3
2`

func TestFail(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testFail))
	out := new(bytes.Buffer)
	err := uniq(in, out)
	if err == nil {
		t.Errorf("test for Fail failed - no error")
	}
}
