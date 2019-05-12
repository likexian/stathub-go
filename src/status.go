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
	"github.com/likexian/simplejson-go"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

// Status storing stat data
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

// Statuses storing list of Status
type Statuses []Status

func (n Statuses) Len() int           { return len(n) }
func (n Statuses) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Statuses) Less(i, j int) bool { return n[i].Name < n[j].Name }

// ReadStatus returns read status from dir
func ReadStatus(dataDir string) Statuses {
	data := []Status{}
	files, err := ioutil.ReadDir(dataDir)
	if err == nil {
		for _, f := range files {
			if FileExists(dataDir + "/" + f.Name() + "/status") {
				d, err := simplejson.Load(dataDir + "/" + f.Name() + "/status")
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

				netRead, _ := d.Get("net_read").Float64()
				netWrite, _ := d.Get("net_write").Float64()
				diskRead, _ := d.Get("disk_read").Float64()
				diskWrite, _ := d.Get("disk_write").Float64()
				netTotal, _ := d.Get("net_total").Float64()
				timeStamp, _ := d.Get("time_stamp").Int()
				uptime, _ := d.Get("uptime").Int()

				s.Load = strings.Fields(s.Load)[0]
				s.Uptime = SecondToHumanTime(int(uptime))
				s.OSRelease = PrettyLinuxVersion(s.OSRelease)

				s.CPURate = Round(s.CPURate, 2)
				s.MemRate = Round(s.MemRate, 2)
				s.SwapRate = Round(s.SwapRate, 2)
				s.DiskRate = Round(s.DiskRate, 2)

				s.NetRead = HumanByte(netRead)
				s.NetWrite = HumanByte(netWrite)
				s.DiskRead = HumanByte(diskRead)
				s.DiskWrite = HumanByte(diskWrite)
				s.NetTotal = HumanByte(netTotal)

				nowDate := time.Now().Format("2006-01-02")
				getDate := time.Unix(int64(timeStamp), 0).Format("2006-01-02")
				if nowDate == getDate {
					s.LastUpdate = time.Unix(int64(timeStamp), 0).Format("15:04:05")
				} else {
					s.LastUpdate = getDate
				}

				s.Status = "success"
				if s.DiskWarn != "" {
					s.DiskWarn = strings.Trim(s.DiskWarn, ";")
					s.Status = "warning"
				}

				diffSeconds := time.Now().Unix() - int64(timeStamp)
				if diffSeconds > 180 {
					s.Status = "danger"
				} else if diffSeconds > 120 {
					s.Status = "warning"
				}

				data = append(data, s)
			}
		}
	}

	sort.Sort(Statuses(data))

	return data
}

// WriteStatus write status to dir
func WriteStatus(dataDir string, data *simplejson.Json) (err error) {
	dataId, _ := data.Get("id").String()
	dataIdDir := dataDir + "/" + dataId[:8]
	if !FileExists(dataIdDir) {
		err = os.Mkdir(dataIdDir, 0755)
		if err != nil {
			return
		}
	}

	current, err := simplejson.Load(dataIdDir + "/current")
	if err == nil {
		oTimeStamp, _ := current.Get("time_stamp").Int()
		oDiskRead, _ := current.Get("disk_read").Float64()
		oDiskWrite, _ := current.Get("disk_write").Float64()
		oNetRead, _ := current.Get("net_read").Float64()
		oNetWrite, _ := current.Get("net_write").Float64()
		oNetTotal, _ := current.Get("net_total").Float64()

		nTimeStamp, _ := data.Get("time_stamp").Int()
		nDiskRead, _ := data.Get("disk_read").Float64()
		nDiskWrite, _ := data.Get("disk_write").Float64()
		nNetRead, _ := data.Get("net_read").Float64()
		nNetWrite, _ := data.Get("net_write").Float64()

		statusSet, _ := current.Map()
		diffSeconds := float64(nTimeStamp - oTimeStamp)
		if diffSeconds <= 0 {
			err = errors.New("report time too short")
			return
		}

		statusSet["disk_read"] = (nDiskRead - oDiskRead) / diffSeconds
		statusSet["disk_write"] = (nDiskWrite - oDiskWrite) / diffSeconds
		statusSet["net_read"] = (nNetRead - oNetRead) / diffSeconds
		statusSet["net_write"] = (nNetWrite - oNetWrite) / diffSeconds

		oNet := oNetRead + oNetWrite
		nNet := nNetRead + nNetWrite
		diff := nNet
		if nNet >= oNet {
			diff = nNet - oNet
		}

		if time.Unix(int64(oTimeStamp), 0).Format("2006-01") == time.Unix(int64(nTimeStamp), 0).Format("2006-01") {
			statusSet["net_total"] = oNetTotal + diff
		} else {
			statusSet["net_total"] = 0
		}
		data.Set("net_total", statusSet["net_total"])

		current.Set("time_stamp", nTimeStamp)
		err = current.Dump(dataIdDir + "/status")
		if err != nil {
			return
		}
	}

	err = data.Dump(dataIdDir + "/current")

	return
}
