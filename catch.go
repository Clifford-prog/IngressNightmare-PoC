package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 当前时间作为文件名（精确到秒）
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("request-%s.log", timestamp)

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Cannot create log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 打印请求信息
	fmt.Fprintln(file, "Method:", r.Method)
	fmt.Fprintln(file, "URL:", r.URL.String())
	fmt.Fprintln(file, "Headers:")
	for k, v := range r.Header {
		fmt.Fprintf(file, "  %s: %s\n", k, v)
	}

	// 打印请求体
	fmt.Fprintln(file, "\nBody:")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(file, "Error reading body:", err)
	} else {
		fmt.Fprintln(file, string(body))
	}

	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/networking/v1/ingresses", handler)

	// 使用自签名证书（你需要提供 cert.pem 和 key.pem）
	certFile := "/usr/local/certificates/cert"
	keyFile := "/usr/local/certificates/key"

	// TLS 配置（可选）
	server := &http.Server{
		Addr: ":7443",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	fmt.Println("Listening on https://0.0.0.0:7443")
	log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}
