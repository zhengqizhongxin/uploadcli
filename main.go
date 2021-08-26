package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var VERSION string = "upload CLI 1.0 build 20210827\r\n"

type inputArgs struct {
	version bool
	upUri   string
	inFile  string
}

var InputArgs inputArgs

func init() {
	flag.StringVar(&InputArgs.upUri, "u", "http://127.0.0.1:8080/up", "upload file url")
	flag.StringVar(&InputArgs.inFile, "f", "demo.txt", "file path")
	flag.BoolVar(&InputArgs.version, "v", false, "show version and exit")
}

func postFile(filename string, targetUrl string) (*http.Response, error) {
	postBodyBuf := bytes.NewBufferString("")
	postBodyWriter := multipart.NewWriter(postBodyBuf)

	_, err := postBodyWriter.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}

	boundary := postBodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	requestReader := io.MultiReader(postBodyBuf, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", targetUrl, requestReader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(postBodyBuf.Len()) + int64(closeBuf.Len())

	return http.DefaultClient.Do(req)
}

func main() {
	flag.Parse()
	if InputArgs.version {
		fmt.Printf("%s", VERSION)
		return
	}
	url := InputArgs.upUri
	file := InputArgs.inFile

	rsp, err := postFile(file, url)
	if err != nil {
		log.Println("Upload failed")
		return
	}

	log.Println("rsp:", rsp.Status)

}
