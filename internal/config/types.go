package config

import (
	"encoding/json"
	"net/http"
	"time"
)

type SslCert struct {
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
	HTTPport string `json:"http_port"`
}

type Configuration struct {
	Address         string
	StaticPath      string          `json:"static"`
	Server          *http.Server    `json:"server"`
	Cert            *SslCert        `json:"cert"`
	GracefulTimeout time.Duration   `json:"close_time"`
	BasicTimeout    time.Duration   `json:"basic_timeout"`
	FileTimeoutRate int64           `json:"file_timeout_rate"` //nanoseconds per 1 KB
	MaxFileTimeout  time.Duration   `json:"max_file_timeout"`
	MinFileTimeout  time.Duration   `json:"min_file_timeout"`
	DatabaseType    string          `json:"db_type"`
	Database        json.RawMessage `json:"db"`
	DatabaseReset   bool            `json:"db_clear"`
	Tika            TikaConf        `json:"tika"`
	FileLimit       int64           `json:"filelimit"`
	Smtp            struct {
		Active   bool
		Identity string
		Username string
		Password string
		Host     string
		Path     string
		From     string
	}
	FreeSpace int `json:"total_free_space"`
	AdminKey  string
	GuestUser *GuestConf
}

type GuestConf struct {
	Name  string
	Pass  string
	Email string
}

type TikaConf struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	Port        string `json:"port"`
	MaxFiles    int    `json:"child_max_files"`
	TaskPulse   int    `json:"child_task_pulse"`
	TaskTimeout int    `json:"child_task_timeout"`
	PingPulse   int    `json:"child_ping_pulse"`
	PingTimeout int    `json:"child_ping_timeout"`
}
