package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var filePath string
	var serverIP string
	var serverPort string

	flag.StringVar(&filePath, "u", "", "Path to the file to upload")
	flag.StringVar(&serverIP, "i", "localhost", "Server IP address")
	flag.StringVar(&serverPort, "p", "8080", "Server port")
	flag.Parse()

	if filePath == "" {
		fmt.Println("Usage: upload.exe -u <file-path> [-i <server-ip>] [-p <server-port>]")
		return
	}

	serverURL := fmt.Sprintf("http://%s:%s/upload", serverIP, serverPort)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Printf("Failed to create form file: %v\n", err)
		return
	}

	if _, err := io.Copy(part, file); err != nil {
		fmt.Printf("Failed to copy file content: %v\n", err)
		return
	}

	writer.Close()

	resp, err := http.Post(serverURL, writer.FormDataContentType(), body)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Server response: %s\n", respBody)
}
