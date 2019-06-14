/*
 * Copyright 2015-2019 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 */

package main

import (
	"errors"
	"strings"

	"github.com/likexian/simplejson-go"
)

// Config storing server and client config
type Config struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	PassWord   string `json:"password"`
	ServerKey  string `json:"server_key"`
	ServerUrl  string `json:"server_url"`
	DaemonUser string `json:"daemon_user"`
	PidFile    string `json:"pid_file"`
	LogFile    string `json:"log_file"`
	BaseDir    string `json:"base_dir"`
	DataDir    string `json:"data_dir"`
	TLSCert    string `json:"tls_cert"`
	TLSKey     string `json:"tls_key"`
	File       string `json:"file"`
}

// SaveConfig save the config to file
func SaveConfig(config Config) (err error) {
	data := simplejson.New(config)
	return data.Dump(config.File)
}

// GetConfig returns the config of file
func GetConfig(fname string) (config Config, err error) {
	config = Config{}

	data, err := simplejson.Load(fname)
	if err != nil {
		return
	}

	config.Id = data.Get("id").MustString("")
	config.Name = data.Get("name").MustString("")
	config.Role = data.Get("role").MustString("")
	config.PassWord = data.Get("password").MustString("")
	config.ServerKey = data.Get("server_key").MustString("")
	config.ServerUrl = data.Get("server_url").MustString("")
	config.DaemonUser = data.Get("daemon_user").MustString("")
	config.PidFile = data.Get("pid_file").MustString("")
	config.LogFile = data.Get("log_file").MustString("")
	config.BaseDir = data.Get("base_dir").MustString("")
	config.DataDir = data.Get("data_dir").MustString("")
	config.TLSCert = data.Get("tls_cert").MustString("")
	config.TLSKey = data.Get("tls_key").MustString("")
	config.File = fname

	if config.Id == "" {
		return config, errors.New("missing id config")
	}

	if config.Role == "" {
		return config, errors.New("missing role config")
	}

	if config.ServerKey == "" {
		return config, errors.New("missing server_key config")
	}

	if config.PidFile == "" {
		return config, errors.New("missing pid_file config")
	}

	if config.LogFile == "" {
		return config, errors.New("missing log_file config")
	}

	if config.BaseDir == "" {
		return config, errors.New("missing base_dir config")
	}

	if config.Role == "server" {
		if config.PassWord == "" {
			return config, errors.New("missing password config")
		}
		if config.DataDir == "" {
			return config, errors.New("missing data_dir config")
		}
		if config.TLSCert == "" {
			return config, errors.New("missing tls_cert config")
		}
		if config.TLSKey == "" {
			return config, errors.New("missing tls_key config")
		}
	} else {
		if config.ServerUrl == "" {
			return config, errors.New("missing server_url config")
		}
	}

	if !strings.HasSuffix(config.BaseDir, "/") {
		config.BaseDir += "/"
	}

	if strings.HasSuffix(config.ServerUrl, "/") {
		config.ServerUrl = config.ServerUrl[:len(config.ServerUrl)-1]
	}

	err = SaveConfig(config)

	return
}

// newServerConfig generate the new server config file
func newServerConfig(fname, id, name, passWord, serverKey string) (err error) {
	config := Config{
		id,
		name,
		"server",
		passWord,
		serverKey,
		DEFAULT_SERVER_URL,
		DEFAULT_PROCESS_USER,
		DEFAULT_PROCESS_LOCK,
		DEFAULT_PROCESS_LOG,
		DEFAULT_BASE_DIR,
		DEFAULT_DATA_DIR,
		DEFAULT_TLS_CERT,
		DEFAULT_TLS_KEY,
		fname,
	}

	return SaveConfig(config)
}

// newClientConfig generate the new cilent config file
func newClientConfig(fname, id, name, serverKey, serverUrl string) (err error) {
	config := Config{
		id,
		name,
		"client",
		"",
		serverKey,
		serverUrl,
		DEFAULT_PROCESS_USER,
		DEFAULT_PROCESS_LOCK,
		DEFAULT_PROCESS_LOG,
		DEFAULT_BASE_DIR,
		DEFAULT_DATA_DIR,
		DEFAULT_TLS_CERT,
		DEFAULT_TLS_KEY,
		fname,
	}

	return SaveConfig(config)
}
