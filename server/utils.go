package main


import(
    "fmt"
    "os"
    "crypto/md5"
)


func FileExists(fname string) bool {
    _, err := os.Stat(fname)
    return !os.IsNotExist(err)
}

func PassWord(key, password string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(key + password)))
}
