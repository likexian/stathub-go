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
	"time"
)

// StatService statSend loop
func StatService() {
	SERVER_LOGGER.Info("start stat service")
	go statSend()
	t := time.NewTicker(60 * time.Second)
	for range t.C {
		go statSend()
	}
}

// statSend get host stat and send to server
func statSend() {
	stat, err := GetStat(SERVER_CONFIG.Id, SERVER_CONFIG.Name)
	if err != nil {
		SERVER_LOGGER.Error("get stat failed: %s", err.Error())
		return
	}

	SERVER_LOGGER.Debug("get stat data: %s", stat)
	for i := 0; i < 3; i++ {
		err := httpSend(SERVER_CONFIG.ServerUrl, SERVER_CONFIG.ServerKey, stat)
		if err != nil {
			SERVER_LOGGER.Error("send stat failed, %s", err.Error())
			time.Sleep(3 * time.Second)
		} else {
			SERVER_LOGGER.Debug("send stat to server successful")
			break
		}
	}
}
