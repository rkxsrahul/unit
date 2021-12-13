package health

import (
	"testing"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	ot "github.com/opentracing/opentracing-go"
)

func TestServiceHealth(t *testing.T) {
	// test case when no configuration is enabled
	err := ServiceHealth(ot.StartSpan("test service health -> test case1"))
	if err != nil {
		t.Log("test case passed")
	} else {
		t.Fatal("expected some error but got ", err)
	}

	//test case enable cockroach and redis configuration and rest of service closed
	config.TomlFile = "/home/nitin/work/src/git.xenonstack.com/stacklabs/stacklabs-auth/testing.toml"
	config.SetConfig()
	err = ServiceHealth(ot.StartSpan("test service health -> test case2"))
	if err != nil {
		t.Log("test case passed")
	} else {
		t.Fatal("expected some error but got ", err)
	}

	//test case enable cockroach, redis configuration, deployment service and rest of service closed
	config.TomlFile = "/home/nitin/work/src/git.xenonstack.com/stacklabs/stacklabs-auth/testing1.toml"
	config.SetConfig()
	err = ServiceHealth(ot.StartSpan("test service health -> test case3"))
	if err != nil {
		t.Log("test case passed")
	} else {
		t.Fatal("expected some error but got ", err)
	}

	//test case when all components are working
	config.TomlFile = "/home/nitin/work/src/git.xenonstack.com/stacklabs/stacklabs-auth/test.toml"
	config.SetConfig()
	err = ServiceHealth(ot.StartSpan("test service health -> test case4"))
	if err == nil {
		t.Log("test case passed")
	} else {
		t.Fatal("expected no error but got ", err)
	}
}
