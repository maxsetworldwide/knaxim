// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package email

import (
	"flag"
	"os"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/config"
)

var to = flag.String("to", "", "address to send test email")

func init() {
	config.V.Email.From = "noreply@maxset.org"
	config.V.Email.Server = "sub5.mail.dreamhost.com:587"
	config.V.Email.Credential.Identity = ""
	config.V.Email.Credential.Username = "ai@maxset.org"
	config.V.Email.Credential.Password = "anM5x3B8"
	config.V.Email.Credential.Host = "sub5.mail.dreamhost.com"
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func buildTo(t *string) []string {
	return strings.Split(*t, ",")
}

func TestSend(t *testing.T) {
	if len(*to) == 0 {
		t.Fatalf("Test error: please specify an email address for which to send a test email.")
	}
	t.Log("sending to: ", to)
	t.Run("Reset", func(t *testing.T) {
		err := SendResetEmail(buildTo(to), "id:2020-02-13-23:56:51:262t", "knaxim.com", "")
		if err != nil {
			t.Fatal("unable to send email: ", err)
		}
	})
	t.Run("Error", func(t *testing.T) {
		emailMsg := "this is a test email error message"
		config.V.ErrorEmail = *to
		err := SendErrorEmail(emailMsg)
		if err != nil {
			t.Fatal("unable to send email: ", err)
		}
	})
}
