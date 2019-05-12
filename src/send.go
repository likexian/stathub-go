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
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/likexian/simplejson-go"
	"io/ioutil"
	"net/http"
	"time"
)

// httpSend send data to stat api
func httpSend(server, key, stat string) (err error) {
	surl := server + "/api/stat"
	skey := Md5(key, stat)

	request, err := http.NewRequest("POST", surl, bytes.NewBuffer([]byte(stat)))
	if err != nil {
		return
	}

	request.Header.Set("X-Client-Key", skey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Stat Hub API Client/"+Version()+" (i@likexian.com)")

	tr := &http.Transport{
		// If not self-signed certificate please disabled this.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(30 * time.Second),
		Transport: tr,
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	jsonData, err := simplejson.Loads(string(data))
	if err != nil {
		return
	}

	status := jsonData.Get("status.code").MustInt(0)
	if status != 1 {
		message := jsonData.Get("status.message").MustString("unknown error")
		return errors.New("server return: " + message)
	}

	return
}
