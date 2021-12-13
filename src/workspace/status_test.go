package workspace

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

func TestStatus(t *testing.T) {
	span := opentracing.StartSpan("simple changepassword")
	status, _ := Status("xenonstack", span)
	if status != 202 {
		t.Error("test case fail")
	}
}
