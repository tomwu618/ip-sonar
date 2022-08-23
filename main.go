package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	log.Print("v0.0.1")

	http.HandleFunc("/", ExampleHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	ip, _ := GetIP(r)

	log.Print(ip)

	w.Write([]byte(ip))
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
