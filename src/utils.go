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
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Byte units
const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

// HumanByte returns readable string for bytes
func HumanByte(bytes float64) string {
	unit := ""
	value := bytes

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	return fmt.Sprintf("%.1f%s", Round(value, 1), unit)
}

// Round returns round number
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

// FileExists returns file is exists
func FileExists(fname string) bool {
	_, err := os.Stat(fname)
	return !os.IsNotExist(err)
}

// ReadFile return text of file
func ReadFile(fname string) (result string, err error) {
	text, err := ioutil.ReadFile(fname)
	if err != nil {
		return
	}
	result = string(text)

	return
}

// WriteFile write text to file
func WriteFile(fname, text string) (err error) {
	ds, _ := filepath.Split(fname)
	if ds != "" && !FileExists(ds) {
		err = os.MkdirAll(ds, 0755)
		if err != nil {
			return
		}
	}

	return ioutil.WriteFile(fname, []byte(text), 0644)
}

// SecondToHumanTime returns readable string for seconds
func SecondToHumanTime(second int) string {
	if second < 60 {
		return fmt.Sprintf("%d sec", second)
	} else if second < 3600 {
		return fmt.Sprintf("%d min", uint64(second/60))
	} else if second < 86400 {
		return fmt.Sprintf("%d hours", uint64(second/3600))
	}

	return fmt.Sprintf("%d days", uint64(second/86400))
}

// PrettyLinuxVersion returns readable linux version
func PrettyLinuxVersion(version string) string {
	find := strings.Index(version, "(")
	if find != -1 {
		version = version[:find]
	}
	version = strings.Replace(version, "release", "", -1)
	version = strings.Replace(version, "GNU", "", -1)
	version = strings.Replace(version, "LINUX", "", -1)
	version = strings.Replace(version, "Linux", "", -1)
	version = strings.Replace(version, "LTS", "", -1)
	version = strings.Replace(version, "/", "", -1)
	version = strings.Replace(version, "  ", " ", -1)
	return version
}

// Md5 returns hex md5 of string
func Md5(str, key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str+key)))
}

// RawInput returns read args from stdin
func RawInput(message string, allowEmpty bool) (result string) {
	fmt.Println(message)

	reader := bufio.NewReader(os.Stdin)
	result, _ = reader.ReadString('\n')
	result = strings.TrimSpace(result)

	if !allowEmpty && result == "" {
		fmt.Println("No data inputed")
		os.Exit(1)
	}

	return
}
