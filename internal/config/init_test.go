/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */
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
