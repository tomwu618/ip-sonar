package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	log.Print("v0.0.1")

	http.HandleFunc("/", ExampleHandler)
	http.HandleFunc("/info", InfoHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	ip, _ := GetIP(r)

	log.Print(ip)

	w.Write([]byte(ip))
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]interface{})

	// 获取所有参数
	params := r.URL.Query()
	for k, v := range params {
		info[k] = v[0]
	}

	// 获取所有头信息
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = v[0]
	}
	info["headers"] = headers

	// 将信息转化为JSON格式
	jsonData, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头并返回JSON数据
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetIP(r *http.Request) (string, error) {
	ip := net.ParseIP(r.Header.Get("X-Real-IP"))
	if ip != nil {
		return ip.String(), nil
	}

	forwardIPs := r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(forwardIPs, ",") {
		ip = net.ParseIP(i)
		if ip != nil {
			return ip.String(), nil
		}
	}

	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	ip = net.ParseIP(remoteIP)
	if ip != nil {
		return ip.String(), nil
	}

	return "", errors.New("no valid ip found")
}
