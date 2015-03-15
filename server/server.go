package main


import (
    "os"
    "fmt"
    "time"
    "strings"
    "net/http"
    "io/ioutil"
    "github.com/likexian/simplejson-go"
)


const (
    DATA_DIR = "/data"
    CONFIG_FILE = "/server.json"
)


var (
    CONFIG_ID = ""
    CONFIG_KEY = ""
)


var (
    SERVER_WORKDIR = ""
    SERVER_START = int64(0)
)


type Config struct {
    Id          string `json:"id"`
    Key         string `json:"key"`
}


type Status struct {
    IP          string
    Status      string
    Uptime      string
    Load        string
    NetRead     int
    NetWrite    int
    DiskRead    int
    DiskWrite   int
    DiskWarn    string
    CPURate     float64
    MemRate     float64
    SwapRate    float64
    DiskRate    float64
    OSRelease   string
    LastUpdate  string
}


func APIStatHandler(w http.ResponseWriter, r *http.Request) {
    ip := strings.Split(r.RemoteAddr, ":")[0]
    client_key := ""
    if test, ok := r.Header["X-Client-Key"]; !ok {
        fmt.Fprintf(w, "Key invalid")
        return
    } else {
        client_key = test[0]
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return
    }

    text := string(body)
    server_key := PassWord(CONFIG_KEY, text)
    if server_key != client_key {
        fmt.Fprintf(w, "Key invalid")
        return
    }

    data, err := simplejson.Loads(text)
    if err != nil {
        return
    }

    data.Set("ip", ip)

    data_id, _ := data.Get("id").String()
    data_id_dir := SERVER_WORKDIR + DATA_DIR + "/" + data_id[:8]
    if !FileExists(data_id_dir) {
        err := os.Mkdir(data_id_dir, 0755)
        if err != nil {
            return
        }
    }

    current, err := simplejson.Load(data_id_dir + "/current")
    simplejson.Dump(data_id_dir + "/current", data)

    if err == nil {
        o_time_stamp, _ := current.Get("time_stamp").Int()
        o_disk_read, _ := current.Get("disk_read").Int()
        o_disk_write, _ := current.Get("disk_write").Int()
        o_net_read, _ := current.Get("net_read").Int()
        o_net_write, _ := current.Get("net_write").Int()

        n_time_stamp, _ := data.Get("time_stamp").Int()
        n_disk_read, _ := data.Get("disk_read").Int()
        n_disk_write, _ := data.Get("disk_write").Int()
        n_net_read, _ := data.Get("net_read").Int()
        n_net_write, _ := data.Get("net_write").Int()

        status_set, _ := current.Map()
        diff_seconds := n_time_stamp - o_time_stamp
        if diff_seconds == 0 {
            return
        }

        status_set["disk_read"] = (n_disk_read - o_disk_read) / diff_seconds
        status_set["disk_write"] = (n_disk_write - o_disk_write) / diff_seconds
        status_set["net_read"] = (n_net_read - o_net_read) / diff_seconds
        status_set["net_write"] = (n_net_write - o_net_write) / diff_seconds

        current.Set("time_stamp", n_time_stamp)
        simplejson.Dump(data_id_dir + "/status", current)
    }

    return
}


func HTTPErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        fmt.Fprint(w, "<title>Not Found</title><h1>Not Found</h1>")
    } else if status == http.StatusInternalServerError {
        fmt.Fprint(w, "<title>Internal Server Error</title><h1>Internal Server Error</h1>")
    }
}


func WriteConfig(id, key string) {
    config := Config{}
    config.Id = id
    config.Key = key

    data := simplejson.Json{}
    data.Data = config
    simplejson.Dump(SERVER_WORKDIR + CONFIG_FILE, &data)
}


func main() {
    pwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    SERVER_WORKDIR = pwd

    SERVER_START = time.Now().Unix()
    if !FileExists(SERVER_WORKDIR + CONFIG_FILE) {
        time_stamp := fmt.Sprintf("%d", SERVER_START)
        id := PassWord(fmt.Sprintf("%s", os.Getpid()), time_stamp)
        key := PassWord(id, time_stamp)
        WriteConfig(id, key)
    }

    config, err := simplejson.Load(SERVER_WORKDIR + CONFIG_FILE)
    if err != nil {
        panic(err)
    }

    CONFIG_ID, _ = config.Get("id").String()
    CONFIG_KEY, _ = config.Get("key").String()

    http.HandleFunc("/api/stat", APIStatHandler)

    err = http.ListenAndServe(":15944", nil)
    if err != nil {
        panic(err)
    }
}
