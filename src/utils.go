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

// Chown do recurse chown go file or folder
func Chown(fname string, uid, gid int) (err error) {
	isDir, err := IsDir(fname)
	if err != nil {
		return
	}

	err = os.Chown(fname, uid, gid)
	if err != nil || !isDir {
		return
	}

	if !strings.HasSuffix(fname, "/") {
		fname += "/"
	}

	fs, err := ioutil.ReadDir(fname)
	if err != nil {
		return
	}

	for _, f := range fs {
		err = Chown(fname+f.Name(), uid, gid)
		if err != nil {
			return
		}
	}

	return
}

// IsDir returns if path is a dir
func IsDir(fname string) (bool, error) {
	f, err := os.Stat(fname)
	if err != nil {
		return false, err
	}

	return f.Mode().IsDir(), nil
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
