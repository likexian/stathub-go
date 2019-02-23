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
    "fmt"
    "os"
    "time"
    "strings"
    "io/ioutil"
    "net/http"
    "text/template"
    "github.com/likexian/simplejson-go"
)


func HttpService() {
    http.HandleFunc("/", IndexHandler)
    http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/logout", LogoutHandler)
    http.HandleFunc("/passwd", PasswdHandler)
    http.HandleFunc("/help", HelpHandler)
    http.HandleFunc("/node", NodeHandler)
    http.HandleFunc("/pkgs/", pkgsHandler)
    http.HandleFunc("/static/", StaticHandler)
    http.HandleFunc("/robots.txt", RobotsTXTHandler)
    http.HandleFunc("/api/stat", APIStatHandler)
    http.HandleFunc("/api/node", APINodeHandler)

    err := http.ListenAndServeTLS(":15944",
        SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSCert, SERVER_CONFIG.BaseDir + SERVER_CONFIG.TLSKey, nil)
    if err != nil {
        panic(err)
    }
}


func IndexHandler(w http.ResponseWriter, r *http.Request) {
    if isRobots(w, r) {
        httpError(w, r, http.StatusForbidden)
        return
    }

    if !isLogin(w, r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    if r.URL.Path != "/" {
        httpError(w, r, http.StatusNotFound)
        return
    }

    tpl, err := template.New("index").Parse(TPL_TEMPLATE["layout.html"])
    if err != nil {
        httpError(w, r, http.StatusInternalServerError)
        return
    }

    tpl, err = tpl.Parse(TPL_TEMPLATE["index.html"])
    if err != nil {
        httpError(w, r, http.StatusInternalServerError)
        return
    }

    if DEBUG {
        tpl, err = template.ParseFiles("template/layout.html", "template/index.html")
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }
    }

    status := ReadStatus(SERVER_CONFIG.DataDir)
    data := map[string]interface{}{
        "data": status,
        "version": TPL_VERSION,
    }
    tpl.Execute(w, data)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if isRobots(w, r) {
        httpError(w, r, http.StatusForbidden)
        return
    }

    if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }

        password := r.PostForm.Get("password")
        if PassWord(SERVER_CONFIG.ServerKey, password) != SERVER_CONFIG.PassWord {
            http.Redirect(w, r, "/login", http.StatusFound)
        } else {
            value := PassWord(SERVER_CONFIG.ServerKey, SERVER_CONFIG.PassWord)
            cookie := http.Cookie{Name: "id", Value: value, HttpOnly: true}
            http.SetCookie(w, &cookie)
            http.Redirect(w, r, "/", http.StatusFound)
        }
    } else {
        tpl, err := template.New("login").Parse(TPL_TEMPLATE["layout.html"])
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }

        tpl, err = tpl.Parse(TPL_TEMPLATE["login.html"])
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }

        if DEBUG {
            tpl, err = template.ParseFiles("template/layout.html", "template/login.html")
            if err != nil {
                httpError(w, r, http.StatusInternalServerError)
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
    if !isLogin(w, r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }

        password := r.PostForm.Get("password")
        if password == "" {
            http.Redirect(w, r, "/passwd", http.StatusFound)
        } else {
            SERVER_CONFIG.PassWord = PassWord(SERVER_CONFIG.ServerKey, password)
            err := SaveConfig(SERVER_CONFIG)
            if err != nil {
                httpError(w, r, http.StatusInternalServerError)
                return
            }
            http.Redirect(w, r, "/", http.StatusFound)
        }
    } else {
        tpl, err := template.New("passwd").Parse(TPL_TEMPLATE["layout.html"])
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }

        tpl, err = tpl.Parse(TPL_TEMPLATE["login.html"])
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }

        if DEBUG {
            tpl, err = template.ParseFiles("template/layout.html", "template/login.html")
            if err != nil {
                httpError(w, r, http.StatusInternalServerError)
                return
            }
        }

        tpl.Execute(w, map[string]string{"action": "passwd"})
    }
}


func HelpHandler(w http.ResponseWriter, r *http.Request) {
    if !isLogin(w, r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }

    tpl, err := template.New("help").Parse(TPL_TEMPLATE["layout.html"])
    if err != nil {
        httpError(w, r, http.StatusInternalServerError)
        return
    }

    tpl, err = tpl.Parse(TPL_TEMPLATE["help.html"])
    if err != nil {
        httpError(w, r, http.StatusInternalServerError)
        return
    }

    if DEBUG {
        tpl, err = template.ParseFiles("template/layout.html", "template/help.html")
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }
    }

    tpl.Execute(w, map[string]string{"server": r.Host, "key": SERVER_CONFIG.ServerKey})
}


func NodeHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key != SERVER_CONFIG.ServerKey {
        httpError(w, r, http.StatusForbidden)
        return
    }

    tpl, err := template.New("node").Parse(TPL_TEMPLATE["node.html"])
    if err != nil {
        httpError(w, r, http.StatusInternalServerError)
        return
    }

    if DEBUG {
        tpl, err = template.ParseFiles("template/node.html")
        if err != nil {
            httpError(w, r, http.StatusInternalServerError)
            return
        }
    }

    tpl.Execute(w, map[string]string{"server": r.Host, "key": SERVER_CONFIG.ServerKey, "version": TPL_VERSION})
}


func RobotsTXTHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "User-agent: *\r\nDisallow: /")
}


func APINodeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    if !isLogin(w, r) {
        result := `{"status": {"code": 0, "message": "login timeout"}}`
        fmt.Fprintf(w, result)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        result := `{"status": {"code": 0, "message": "data error"}}`
        fmt.Fprintf(w, result)
        return
    }

    data, err := simplejson.Loads(string(body))
    if err != nil {
        result := `{"status": {"code": 0, "message": "data invalid"}}`
        fmt.Fprintf(w, result)
        return
    }

    data_id, _ := data.Get("id").String()
    data_id_dir := SERVER_CONFIG.BaseDir + SERVER_CONFIG.DataDir + "/" + data_id[3:]
    if !FileExists(data_id_dir) {
        result := `{"status": {"code": 0, "message": "node id invalid"}}`
        fmt.Fprintf(w, result)
        return
    }

    err = os.RemoveAll(data_id_dir)
    if err != nil {
        result := `{"status": {"code": 0, "message": "` + err.Error() + `"}}`
        fmt.Fprintf(w, result)
        return
    }

    result := `{"status": {"code": 1, "message": "ok"}}`
    fmt.Fprintf(w, result)
}


func APIStatHandler(w http.ResponseWriter, r *http.Request) {
    ip := strings.Split(r.RemoteAddr, ":")[0]
    if test, ok := r.Header["X-Real-Ip"]; ok {
        ip = test[0]
    }

    client_key := ""
    if test, ok := r.Header["X-Client-Key"]; !ok {
        result := `{"status": {"code": 0, "message": "key invalid"}}`
        fmt.Fprintf(w, result)
        return
    } else {
        client_key = test[0]
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        result := `{"status": {"code": 0, "message": "body invalid"}}`
        fmt.Fprintf(w, result)
        return
    }

    text := string(body)
    server_key := PassWord(SERVER_CONFIG.ServerKey, text)
    if server_key != client_key {
        result := `{"status": {"code": 0, "message": "key invalid"}}`
        fmt.Fprintf(w, result)
        return
    }

    data, err := simplejson.Loads(text)
    if err != nil {
        result := `{"status": {"code": 0, "message": "body invalid"}}`
        fmt.Fprintf(w, result)
        return
    }

    data.Set("ip", ip)
    name, _ := data.Get("host_name").String()
    data.Set("host_name", strings.Split(name, ".")[0])

    err = WriteStatus(SERVER_CONFIG.DataDir, data)
    if err != nil {
        result := `{"status": {"code": 0, "message": "` + err.Error() + `"}}`
        fmt.Fprintf(w, result)
        return
    }

    result := `{"status": {"code": 1, "message": "ok"}}`
    fmt.Fprintf(w, result)
}


func StaticHandler(w http.ResponseWriter, r *http.Request) {
    n := strings.LastIndex(r.URL.Path, ".")
    if n == -1 {
        httpError(w, r, http.StatusNotFound)
        return
    }

    ext := r.URL.Path[n + 1:]
    if ext == "css" {
        w.Header().Set("Content-Type", "text/css; charset=utf-8")
    } else if (ext == "js") {
        w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
    }

    if DEBUG {
        http.ServeFile(w, r, r.URL.Path[1:])
    } else {
        if test, ok := TPL_STATIC[r.URL.Path[8:]]; ok {
            fmt.Fprint(w, test)
        } else {
            httpError(w, r, http.StatusNotFound)
        }
    }
}


func pkgsHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
}


func getHTTPHeader(r *http.Request, name string) string {
    if line, ok := r.Header[name]; ok {
        return line[0]
    }

    return ""
}


func httpError(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    if status == http.StatusForbidden {
        fmt.Fprint(w, "<title>Forbidden</title><h1>Forbidden</h1>")
    } else if status == http.StatusNotFound {
        fmt.Fprint(w, "<title>Not Found</title><h1>Not Found</h1>")
    } else if status == http.StatusInternalServerError {
        fmt.Fprint(w, "<title>Internal Server Error</title><h1>Internal Server Error</h1>")
    }
}


func isLogin(w http.ResponseWriter, r *http.Request) bool {
    cookie, err := r.Cookie("id")
    if err != nil || cookie.Value == "" {
        return false
    } else {
        value := PassWord(SERVER_CONFIG.ServerKey, SERVER_CONFIG.PassWord)
        if value != cookie.Value {
            return false
        }
    }

    return true
}


func isRobots(w http.ResponseWriter, r *http.Request) bool {
    agent := strings.ToLower(getHTTPHeader(r, "User-Agent"))
    robots := []string{"bot", "spider", "archiver", "yahoo! slurp", "haosou"}
    for _, v := range robots {
        if strings.Contains(agent, v) {
            return true
        }
    }

    return false
}
