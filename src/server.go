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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/likexian/gokit/xfile"
	"github.com/likexian/gokit/xhash"
	"github.com/likexian/simplejson-go"
)

type ApiStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiResult struct {
	Status ApiStatus `json:"status"`
}

// HttpService start http service
func HttpService() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/passwd", passwdHandler)
	http.HandleFunc("/help", helpHandler)
	http.HandleFunc("/node", nodeHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/robots.txt", robotsTXTHandler)
	http.HandleFunc("/api/stat", apiStatHandler)
	http.HandleFunc("/api/node", apiNodeHandler)

	SERVER_LOGGER.Info("start http service")
	err := http.ListenAndServeTLS(":15944",
		SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSCert, SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSKey, nil)
	if err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
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
		"data":    status,
		"version": Version(),
	}
	_ = tpl.Execute(w, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		if xhash.Md5(SERVER_CONFIG.ServerKey, password).Hex() != SERVER_CONFIG.PassWord {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			value := xhash.Md5(SERVER_CONFIG.ServerKey, SERVER_CONFIG.PassWord).Hex()
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

		_ = tpl.Execute(w, map[string]string{"action": "login"})
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	expires := time.Now()
	expires = expires.AddDate(0, 0, -1)
	cookie := http.Cookie{Name: "id", Value: "", Expires: expires, HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func passwdHandler(w http.ResponseWriter, r *http.Request) {
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
			SERVER_CONFIG.PassWord = xhash.Md5(SERVER_CONFIG.ServerKey, password).Hex()
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

		_ = tpl.Execute(w, map[string]string{"action": "passwd"})
	}
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
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

	_ = tpl.Execute(w, map[string]string{"server": r.Host, "key": SERVER_CONFIG.ServerKey})
}

func nodeHandler(w http.ResponseWriter, r *http.Request) {
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

	_ = tpl.Execute(w, map[string]string{"server": r.Host, "key": SERVER_CONFIG.ServerKey, "version": Version()})
}

func robotsTXTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User-agent: *\r\nDisallow: /")
}

func apiNodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := ApiResult{
		Status: ApiStatus{},
	}

	defer func() {
		text, _ := simplejson.Dumps(result)
		_, _ = fmt.Fprint(w, text)
	}()

	if !isLogin(w, r) {
		result.Status.Code = 1
		result.Status.Message = "login timeout"
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result.Status.Code = 1
		result.Status.Message = "data error"
		return
	}

	data, err := simplejson.Loads(string(body))
	if err != nil {
		result.Status.Code = 1
		result.Status.Message = "data invalid"
		return
	}

	dataId, _ := data.Get("id").String()
	dataIdDir := SERVER_CONFIG.BaseDir + SERVER_CONFIG.DataDir + "/" + dataId[3:]
	if !xfile.Exists(dataIdDir) {
		result.Status.Code = 1
		result.Status.Message = "node id invalid"
		return
	}

	err = os.RemoveAll(dataIdDir)
	if err != nil {
		result.Status.Code = 1
		result.Status.Message = err.Error()
		return
	}
}

func apiStatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := ApiResult{
		Status: ApiStatus{},
	}

	defer func() {
		text, _ := simplejson.Dumps(result)
		_, _ = fmt.Fprint(w, text)
	}()

	ip := getHTTPHeader(r, "X-Real-Ip")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	clientKey := getHTTPHeader(r, "X-Client-Key")
	if clientKey == "" {
		result.Status.Code = 1
		result.Status.Message = "key empty"
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result.Status.Code = 1
		result.Status.Message = "body error"
		return
	}

	text := string(body)
	serverKey := xhash.Md5(SERVER_CONFIG.ServerKey, text).Hex()
	if serverKey != clientKey {
		result.Status.Code = 1
		result.Status.Message = "key invalid"
		return
	}

	data, err := simplejson.Loads(text)
	if err != nil {
		result.Status.Code = 1
		result.Status.Message = "body invalid"
		return
	}

	data.Set("ip", ip)
	name, _ := data.Get("host_name").String()
	data.Set("host_name", strings.Split(name, ".")[0])

	err = WriteStatus(SERVER_CONFIG.DataDir, data)
	if err != nil {
		result.Status.Code = 1
		result.Status.Message = err.Error()
		return
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	n := strings.LastIndex(r.URL.Path, ".")
	if n == -1 {
		httpError(w, r, http.StatusNotFound)
		return
	}

	ext := r.URL.Path[n+1:]
	if ext == "css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if ext == "js" {
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

func getHTTPHeader(r *http.Request, name string) string {
	if line, ok := r.Header[name]; ok {
		return line[0]
	}

	return ""
}

// httpError returns a http error
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

// isLogin returns request has login
func isLogin(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("id")
	if err != nil || cookie.Value == "" {
		return false
	}

	value := xhash.Md5(SERVER_CONFIG.ServerKey, SERVER_CONFIG.PassWord).Hex()
	return value == cookie.Value
}

// isRobots returns is a robot request
func isRobots(w http.ResponseWriter, r *http.Request) bool {
	agent := strings.ToLower(getHTTPHeader(r, "User-Agent"))
	robots := []string{"bot", "spider", "archiver", "crawler", "search"}
	for _, v := range robots {
		if strings.Contains(agent, v) {
			return true
		}
	}

	return false
}
