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
    "fmt"
    "os"
    "time"
    "flag"
    "path/filepath"
    "github.com/likexian/daemon-go"
    "github.com/likexian/host-stat-go"
    "github.com/likexian/logger-go"
)


var (
    // SERVER_START is server start timestamp
    SERVER_START   = int64(0)
    // SERVER_CONFIG is server config data
    SERVER_CONFIG  = Config{}
    // SERVER_LOGGER is server logger
    SERVER_LOGGER  = logger.New(os.Stderr, logger.INFO)
)


func main() {
    SERVER_START = time.Now().Unix()

    if DEBUG {
        SERVER_LOGGER = logger.New(os.Stderr, logger.DEBUG)
    }

    show_version    := flag.Bool("v", false, "show current version")
    config_file     := flag.String("c", "", "set configuration file")
    init_server     := flag.Bool("init-server", false, "init server configuration")
    init_client     := flag.Bool("init-client", false, "init client configuration")
    server_key      := flag.String("server-key", "", "set server key, required when init client")
    server_url      := flag.String("server-url", "", "set server url, required when init client")

    flag.Parse()

    if *show_version {
        version := fmt.Sprintf("StatHub v%s-%s\n%s\n%s", TPL_VERSION, TPL_REVHEAD, TPL_LICENSE, TPL_AUTHOR)
        fmt.Println(version)
        os.Exit(0)
    }

    if *config_file == "" {
        flag.Usage()
        os.Exit(-1)
    }

    if *init_server {
        hostInfo, _ := hoststat.GetHostInfo()
        time_stamp := fmt.Sprintf("%d", SERVER_START)
        id := PassWord(fmt.Sprintf("%d", os.Getpid()), time_stamp)
        key := PassWord(id, time_stamp)
        password := PassWord(key, "likexian")
        err := newServerConfig(*config_file, id, hostInfo.HostName, password, key)
        if err != nil {
            SERVER_LOGGER.Critical(err.Error())
            os.Exit(-1)
        } else {
            SERVER_LOGGER.Info("init server configuration successful")
            os.Exit(0)
        }
    }

    if *init_client {
        if *server_key == "" {
            SERVER_LOGGER.Critical("server key is required, set it by --server-key.")
            os.Exit(-1)
        }
        if *server_url == "" {
            SERVER_LOGGER.Critical("server url is required, set it by --server-url.")
            os.Exit(-1)
        }
        hostInfo, _ := hoststat.GetHostInfo()
        time_stamp := fmt.Sprintf("%d", SERVER_START)
        id := PassWord(fmt.Sprintf("%d", os.Getpid()), time_stamp)
        err := newClientConfig(*config_file, id, hostInfo.HostName, *server_key, *server_url)
        if err != nil {
            SERVER_LOGGER.Critical(err.Error())
            os.Exit(-1)
        } else {
            SERVER_LOGGER.Info("init client configuration successful")
            os.Exit(0)
        }
    }

    if !FileExists(*config_file) {
        SERVER_LOGGER.Critical(fmt.Sprintf("configuration file %s is not found.\n", *config_file))
        os.Exit(-1)
    }

    var err error
    SERVER_CONFIG, err = GetConfig(*config_file)
    if err != nil {
        SERVER_LOGGER.Critical(fmt.Sprintf("configuration load failed, %s", err.Error()))
        os.Exit(-1)
    }

    if SERVER_CONFIG.Role == "server" {
        if !FileExists(SERVER_CONFIG.BaseDir + SERVER_CONFIG.DataDir) {
            err := os.MkdirAll(SERVER_CONFIG.BaseDir + SERVER_CONFIG.DataDir, 0755)
            if err != nil {
                SERVER_LOGGER.Critical(err.Error())
                os.Exit(-1)
            }
        }
        if !FileExists(SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSCert) {
            err := WriteFile(SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSCert, TPL_CERT["cert.pem"])
            if err != nil {
                SERVER_LOGGER.Critical(err.Error())
                os.Exit(-1)
            }
        }
        if !FileExists(SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSKey) {
            err := WriteFile(SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSKey, TPL_CERT["cert.key"])
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
        c := daemon.Config {
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
