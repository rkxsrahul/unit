package config

import (
	"log"
	"os"
	"testing"
)

func TestConfigurationWithEnv(t *testing.T) {
	// test case when no variable is set
	ConfigurationWithEnv()
	if Conf.Service.Port == "8000" {
		t.Log("test case passed")
	} else {
		t.Fatal("expected to get value for service port is 8000 but got ", Conf.Service.Port)
	}

	// test case when variable is set
	os.Setenv("AKIRA_AUTH_PORT", "9000")
	ConfigurationWithEnv()
	if Conf.Service.Port == "9000" {
		t.Log("test case passed")
	} else {
		t.Fatal("expected to get value for service port is 9000 but got ", Conf.Service.Port)
	}
}

func TestConfigurationWithToml(t *testing.T) {
	// test case when no file is passed
	err := ConfigurationWithToml("")
	if err != nil {
		t.Log("test case passed")
	} else {
		t.Fatal("expected some error but got ", err)
	}

	//test case when file is passed
	err = ConfigurationWithToml("/home/nitin/work/src/git.xenonstack.com/stacklabs/stacklabs-auth/testing.toml")
	log.Println(err)
	if err == nil {
		t.Log("test case passed")
	} else {
		t.Fatal("expected no error but got ", err)
	}

	//test case when file is passed with no port
	Conf.Service.Port = ""
	err = ConfigurationWithToml("/home/nitin/work/src/git.xenonstack.com/stacklabs/stacklabs-auth/testing1.toml")
	log.Println(err)
	if Conf.Service.Port == "8000" {
		t.Log("test case passed")
	} else {
		t.Fatal("expected to get value for service port is 8000 but got ", Conf.Service.Port)
	}
}

func TestSetConfig(t *testing.T) {
	// test case for configuration through environment
	TomlFile = ""
	SetConfig()
	// test case for configuration through toml
	TomlFile = "/home/nitin/work/src/git.xenonstack.com/stacklabs/stacklabs-auth/testing.toml"
	SetConfig()
}
