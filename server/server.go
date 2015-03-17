package main


import (
    "os"
    "fmt"
    "time"
    "strings"
    "net/http"
    "io/ioutil"
    "text/template"
    "github.com/likexian/simplejson-go"
)


const (
    DATA_DIR = "/data"
    CONFIG_FILE = "/server.json"
)


var (
    CONFIG_ID = ""
    CONFIG_KEY = ""
    CONFIG_PASSWORD = ""
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
    Name        string
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


func IndexHandler(w http.ResponseWriter, r *http.Request) {
    if !IsLogin(w, r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    if r.URL.Path != "/" {
        HTTPErrorHandler(w, r, http.StatusNotFound)
        return
    }

    tpl, err := template.ParseFiles("template/layout.html", "template/index.html")
    if err != nil {
        HTTPErrorHandler(w, r, http.StatusInternalServerError)
        return
    }

    data := []Status{}
    files, err := ioutil.ReadDir(SERVER_WORKDIR + DATA_DIR)
    if err == nil {
        for _, f := range files {
            if FileExists(SERVER_WORKDIR + DATA_DIR + "/" + f.Name() + "/status") {
                d, err := simplejson.Load(SERVER_WORKDIR + DATA_DIR + "/" + f.Name() + "/status")
                if err != nil {
                    continue
                }

                s := Status{}
                s.IP, _ = d.Get("ip").String()
                s.Name, _ = d.Get("host_name").String()
                uptime, _ := d.Get("uptime").Int()
                s.NetRead, _ = d.Get("net_read").Int()
                s.Load, _ = d.Get("load").String()
                s.NetWrite, _ = d.Get("net_write").Int()
                s.DiskRead, _ = d.Get("disk_read").Int()
                s.DiskWrite, _ = d.Get("disk_write").Int()
                s.DiskWarn, _ = d.Get("disk_warn").String()
                s.CPURate, _ = d.Get("cpu_rate").Float64()
                s.MemRate, _ = d.Get("mem_rate").Float64()
                s.SwapRate, _ = d.Get("swap_rate").Float64()
                s.DiskRate, _ = d.Get("disk_rate").Float64()
                s.OSRelease, _ = d.Get("os_release").String()
                time_stamp, _ := d.Get("time_stamp").Int()

                s.Uptime = SecondToHumanTime(int(uptime))
                s.OSRelease = PrettyLinuxVersion(s.OSRelease)

                now_date := time.Now().Format("2006-01-02")
                get_date := time.Unix(int64(time_stamp), 0).Format("2006-01-02")
                if now_date == get_date {
                    s.LastUpdate = time.Unix(int64(time_stamp), 0).Format("15:04:05")
                } else {
                    s.LastUpdate = get_date
                }

                if s.DiskWarn != "" {
                    s.Status = "warning"
                }

                diff_seconds := time.Now().Unix() - int64(time_stamp)
                if diff_seconds > 120 {
                    s.Status = "danger"
                } else if diff_seconds > 90 {
                    s.Status = "warning"
                }

                data = append(data, s)
            }
        }
    }

    tpl.Execute(w, data)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        password := r.PostForm.Get("password")
        if PassWord(CONFIG_KEY, password) != CONFIG_PASSWORD {
            http.Redirect(w, r, "/login", http.StatusFound)
        } else {
            value := PassWord(CONFIG_KEY, CONFIG_PASSWORD)
            cookie := http.Cookie{Name: "id", Value: value, HttpOnly: true}
            http.SetCookie(w, &cookie)
            http.Redirect(w, r, "/", http.StatusFound)
        }
    } else {
        tpl, err := template.ParseFiles("template/layout.html", "template/login.html")
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        tpl.Execute(w, map[string]string{"action": "login"})
    }
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
    name, _ := data.Get("host_name").String()
    data.Set("host_name", strings.Split(name, ".")[0])

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


func BootstrapCSSHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/css; charset=utf-8")
    http.ServeFile(w, r, r.URL.Path[1:])
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


func IsLogin(w http.ResponseWriter, r *http.Request) bool {
    cookie, err := r.Cookie("id")
    if err != nil || cookie.Value == "" {
        return false
    } else {
        value := PassWord(CONFIG_KEY, CONFIG_PASSWORD)
        if value != cookie.Value {
            return false
        }
    }

    return true
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
    CONFIG_PASSWORD = "7be84a051a01334edf5cf935cad4cc6c"

    http.HandleFunc("/", IndexHandler)
    http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/static/bootstrap.css", BootstrapCSSHandler)
    http.HandleFunc("/api/stat", APIStatHandler)

    err = http.ListenAndServe(":15944", nil)
    if err != nil {
        panic(err)
    }
}
