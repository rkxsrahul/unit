package admin

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

// test to list workspace
func TestListWorkSpaces(t *testing.T) {
	span := opentracing.StartSpan("simple listworkspace")

	//passing span
	_, err := ListWorkSpaces(span)
	if err != nil {
		t.Error("test case fail")
	}

}
