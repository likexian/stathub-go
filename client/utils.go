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
    "os"
    "fmt"
    "math"
    "bufio"
    "strings"
    "crypto/md5"
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


func RawInput(message string, allow_empty bool) (result string) {
    fmt.Println(message)

    reader := bufio.NewReader(os.Stdin)
    result, _ = reader.ReadString('\n')
    result = strings.Trim(strings.Trim(result, "\n"), " ")

    if !allow_empty && result == "" {
        fmt.Println("No data inputed\n")
        os.Exit(1)
    }

    return
}


func PassWord(key, password string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(key + password)))
}
