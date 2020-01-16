package config

import (
	"encoding/json"
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

type Configuration struct {
	Address         string
	StaticPath      string          `json:"static"`
	Server          *http.Server    `json:"server"`
	Cert            *Ssl            `json:"cert"`
	GracefulTimeout time.Duration   `json:"close_time"`
	BasicTimeout    time.Duration   `json:"basic_timeout"`
	FileTimeoutRate int64           `json:"file_timeout_rate"` //nanoseconds per 1 KB
	MaxFileTimeout  time.Duration   `json:"max_file_timeout"`
	MinFileTimeout  time.Duration   `json:"min_file_timeout"`
	DatabaseType    string          `json:"db_type"`
	Database        json.RawMessage `json:"db"`
	DatabaseReset   bool            `json:"db_clear"`
	Tika            Tika            `json:"tika"`
	FileLimit       int64           `json:"filelimit"`
	FreeSpace       int             `json:"total_free_space"`
	AdminKey        string
	GuestUser       *Guest
	SetupTimeout    time.Duration
}
