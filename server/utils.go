/*
 * A smart Hub for holding server stat
 * http://www.likexian.com/
 *
 * Copyright 2015, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package main


import(
    "fmt"
    "os"
    "math"
    "errors"
    "strings"
    "strconv"
    "crypto/md5"
    "io/ioutil"
)


func Round(data float64, precision int) (result float64) {
    pow := math.Pow(10, float64(precision))
    digit := pow * data
    _, div := math.Modf(digit)

    if div >= 0.5 {
        result = math.Ceil(digit)
    } else {
        result = math.Floor(digit)
    }
    result = result / pow

    return
}


func FileExists(fname string) bool {
    _, err := os.Stat(fname)
    return !os.IsNotExist(err)
}


func ReadFile(fname string) (result string, err error) {
    text, err := ioutil.ReadFile(fname)
    if err != nil {
        return
    }
    result = string(text)

    return
}


func WriteFile(fname, text string) (err error) {
    content := []byte(text)
    err = ioutil.WriteFile(fname, content, 0644)
    return
}


func SecondToHumanTime(second int) (string) {
    if second < 60 {
        return fmt.Sprintf("%d sec", second)
    } else if second < 3600 {
        return fmt.Sprintf("%d min", uint64(second / 60))
    } else if second < 86400 {
        return fmt.Sprintf("%d hours", uint64(second / 3600))
    } else {
        return fmt.Sprintf("%d days", uint64(second / 86400))
    }
}


func PrettyLinuxVersion(version string) (string) {
    find := strings.Index(version, "(")
    if find != -1 {
        version = version[:find]
    }
    version = strings.Replace(version, "release", "", -1)
    version = strings.Replace(version, "GNU", "", -1)
    version = strings.Replace(version, "LINUX", "", -1)
    version = strings.Replace(version, "Linux", "", -1)
    version = strings.Replace(version, "/", "", -1)
    version = strings.Replace(version, "  ", " ", -1)
    return version
}


func PassWord(key, password string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(key + password)))
}


func LookupUser(name string) (uid, gid int, err error) {
    text, err := ReadFile("/etc/passwd")
    if err != nil {
        return
    }

    s_uid := ""
    s_gid := ""
    lines := strings.Split(text, "\n")
    for _, v := range lines {
        ls := strings.Split(v, ":")
        if ls[0] == name {
            s_uid = ls[2]
            s_gid = ls[3]
        }
    }

    if s_uid == "" || s_gid == "" {
        err = errors.New("User not exits")
        return
    }

    gid, err = strconv.Atoi(s_gid)
    if err != nil {
        return
    }

    uid, err = strconv.Atoi(s_uid)
    if err != nil {
        return
    }

    return
}
