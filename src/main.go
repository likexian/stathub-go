/*
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 *
 * Copyright 2015-2019, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/likexian/daemon-go"
	"github.com/likexian/logger-go"
	"os"
	"path/filepath"
	"time"
)

var (
	// SERVER_START is server start timestamp
	SERVER_START = int64(0)
	// SERVER_CONFIG is server config data
	SERVER_CONFIG = Config{}
	// SERVER_LOGGER is server logger
	SERVER_LOGGER = logger.New(os.Stderr, logger.INFO)
)

func main() {
	SERVER_START = time.Now().Unix()

	if DEBUG {
		SERVER_LOGGER = logger.New(os.Stderr, logger.DEBUG)
	}

	showVersion := flag.Bool("v", false, "show current version")
	configFile := flag.String("c", "", "set configuration file")
	initServer := flag.Bool("init-server", false, "init server configuration")
	initClient := flag.Bool("init-client", false, "init client configuration")
	serverKey := flag.String("server-key", "", "set server key, required when init client")
	serverUrl := flag.String("server-url", "", "set server url, required when init client")

	flag.Parse()

	if *showVersion {
		version := fmt.Sprintf("StatHub v%s-%s\n%s\n%s", Version(), TPL_REVHEAD, License(), Author())
		fmt.Println(version)
		os.Exit(0)
	}

	if *configFile == "" {
		flag.Usage()
		os.Exit(-1)
	}

	if *initServer {
		timeStamp := fmt.Sprintf("%d", SERVER_START)
		id := Md5(fmt.Sprintf("%d", os.Getpid()), timeStamp)
		key := Md5(id, timeStamp)
		password := Md5(key, "likexian")
		err := newServerConfig(*configFile, id, "", password, key)
		if err != nil {
			SERVER_LOGGER.Critical(err.Error())
			os.Exit(-1)
		} else {
			SERVER_LOGGER.Info("init server configuration successful")
			os.Exit(0)
		}
	}

	if *initClient {
		if *serverKey == "" {
			SERVER_LOGGER.Critical("server key is required, set it by --server-key.")
			os.Exit(-1)
		}
		if *serverUrl == "" {
			SERVER_LOGGER.Critical("server url is required, set it by --server-url.")
			os.Exit(-1)
		}
		timeStamp := fmt.Sprintf("%d", SERVER_START)
		id := Md5(fmt.Sprintf("%d", os.Getpid()), timeStamp)
		err := newClientConfig(*configFile, id, "", *serverKey, *serverUrl)
		if err != nil {
			SERVER_LOGGER.Critical(err.Error())
			os.Exit(-1)
		} else {
			SERVER_LOGGER.Info("init client configuration successful")
			os.Exit(0)
		}
	}

	if !FileExists(*configFile) {
		SERVER_LOGGER.Critical(fmt.Sprintf("configuration file %s is not found.\n", *configFile))
		os.Exit(-1)
	}

	var err error
	SERVER_CONFIG, err = GetConfig(*configFile)
	if err != nil {
		SERVER_LOGGER.Critical(fmt.Sprintf("configuration load failed, %s", err.Error()))
		os.Exit(-1)
	}

	if SERVER_CONFIG.Role == "server" {
		if !FileExists(SERVER_CONFIG.BaseDir + SERVER_CONFIG.DataDir) {
			err := os.MkdirAll(SERVER_CONFIG.BaseDir+SERVER_CONFIG.DataDir, 0755)
			if err != nil {
				SERVER_LOGGER.Critical(err.Error())
				os.Exit(-1)
			}
		}
		if !FileExists(SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSCert) {
			err := WriteFile(SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSCert, TPL_CERT["cert.pem"])
			if err != nil {
				SERVER_LOGGER.Critical(err.Error())
				os.Exit(-1)
			}
		}
		if !FileExists(SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSKey) {
			err := WriteFile(SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSKey, TPL_CERT["cert.key"])
			if err != nil {
				SERVER_LOGGER.Critical(err.Error())
				os.Exit(-1)
			}
		}
	}

	for _, v := range []string{SERVER_CONFIG.PidFile, SERVER_CONFIG.LogFile} {
		ds, _ := filepath.Split(SERVER_CONFIG.BaseDir + v)
		if ds != "" && !FileExists(ds) {
			err := os.MkdirAll(ds, 0755)
			if err != nil {
				SERVER_LOGGER.Critical(err.Error())
				os.Exit(-1)
			}
		}
	}

	if !DEBUG {
		c := daemon.Config{
			Pid:   SERVER_CONFIG.BaseDir + SERVER_CONFIG.PidFile,
			Log:   SERVER_CONFIG.BaseDir + SERVER_CONFIG.LogFile,
			User:  SERVER_CONFIG.DaemonUser,
			Chdir: "",
		}
		err := c.Daemon()
		if err != nil {
			SERVER_LOGGER.Critical(err.Error())
			os.Exit(-1)
		}
	}

	SERVER_LOGGER.Info("server start at %d", SERVER_START)
	if SERVER_CONFIG.Role == "server" {
		go StatService()
		HttpService()
	} else {
		StatService()
	}
}
