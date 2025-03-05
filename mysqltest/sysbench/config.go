// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysbench

import (
	"strconv"

	"github.com/cybergarage/go-sqltest/sqltest/sysbench"
)

const (
	// https://github.com/akopytov/sysbench
	MySQLHost     = "mysql-host"
	MySQLPort     = "mysql-port"
	MySQLUser     = "mysql-user"
	MySQLPassword = "mysql-password"
	MySQLDB       = "mysql-db"
	MySQLSSL      = "mysql-ssl"
	MySQLDebug    = "mysql-debug"
)

// Config represents a sysbench config.
type Config struct {
	sysbench.Config
}

// NewDefaultConfig returns a new default config.
func NewDefaultConfig() *Config {
	cfg := &Config{
		Config: sysbench.NewDefaultConfig(),
	}
	cfg.SetHost("127.0.0.1")
	cfg.SetPort(3306)
	cfg.SetSSL(true)
	cfg.SetDebug(true)
	return cfg
}

func (cfg *Config) setBool(key string, value bool) {
	if value {
		cfg.Config.Set(key, "on")
		return
	}
	cfg.Config.Set(key, "off")
}

// SetHost sets the host.
func (cfg *Config) SetHost(host string) {
	cfg.Config.Set(MySQLHost, host)
}

// SetPort sets the port.
func (cfg *Config) SetPort(port int) {
	cfg.Config.Set(MySQLPort, strconv.Itoa(port))
}

// SetUser sets the user.
func (cfg *Config) SetUser(user string) {
	cfg.Config.Set(MySQLUser, user)
}

// SetPassword sets the password.
func (cfg *Config) SetPassword(password string) {
	cfg.Config.Set(MySQLPassword, password)
}

// SetDB sets the db.
func (cfg *Config) SetDB(db string) {
	cfg.Config.Set(MySQLDB, db)
}

// SetSSL sets the ssl.
func (cfg *Config) SetSSL(ssl bool) {
	cfg.setBool(MySQLSSL, ssl)
}

// SetDebug sets the debug.
func (cfg *Config) SetDebug(debug bool) {
	cfg.setBool(MySQLDebug, debug)
}
