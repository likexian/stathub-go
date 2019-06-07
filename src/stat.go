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
	"fmt"
	"time"

	"github.com/likexian/gokit/xhuman"
	"github.com/likexian/host-stat-go"
	"github.com/likexian/simplejson-go"
)

// Stat storing stat data
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

// GetStat return stat data
func GetStat(id string, name string) (result string, err error) {
	stat := Stat{}
	stat.Id = id
	stat.TimeStamp = time.Now().Unix()

	hostInfo, err := hoststat.GetHostInfo()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get host info failed: %s", err.Error())
	}
	stat.OSRelease = hostInfo.Release + " " + hostInfo.OSBit

	if name == "" {
		stat.HostName = hostInfo.HostName
	} else {
		stat.HostName = name
	}

	cpuInfo, err := hoststat.GetCPUInfo()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get cpu info failed: %s", err.Error())
	}
	stat.CPUName = cpuInfo.ModelName
	stat.CPUCore = cpuInfo.CoreCount

	cpuStat, err := hoststat.GetCPUStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get cpu stat failed: %s", err.Error())
	}
	stat.CPURate = xhuman.Round(100-cpuStat.IdleRate, 2)

	memStat, err := hoststat.GetMemStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get mem stat failed: %s", err.Error())
	}
	stat.MemRate = memStat.MemRate
	stat.SwapRate = memStat.SwapRate

	diskStat, err := hoststat.GetDiskStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get disk stat failed: %s", err.Error())
	}
	diskTotal := uint64(0)
	diskUsed := uint64(0)
	for _, v := range diskStat {
		diskTotal += v.Total
		diskUsed += v.Used
		if v.UsedRate > 90 {
			stat.DiskWarn += fmt.Sprintf("%s %.2f%%;", v.Mount, v.UsedRate)
		}
	}
	stat.DiskRate = xhuman.Round(float64(diskUsed)*100/float64(diskTotal), 2)

	ioStat, err := hoststat.GetIOStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get io stat failed: %s", err.Error())
	}
	diskRead := uint64(0)
	diskWrite := uint64(0)
	for _, v := range ioStat {
		diskRead += v.ReadBytes
		diskWrite += v.WriteBytes
	}
	stat.DiskRead = diskRead
	stat.DiskWrite = diskWrite

	netStat, err := hoststat.GetNetStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get net stat failed: %s", err.Error())
	}
	netWrite := uint64(0)
	netRead := uint64(0)
	for _, v := range netStat {
		if v.Device != "lo" {
			netWrite += v.TXBytes
			netRead += v.RXBytes
		}
	}
	stat.NetWrite = netWrite
	stat.NetRead = netRead

	uptimeStat, err := hoststat.GetUptimeStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get uptime stat failed: %s", err.Error())
	}
	stat.Uptime = uint64(uptimeStat.Uptime)

	loadStat, err := hoststat.GetLoadStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get load stat failed: %s", err.Error())
	}
	stat.Load = fmt.Sprintf("%.2f %.2f %.2f", loadStat.LoadNow, loadStat.LoadPre, loadStat.LoadFar)

	data := simplejson.New(stat)
	result, err = data.Dumps()

	return
}
