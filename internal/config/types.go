package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Ssl struct {
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
	HTTPport string `json:"http_port"`
}

type Guest struct {
	Name  string
	Pass  string
	Email string
}

type Tika struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	Port        string `json:"port"`
	MaxFiles    int    `json:"child_max_files"`
	TaskPulse   int    `json:"child_task_pulse"`
	TaskTimeout int    `json:"child_task_timeout"`
	PingPulse   int    `json:"child_ping_pulse"`
	PingTimeout int    `json:"child_ping_timeout"`
}

type SMTP struct {
	From       string `json:"from"`
	Server     string `json:"server"`
	Credential struct {
		Identity string
		Username string
		Password string
		Host     string
	} `json:"cred"`
}

type Duration struct {
	time.Duration
}

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

type Configuration struct {
	Address         string
	StaticPath      string          `json:"static"`
	IndexPath       string          `json:"index"`
	Server          *http.Server    `json:"server"`
	Cert            *Ssl            `json:"cert"`
	GracefulTimeout Duration        `json:"close_time"`
	BasicTimeout    Duration        `json:"basic_timeout"`
	FileTimeoutRate int64           `json:"file_timeout_rate"` //nanoseconds per 1 KB
	MaxFileTimeout  Duration        `json:"max_file_timeout"`
	MinFileTimeout  Duration        `json:"min_file_timeout"`
	DatabaseType    string          `json:"db_type"`
	Database        json.RawMessage `json:"db"`
	DatabaseReset   bool            `json:"db_clear"`
	Tika            Tika            `json:"tika"`
	GotenPath       string          `json:"gotenpath"`
	FileLimit       int64           `json:"filelimit"`
	FreeSpace       int             `json:"total_free_space"`
	AdminKey        string
	GuestUser       *Guest
	SetupTimeout    Duration
	UserTimeouts    struct {
		Inactivity Duration
		Total      Duration
	}
	Email SMTP
}
