package workspace

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

func TestDeleteWorkspace(t *testing.T) {

	span := opentracing.StartSpan("simple changepassword")
	DeleteWorkspace("xenonstack", token, span)

}

func TestDeleteDedicatedWorkspace(t *testing.T) {
	err := deleteDedicatedWorkspace("xenonstack", token)
	if err == nil {
		t.Error("test case fail")
	}
}
