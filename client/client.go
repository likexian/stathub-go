/*
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 *
 * Copyright 2015, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "crypto/tls"
    "os"
    "time"
    "strings"
    "github.com/likexian/host-stat-go"
    "github.com/likexian/simplejson-go"
)

const (
    CONFIG_FILE = "/client.json"
)

var (
    CLIENT_WORKDIR = ""
    CLIENT_START   = int64(0)
)

type Config struct {
    Id     string `json:"id"`
    Name   string `json:"name"`
    Server string `json:"server"`
    Key    string `json:"key"`
}

type Stat struct {
    Id        string  `json:"id"`
    TimeStamp int64   `json:"time_stamp"`
    HostName  string  `json:"host_name"`
    OSRelease string  `json:"os_release"`
    CPUName   string  `json:"cpu_name"`
    CPUCore   uint64  `json:"cpu_core"`
    Uptime    uint64  `json:"uptime"`
    Load      string  `json:"load"`
    CPURate   float64 `json:"cpu_rate"`
    MemRate   float64 `json:"mem_rate"`
    SwapRate  float64 `json:"swap_rate"`
    DiskRate  float64 `json:"disk_rate"`
    DiskWarn  string  `json:"disk_warn"`
    DiskRead  uint64  `json:"disk_read"`
    DiskWrite uint64  `json:"disk_write"`
    NetRead   uint64  `json:"net_read"`
    NetWrite  uint64  `json:"net_write"`
}

func Version() string {
    return "0.13.2"
}

func Author() string {
    return "[Li Kexian](https://www.likexian.com/)"
}

func License() string {
    return "Apache License, Version 2.0"
}

func main() {
    if (len(os.Args) > 1) {
        if (os.Args[1] == "-v" || os.Args[1] == "--version") {
            version := fmt.Sprintf("StatHub Client v%s\n%s\n%s", Version(), License(), Author())
            fmt.Println(version)
            os.Exit(0)
        }
    }

    pwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    CLIENT_WORKDIR := pwd

    CLIENT_START = time.Now().Unix()
    if !FileExists(CLIENT_WORKDIR + CONFIG_FILE) {
        SettingConfig(CLIENT_START)
    }

start:
    config, err := simplejson.Load(CLIENT_WORKDIR + CONFIG_FILE)
    if err != nil {
        return
    }

    id, _ := config.Get("id").String()
    name, _ := config.Get("name").String()
    server, _ := config.Get("server").String()
    key, _ := config.Get("key").String()

    stat := GetStat(id, name, CLIENT_START)
    surl := server + "/api/stat"
    skey := PassWord(key, stat)

    request, err := http.NewRequest("POST", surl, bytes.NewBuffer([]byte(stat)))
    request.Header.Set("X-Client-Key", skey)
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("User-Agent", "Stat Hub API Client/0.1.0 (i@likexian.com)")

    client := &http.Client{}
    tr := &http.Transport{
        // If not self-signed certificate please disabled this.
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client.Transport = tr

    response, err := client.Do(request)
    if err != nil {
        // auto select http or https
        if strings.Contains(err.Error(), "malformed HTTP response") {
            server = "https://" + strings.Split(server, "://")[1]
            WriteConfig(id, name, server, key)
            goto start
        } else if strings.Contains(err.Error(), "oversized record received with length") {
            server = "http://" + strings.Split(server, "://")[1]
            WriteConfig(id, name, server, key)
            goto start
        } else {
            fmt.Println(err)
            os.Exit(1)
        }
    }
    defer response.Body.Close()

    data, _ := ioutil.ReadAll(response.Body)
    text := string(data)
    if text != "" {
        fmt.Println(text)
        os.Exit(1)
    }
}

func SettingConfig(time_stamp int64) {
    host_info, _ := host_stat.GetHostInfo()
    host_name := host_info.HostName

    name := RawInput(fmt.Sprintf("> Please enter the NAME of THIS node [%s]:", host_name), true)
    if name == "" {
        name = host_name
    }

    server := RawInput("> Please enter the URL of SERVER :", false)
    key := RawInput("> Please enter the KEY of SERVER :", false)

    if len(server) <= 7 {
        server = "http://" + server
    }

    if server[:7] != "http://" && server[:8] != "https://" {
        server = "https://" + server
    }

    if server[len(server)-1:] == "/" {
        server = server[:len(server)-1]
    }

    random := fmt.Sprintf("%s%s", os.Getpid(), time_stamp)
    id := PassWord(key, random)

    WriteConfig(id, name, server, key)
}

func WriteConfig(id, name, server, key string) {
    config := Config{}
    config.Id = id
    config.Name = name
    config.Server = server
    config.Key = key

    data := simplejson.Json{}
    data.Data = config
    simplejson.Dump(CLIENT_WORKDIR + CONFIG_FILE, &data)
}

func GetStat(id string, name string, time_stamp int64) string {
    stat := Stat{}
    stat.Id = id
    stat.TimeStamp = time_stamp

    host_info, _ := host_stat.GetHostInfo()
    stat.OSRelease = host_info.Release + " " + host_info.OSBit
    if name == "" {
        stat.HostName = host_info.HostName
    } else {
        stat.HostName = name
    }

    cpu_info, _ := host_stat.GetCPUInfo()
    stat.CPUName = cpu_info.ModelName
    stat.CPUCore = cpu_info.CoreCount

    cpu_stat, _ := host_stat.GetCPUStat()
    stat.CPURate = Round(100-cpu_stat.IdleRate, 2)

    mem_stat, _ := host_stat.GetMemStat()
    stat.MemRate = mem_stat.MemRate
    stat.SwapRate = mem_stat.SwapRate

    disk_stat, _ := host_stat.GetDiskStat()
    disk_total := uint64(0)
    disk_used := uint64(0)
    for _, v := range disk_stat {
        disk_total += v.Total
        disk_used += v.Used
        if v.UsedRate > 90 {
            stat.DiskWarn += fmt.Sprintf("%s %.2f;", v.Mount, v.UsedRate)
        }
    }
    stat.DiskRate = Round(float64(disk_used)/float64(disk_total), 2)

    io_stat, _ := host_stat.GetIOStat()
    disk_read := uint64(0)
    disk_write := uint64(0)
    for _, v := range io_stat {
        disk_read += v.ReadBytes
        disk_write += v.WriteBytes
    }
    stat.DiskRead = disk_read
    stat.DiskWrite = disk_write

    net_stat, _ := host_stat.GetNetStat()
    net_write := uint64(0)
    net_read := uint64(0)
    for _, v := range net_stat {
        if v.Device != "lo" {
            net_write += v.TXBytes
            net_read += v.RXBytes
        }
    }
    stat.NetWrite = net_write
    stat.NetRead = net_read

    uptime_stat, _ := host_stat.GetUptimeStat()
    stat.Uptime = uint64(uptime_stat.Uptime)

    load_stat, _ := host_stat.GetLoadStat()
    stat.Load = fmt.Sprintf("%.2f %.2f %.2f", load_stat.LoadNow, load_stat.LoadPre, load_stat.LoadFar)

    json := simplejson.Json{}
    json.Data = stat
    result, _ := simplejson.Dumps(&json)

    return result
}
