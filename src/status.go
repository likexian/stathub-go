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
    "os"
    "sort"
    "time"
    "errors"
    "io/ioutil"
    "github.com/likexian/simplejson-go"
)


type Status struct {
    Id         string
    IP         string
    Name       string
    Status     string
    Uptime     string
    Load       string
    NetRead    string
    NetWrite   string
    DiskRead   string
    DiskWrite  string
    DiskWarn   string
    CPURate    float64
    MemRate    float64
    SwapRate   float64
    DiskRate   float64
    NetTotal   string
    OSRelease  string
    LastUpdate string
}


type Statuses []Status


func (n Statuses) Len() int {return len(n)}
func (n Statuses) Swap(i, j int) {n[i], n[j] = n[j], n[i]}
func (n Statuses) Less(i, j int) bool {return n[i].Name < n[j].Name}


func ReadStatus(data_dir string) Statuses {
    data := []Status{}
    files, err := ioutil.ReadDir(data_dir)
    if err == nil {
        for _, f := range files {
            if FileExists(data_dir + "/" + f.Name() + "/status") {
                d, err := simplejson.Load(data_dir + "/" + f.Name() + "/status")
                if err != nil {
                    continue
                }

                s := Status{}
                s.Id = f.Name()
                s.IP, _ = d.Get("ip").String()
                s.Name, _ = d.Get("host_name").String()
                s.Load, _ = d.Get("load").String()
                s.DiskWarn, _ = d.Get("disk_warn").String()
                s.CPURate, _ = d.Get("cpu_rate").Float64()
                s.MemRate, _ = d.Get("mem_rate").Float64()
                s.SwapRate, _ = d.Get("swap_rate").Float64()
                s.DiskRate, _ = d.Get("disk_rate").Float64()
                s.OSRelease, _ = d.Get("os_release").String()

                net_read, _ := d.Get("net_read").Float64()
                net_write, _ := d.Get("net_write").Float64()
                disk_read, _ := d.Get("disk_read").Float64()
                disk_write, _ := d.Get("disk_write").Float64()
                net_total, _ := d.Get("net_total").Float64()
                time_stamp, _ := d.Get("time_stamp").Int()
                uptime, _ := d.Get("uptime").Int()

                s.Uptime = SecondToHumanTime(int(uptime))
                s.OSRelease = PrettyLinuxVersion(s.OSRelease)

                s.NetRead = HumanByte(net_read)
                s.NetWrite = HumanByte(net_write)
                s.DiskRead = HumanByte(disk_read)
                s.DiskWrite = HumanByte(disk_write)
                s.NetTotal = HumanByte(net_total)

                now_date := time.Now().Format("2006-01-02")
                get_date := time.Unix(int64(time_stamp), 0).Format("2006-01-02")
                if now_date == get_date {
                    s.LastUpdate = time.Unix(int64(time_stamp), 0).Format("15:04:05")
                } else {
                    s.LastUpdate = get_date
                }

                s.Status = "success"
                if s.DiskWarn != "" {
                    s.Status = "warning"
                }

                diff_seconds := time.Now().Unix() - int64(time_stamp)
                if diff_seconds > 180 {
                    s.Status = "danger"
                } else if diff_seconds > 120 {
                    s.Status = "warning"
                }

                data = append(data, s)
            }
        }
    }

    sort.Sort(Statuses(data))

    return data
}


func WriteStatus(data_dir string, data *simplejson.Json) (err error) {
    data_id, _ := data.Get("id").String()
    data_id_dir := data_dir + "/" + data_id[:8]
    if !FileExists(data_id_dir) {
        err = os.Mkdir(data_id_dir, 0755)
        if err != nil {
            return
        }
    }

    current, err := simplejson.Load(data_id_dir + "/current")
    if err == nil {
        o_time_stamp, _ := current.Get("time_stamp").Int()
        o_disk_read, _ := current.Get("disk_read").Float64()
        o_disk_write, _ := current.Get("disk_write").Float64()
        o_net_read, _ := current.Get("net_read").Float64()
        o_net_write, _ := current.Get("net_write").Float64()
        o_net_total, _ := current.Get("net_total").Float64()

        n_time_stamp, _ := data.Get("time_stamp").Int()
        n_disk_read, _ := data.Get("disk_read").Float64()
        n_disk_write, _ := data.Get("disk_write").Float64()
        n_net_read, _ := data.Get("net_read").Float64()
        n_net_write, _ := data.Get("net_write").Float64()

        status_set, _ := current.Map()
        diff_seconds := float64(n_time_stamp - o_time_stamp)
        if diff_seconds <= 0 {
            err = errors.New("report time too short")
            return
        }

        status_set["disk_read"] = (n_disk_read - o_disk_read) / diff_seconds
        status_set["disk_write"] = (n_disk_write - o_disk_write) / diff_seconds
        status_set["net_read"] = (n_net_read - o_net_read) / diff_seconds
        status_set["net_write"] = (n_net_write - o_net_write) / diff_seconds

        o_net := o_net_read + o_net_write
        n_net := n_net_read + n_net_write
        diff := n_net
        if n_net >= o_net {
            diff = n_net - o_net
        }

        if (time.Unix(int64(o_time_stamp), 0).Format("2006-01") == time.Unix(int64(n_time_stamp), 0).Format("2006-01")) {
            status_set["net_total"] = o_net_total + diff
        } else {
            status_set["net_total"] = 0
        }
        data.Set("net_total", status_set["net_total"])

        current.Set("time_stamp", n_time_stamp)
        _, err = current.Dump(data_id_dir + "/status")
        if err != nil {
            return
        }
    }

    _, err = data.Dump(data_id_dir + "/current")
    if err != nil {
        return
    }

    return
}
