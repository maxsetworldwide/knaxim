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

package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// Ssl is SSL configuration values
type Ssl struct {
	CertFile string `json:"cert" yaml:"cert"`
	KeyFile  string `json:"key" yaml:"key"`
	HTTPport string `json:"http_port" yaml:"http_port"`
}

// Guest user setup configuration
type Guest struct {
	Name  string
	Pass  string
	Email string
}

// Tika connection and configuration values
type Tika struct {
	Type        string `json:"type" yaml:"type"`
	Path        string `json:"path" yaml:"path"`
	Port        string `json:"port" yaml:"port"`
	MaxFiles    int    `json:"child_max_files" yaml:"child_max_files"`
	TaskPulse   int    `json:"child_task_pulse" yaml:"child_task_pulse"`
	TaskTimeout int    `json:"child_task_timeout" yaml:"child_task_timeout"`
	PingPulse   int    `json:"child_ping_pulse" yaml:"child_ping_pulse"`
	PingTimeout int    `json:"child_ping_timeout" yaml:"child_ping_timeout"`
}

// SMTP configuration values
type SMTP struct {
	From       string `json:"from" yaml:"from"`
	Server     string `json:"server" yaml:"server"`
	Credential struct {
		Identity string
		Username string
		Password string
		Host     string
	} `json:"cred" yaml:"cred"`
}

// Duration type that has custom UnmarshalJSON to allow use of
// time.ParseDuration for more easily written and read duration values
// in configuration
type Duration struct {
	time.Duration
}

// UnmarshalJSON allows use of time.ParseDuration if associated value is // a string
func (d *Duration) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.New("no data when parsing data")
	}
	var durstr string
	if err := json.Unmarshal(b, &durstr); err == nil {
		d.Duration, err = time.ParseDuration(durstr)
		return err
	}
	return json.Unmarshal(b, &(d.Duration))
}

// UnmarshalYAML Decodes Yaml value into a time.Duration
func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	return n.Decode(&(d.Duration))
}

// Configuration struct that is populated by the Configuration file
type Configuration struct {
	Address              string
	StaticPath           string       `json:"static" yaml:"static"`
	IndexPath            string       `json:"index" yaml:"index"`
	Server               *http.Server `json:"server" yaml:"server"`
	Cert                 *Ssl         `json:"cert" yaml:"cert"`
	GracefulTimeout      Duration     `json:"close_time" yaml:"close_time"`
	BasicTimeout         Duration     `json:"basic_timeout" yaml:"basic_timeout"`
	FileTimeoutRate      int64        `json:"file_timeout_rate" yaml:"file_timeout_rate"` //nanoseconds per 1 KB
	MaxFileTimeout       Duration     `json:"max_file_timeout" yaml:"max_file_timeout"`
	MinFileTimeout       Duration     `json:"min_file_timeout" yaml:"min_file_timeout"`
	ActiveFileProcessing int
	DatabaseType         string `json:"db_type" yaml:"db_type"`
	Database             Raw    `json:"db" yaml:"db"`
	DatabaseReset        bool   `json:"db_clear" yaml:"db_clear"`
	Tika                 Tika   `json:"tika" yaml:"tika"`
	GotenPath            string `json:"gotenpath" yaml:"gotenpath"`
	FileLimit            int64  `json:"filelimit" yaml:"filelimit"`
	FreeSpace            int    `json:"total_free_space" yaml:"total_free_space"`
	MaxFileCount         int64  `json:"maxfilecount" yaml:"maxfilecount"`
	AdminKey             string
	GuestUser            *Guest
	SetupTimeout         Duration
	UserTimeouts         struct {
		Inactivity Duration
		Total      Duration
	}
	Email       SMTP
	ErrorEmail  string `json:"error_email" yaml:"error_email"`
	LogPath     string `json:"log_path" yaml:"log_path"`
	PrivateMode bool
}

// MarshalJSON extracts data fields from configuration to generate json configuration
func (c Configuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Address": c.Address,
		"static":  c.StaticPath,
		"index":   c.IndexPath,
		"Server": map[string]interface{}{
			"Addr":              c.Server.Addr,
			"ReadHeaderTimeout": c.Server.ReadHeaderTimeout,
		},
		"cert":                 c.Cert,
		"close_time":           c.GracefulTimeout,
		"basic_timeout":        c.BasicTimeout,
		"file_timeout_rate":    c.FileTimeoutRate,
		"max_file_timeout":     c.MaxFileTimeout,
		"min_file_timeout":     c.MinFileTimeout,
		"ActiveFileProcessing": c.ActiveFileProcessing,
		"db_type":              c.DatabaseType,
		"db":                   c.Database,
		"db_clear":             c.DatabaseReset,
		"tika":                 c.Tika,
		"gotenpath":            c.GotenPath,
		"filelimit":            c.FileLimit,
		"total_free_space":     c.FreeSpace,
		"maxfilecount":         c.MaxFileCount,
		"AdminKey":             c.AdminKey,
		"GuestUser":            c.GuestUser,
		"SetupTimeout":         c.SetupTimeout,
		"UserTimeouts":         c.UserTimeouts,
		"Email":                c.Email,
		"error_email":          c.ErrorEmail,
		"log_path":             c.LogPath,
		"PrivateMode":          c.PrivateMode,
	})
}

// Raw represents a value not to be decoded. Primarily for data fields that can hold a variety of data types
type Raw struct {
	JSON json.RawMessage
	YAML *yaml.Node
}

// MarshalJSON output the contents of the JSON field
func (r *Raw) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.JSON)
}

// UnmarshalJSON saves the byte slice in the JSON field
func (r *Raw) UnmarshalJSON(b []byte) error {
	r.JSON = json.RawMessage(b)
	return nil
}

// MarshalYAML outputs the content of the YAML field
func (r *Raw) MarshalYAML() (interface{}, error) {
	return r.YAML, nil
}

// UnmarshalYAML saves the YAML node in the YAML field
func (r *Raw) UnmarshalYAML(n *yaml.Node) error {
	r.YAML = n
	return nil
}
