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
    "time"
    "github.com/likexian/host-stat-go"
    "github.com/likexian/simplejson-go"
)


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


func GetStat(id string, name string) string {
    stat := Stat{}
    stat.Id = id
    stat.TimeStamp = time.Now().Unix()

    host_info, _ := hoststat.GetHostInfo()
    stat.OSRelease = host_info.Release + " " + host_info.OSBit
    if name == "" {
        stat.HostName = host_info.HostName
    } else {
        stat.HostName = name
    }

    cpu_info, _ := hoststat.GetCPUInfo()
    stat.CPUName = cpu_info.ModelName
    stat.CPUCore = cpu_info.CoreCount

    cpu_stat, _ := hoststat.GetCPUStat()
    stat.CPURate = Round(100 - cpu_stat.IdleRate, 2)

    mem_stat, _ := hoststat.GetMemStat()
    stat.MemRate = mem_stat.MemRate
    stat.SwapRate = mem_stat.SwapRate

    disk_stat, _ := hoststat.GetDiskStat()
    disk_total := uint64(0)
    disk_used := uint64(0)
    for _, v := range disk_stat {
        disk_total += v.Total
        disk_used += v.Used
        if v.UsedRate > 90 {
            stat.DiskWarn += fmt.Sprintf("%s %.2f;", v.Mount, v.UsedRate)
        }
    }
    stat.DiskRate = Round(float64(disk_used) / float64(disk_total), 2)

    io_stat, _ := hoststat.GetIOStat()
    disk_read := uint64(0)
    disk_write := uint64(0)
    for _, v := range io_stat {
        disk_read += v.ReadBytes
        disk_write += v.WriteBytes
    }
    stat.DiskRead = disk_read
    stat.DiskWrite = disk_write

    net_stat, _ := hoststat.GetNetStat()
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

    uptime_stat, _ := hoststat.GetUptimeStat()
    stat.Uptime = uint64(uptime_stat.Uptime)

    load_stat, _ := hoststat.GetLoadStat()
    stat.Load = fmt.Sprintf("%.2f %.2f %.2f", load_stat.LoadNow, load_stat.LoadPre, load_stat.LoadFar)

    data := simplejson.New(stat)
    result, _ := data.Dumps()

    return result
}
