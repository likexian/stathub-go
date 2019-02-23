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
)


// StatService statSend loop
func StatService() {
    go statSend()
    t := time.NewTicker(60 * time.Second)
    for range t.C {
        go statSend()
    }
}


// statSend get host stat and send to server
func statSend() {
    stat := GetStat(SERVER_CONFIG.Id, SERVER_CONFIG.Name)
    for i:=0; i<3; i++ {
        err := httpSend(SERVER_CONFIG.ServerUrl, SERVER_CONFIG.ServerKey, stat)
        if err != nil {
            fmt.Println(err)
        } else {
            break
        }
    }
}
