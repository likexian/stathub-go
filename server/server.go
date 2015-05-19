/*
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 *
 * Copyright 2015, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "text/template"
    "time"
    "github.com/likexian/simplejson-go"
)

const (
    DEBUG        = false
    DATA_DIR     = "/data"
    CONFIG_FILE  = "/server.json"
    CLIENT_FILE  = "/client"
    PROCESS_USER = "nobody"
    PROCESS_LOCK = "/stathub.pid"
    PROCESS_LOG  = "/stathub.log"
    TLS_CERT     = "/cert/cert.pem"
    TLS_KEY      = "/cert/cert.key"
)

var (
    CONFIG_ID       = ""
    CONFIG_KEY      = ""
    CONFIG_PASSWORD = ""
    CONFIG_ISTLS    = false
)

var (
    SERVER_WORKDIR = ""
    SERVER_START   = int64(0)
)

type Config struct {
    Id       string `json:"id"`
    Key      string `json:"key"`
    PassWord string `json:"password"`
    IsTLS    bool   `json:"istls"`
}

type Status struct {
    Id         string
    IP         string
    Name       string
    Status     string
    Uptime     string
    Load       string
    NetRead    string
    NetWrite   string
    DiskRead   string
    DiskWrite  string
    DiskWarn   string
    CPURate    float64
    MemRate    float64
    SwapRate   float64
    DiskRate   float64
    NetTotal   string
    OSRelease  string
    LastUpdate string
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

    tpl, err := template.New("index").Parse(Template_Layout)
    if err != nil {
        HTTPErrorHandler(w, r, http.StatusInternalServerError)
        return
    }

    tpl, err = tpl.Parse(Template_Index)
    if err != nil {
        HTTPErrorHandler(w, r, http.StatusInternalServerError)
        return
    }

    if DEBUG {
        tpl, err = template.ParseFiles("template/layout.html", "template/index.html")
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }
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
                s.Id = f.Name()
                s.IP, _ = d.Get("ip").String()
                s.Name, _ = d.Get("host_name").String()
                s.Load, _ = d.Get("load").String()
                s.DiskWarn, _ = d.Get("disk_warn").String()
                s.CPURate, _ = d.Get("cpu_rate").Float64()
                s.MemRate, _ = d.Get("mem_rate").Float64()
                s.SwapRate, _ = d.Get("swap_rate").Float64()
                s.DiskRate, _ = d.Get("disk_rate").Float64()
                s.OSRelease, _ = d.Get("os_release").String()

                net_read, _ := d.Get("net_read").Float64()
                net_write, _ := d.Get("net_write").Float64()
                disk_read, _ := d.Get("disk_read").Float64()
                disk_write, _ := d.Get("disk_write").Float64()
                net_total, _ := d.Get("net_total").Float64()
                time_stamp, _ := d.Get("time_stamp").Int()
                uptime, _ := d.Get("uptime").Int()

                s.Uptime = SecondToHumanTime(int(uptime))
                s.OSRelease = PrettyLinuxVersion(s.OSRelease)

                s.NetRead = HumanByte(net_read)
                s.NetWrite = HumanByte(net_write)
                s.DiskRead = HumanByte(disk_read)
                s.DiskWrite = HumanByte(disk_write)
                s.NetTotal = HumanByte(net_total)

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
                if diff_seconds > 180 {
                    s.Status = "danger"
                } else if diff_seconds > 120 {
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
        tpl, err := template.New("login").Parse(Template_Layout)
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        tpl, err = tpl.Parse(Template_Login)
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        if DEBUG {
            tpl, err = template.ParseFiles("template/layout.html", "template/login.html")
            if err != nil {
                HTTPErrorHandler(w, r, http.StatusInternalServerError)
                return
            }
        }

        tpl.Execute(w, map[string]string{"action": "login"})
    }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    expires := time.Now()
    expires = expires.AddDate(0, 0, -1)
    cookie := http.Cookie{Name: "id", Value: "", Expires: expires, HttpOnly: true}
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/login", http.StatusFound)
    return
}

func PasswdHandler(w http.ResponseWriter, r *http.Request) {
    if !IsLogin(w, r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        password := r.PostForm.Get("password")
        if password == "" {
            http.Redirect(w, r, "/passwd", http.StatusFound)
        } else {
            CONFIG_PASSWORD = PassWord(CONFIG_KEY, password)
            WriteConfig(CONFIG_ID, CONFIG_KEY, CONFIG_PASSWORD, CONFIG_ISTLS)
            http.Redirect(w, r, "/", http.StatusFound)
        }
    } else {
        tpl, err := template.New("passwd").Parse(Template_Layout)
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        tpl, err = tpl.Parse(Template_Login)
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }

        if DEBUG {
            tpl, err = template.ParseFiles("template/layout.html", "template/login.html")
            if err != nil {
                HTTPErrorHandler(w, r, http.StatusInternalServerError)
                return
            }
        }

        tpl.Execute(w, map[string]string{"action": "passwd"})
    }
}

func HelpHandler(w http.ResponseWriter, r *http.Request) {
    if !IsLogin(w, r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    tpl, err := template.New("help").Parse(Template_Layout)
    if err != nil {
        HTTPErrorHandler(w, r, http.StatusInternalServerError)
        return
    }

    tpl, err = tpl.Parse(Template_Help)
    if err != nil {
        HTTPErrorHandler(w, r, http.StatusInternalServerError)
        return
    }

    if DEBUG {
        tpl, err = template.ParseFiles("template/layout.html", "template/help.html")
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }
    }

    tpl.Execute(w, map[string]string{"server": r.Host, "key": CONFIG_KEY})
}

func NodeHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key != CONFIG_KEY {
        HTTPErrorHandler(w, r, http.StatusForbidden)
        return
    }

    tpl, err := template.New("node").Parse(Template_Node)
    if err != nil {
        HTTPErrorHandler(w, r, http.StatusInternalServerError)
        return
    }

    if DEBUG {
        tpl, err = template.ParseFiles("template/node.html")
        if err != nil {
            HTTPErrorHandler(w, r, http.StatusInternalServerError)
            return
        }
    }

    id := PassWord(CONFIG_KEY, fmt.Sprintf("%s", time.Now().Unix()))
    tpl.Execute(w, map[string]string{"id": id, "server": r.Host, "key": CONFIG_KEY})
}

func Client32Handler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, SERVER_WORKDIR+CLIENT_FILE+"_i686")
}

func Client64Handler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, SERVER_WORKDIR+CLIENT_FILE+"_x86_64")
}

func APIStatHandler(w http.ResponseWriter, r *http.Request) {
    ip := strings.Split(r.RemoteAddr, ":")[0]
    if test, ok := r.Header["X-Real-Ip"]; ok {
        ip = test[0]
    }

    client_key := ""
    if test, ok := r.Header["X-Client-Key"]; !ok {
        fmt.Fprintf(w, "Key invalid")
        return
    } else {
        client_key = test[0]
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "Body error")
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
        fmt.Fprintf(w, "Body invalid")
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
    if err == nil {
        o_time_stamp, _ := current.Get("time_stamp").Int()
        o_disk_read, _ := current.Get("disk_read").Float64()
        o_disk_write, _ := current.Get("disk_write").Float64()
        o_net_read, _ := current.Get("net_read").Float64()
        o_net_write, _ := current.Get("net_write").Float64()
        o_net_total, _ := current.Get("net_total").Float64()

        n_time_stamp, _ := data.Get("time_stamp").Int()
        n_disk_read, _ := data.Get("disk_read").Float64()
        n_disk_write, _ := data.Get("disk_write").Float64()
        n_net_read, _ := data.Get("net_read").Float64()
        n_net_write, _ := data.Get("net_write").Float64()

        status_set, _ := current.Map()
        diff_seconds := float64(n_time_stamp - o_time_stamp)
        if diff_seconds <= 0 {
            return
        }

        status_set["disk_read"] = (n_disk_read - o_disk_read) / diff_seconds
        status_set["disk_write"] = (n_disk_write - o_disk_write) / diff_seconds
        status_set["net_read"] = (n_net_read - o_net_read) / diff_seconds
        status_set["net_write"] = (n_net_write - o_net_write) / diff_seconds

        o_net := o_net_read + o_net_write
        n_net := n_net_read + n_net_write
        diff := n_net
        if n_net >= o_net {
            diff = n_net - o_net
        }

        if (time.Unix(int64(o_time_stamp), 0).Format("2006-01") == time.Unix(int64(n_time_stamp), 0).Format("2006-01")) {
            status_set["net_total"] = o_net_total + diff
        } else {
            status_set["net_total"] = 0
        }
        data.Set("net_total", status_set["net_total"])

        current.Set("time_stamp", n_time_stamp)
        simplejson.Dump(data_id_dir + "/status", current)
    }

    simplejson.Dump(data_id_dir + "/current", data)

    return
}

func HTTPErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    if status == http.StatusForbidden {
        fmt.Fprint(w, "<title>Forbidden</title><h1>Forbidden</h1>")
    } else if status == http.StatusNotFound {
        fmt.Fprint(w, "<title>Not Found</title><h1>Not Found</h1>")
    } else if status == http.StatusInternalServerError {
        fmt.Fprint(w, "<title>Internal Server Error</title><h1>Internal Server Error</h1>")
    }
}

func BootstrapCSSHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/css; charset=utf-8")
    if DEBUG {
        http.ServeFile(w, r, r.URL.Path[1:])
    } else {
        fmt.Fprint(w, Template_Bootstrap)
    }
}

func WriteConfig(id, key, password string, istls bool) {
    config := Config{}
    config.Id = id
    config.Key = key
    config.PassWord = password
    config.IsTLS = istls

    data := simplejson.Json{}
    data.Data = config
    simplejson.Dump(SERVER_WORKDIR+CONFIG_FILE, &data)
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
    uid, gid, err := LookupUser(PROCESS_USER)
    if err != nil {
        panic(err)
    }

    pwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    SERVER_WORKDIR = pwd

    if !FileExists(SERVER_WORKDIR + DATA_DIR) {
        err := os.Mkdir(SERVER_WORKDIR+DATA_DIR, 0755)
        if err != nil {
            return
        }
    }
    os.Chown(SERVER_WORKDIR + DATA_DIR, uid, gid)

    if !FileExists(SERVER_WORKDIR + DATA_DIR + PROCESS_LOG) {
        WriteFile(SERVER_WORKDIR + DATA_DIR + PROCESS_LOG, "")
    }
    // os.Chown(SERVER_WORKDIR + DATA_DIR + PROCESS_LOG, uid, gid)

    SERVER_START = time.Now().Unix()
    if !FileExists(SERVER_WORKDIR + CONFIG_FILE) {
        time_stamp := fmt.Sprintf("%d", SERVER_START)
        id := PassWord(fmt.Sprintf("%s", os.Getpid()), time_stamp)
        key := PassWord(id, time_stamp)
        password := PassWord(key, "likexian")
        istls := true
        WriteConfig(id, key, password, istls)
    }

    cert_dir := strings.Split(TLS_CERT, "/")[1]
    if !FileExists(SERVER_WORKDIR + "/" + cert_dir) {
        err := os.Mkdir(SERVER_WORKDIR + "/" + cert_dir, 0755)
        if err != nil {
            panic(err)
        }
    }

    if !FileExists(SERVER_WORKDIR + TLS_CERT) {
        err := WriteFile(SERVER_WORKDIR + TLS_CERT, Default_TLS_CERT)
        if err != nil {
            panic(err)
        }
    }

    if !FileExists(SERVER_WORKDIR + TLS_KEY) {
        err := WriteFile(SERVER_WORKDIR + TLS_KEY, Default_TLS_KEY)
        if err != nil {
            panic(err)
        }
    }

    if !DEBUG {
        daemon := Daemon(SERVER_WORKDIR + DATA_DIR + PROCESS_LOCK, SERVER_WORKDIR + DATA_DIR + PROCESS_LOG, uid, gid, 0, 0)
        if daemon != 0 {
            os.Exit(-1)
        }
    }

    config, err := simplejson.Load(SERVER_WORKDIR + CONFIG_FILE)
    if err != nil {
        panic(err)
    }

    CONFIG_ID, _ = config.Get("id").String()
    CONFIG_KEY, _ = config.Get("key").String()
    CONFIG_PASSWORD, _ = config.Get("password").String()
    CONFIG_ISTLS, _ = config.Get("istls").Bool()

    http.HandleFunc("/", IndexHandler)
    http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/logout", LogoutHandler)
    http.HandleFunc("/passwd", PasswdHandler)
    http.HandleFunc("/help", HelpHandler)
    http.HandleFunc("/node", NodeHandler)
    http.HandleFunc("/static/client_i686", Client32Handler)
    http.HandleFunc("/static/client_x86_64", Client64Handler)
    http.HandleFunc("/static/bootstrap.css", BootstrapCSSHandler)
    http.HandleFunc("/api/stat", APIStatHandler)

    if CONFIG_ISTLS {
        err = http.ListenAndServeTLS(":15944", SERVER_WORKDIR + TLS_CERT, SERVER_WORKDIR + TLS_KEY, nil)
        if err != nil {
            panic(err)
        }
    } else {
        err = http.ListenAndServe(":15944", nil)
        if err != nil {
            panic(err)
        }
    }
}
