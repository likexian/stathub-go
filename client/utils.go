package main


import(
    "fmt"
    "math"
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


func PassWord(key, password string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(key + password)))
}
