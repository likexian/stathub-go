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
    "bytes"
    "time"
    "io/ioutil"
    "net/http"
    "crypto/tls"
    "errors"
    "github.com/likexian/simplejson-go"
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
    request.Header.Set("User-Agent", "Stat Hub API Client/" + Version() + " (i@likexian.com)")

    tr := &http.Transport{
        // If not self-signed certificate please disabled this.
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{
        Timeout: time.Duration(30 * time.Second),
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
        return errors.New(message)
    }

    return
}
