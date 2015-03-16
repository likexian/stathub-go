package main


import(
    "fmt"
    "os"
    "strings"
    "crypto/md5"
)


func FileExists(fname string) bool {
    _, err := os.Stat(fname)
    return !os.IsNotExist(err)
}

func PassWord(key, password string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(key + password)))
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
