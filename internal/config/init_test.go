package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadingConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testconfig*.json")
	if err != nil {
		t.Fatal("unable to make temporary file ", err.Error())
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(testconfigjson); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	if err := ParseConfig(tmpfile.Name()); err != nil {
		t.Log(tmpfile.Name())
		t.Fatal(err)
	}
}

func TestLoadingYAML(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testconfig*.yml")
	if err != nil {
		t.Fatal("unable to make temporary file ", err.Error())
	}
	//defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(testconfigyaml); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	if err := ParseConfig(tmpfile.Name()); err != nil {
		t.Log(tmpfile.Name())
		t.Fatal(err)
	}
}
