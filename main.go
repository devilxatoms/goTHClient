package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Header struct {
	name  string
	value string
}

type Param struct {
	name  string
	value string
}

type Request struct {
	method  string
	url     string
	body    *bytes.Buffer
	headers []Header
	params  []Param
}

var tokenFlag string
var reportPathFlag string
var serverFlag string
var scanTypeFlag string

func buildMultipartBody(reportPath, paramName string) *bytes.Buffer {
	file, err := os.Open(reportPath)
	if err != nil {
		log.Fatal(err)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		log.Fatal(err)
	}
	part.Write(fileContents)

	writer.Close()

	if err != nil {
		log.Fatal(err)
	}
	return body
}

func (r *Request) AddHeader(name, value string) {
	r.headers = append(r.headers, Header{name: name, value: value})
}

func (r *Request) AddParam(name, value string) {
	r.params = append(r.params, Param{name: name, value: value})
}

func callApi(request *Request) {
	log.Printf("1 Current Lenght: %v", strconv.FormatInt(int64(request.body.Len()), 10))
	client := &http.Client{}
	req, err := http.NewRequest(request.method, request.url, request.body)
	if err != nil {
		log.Fatal(err)
	}

	// add headers
	for _, header := range request.headers {
		req.Header.Add(header.name, header.value)
	}

	log.Printf("2 Current Lenght: %v", strconv.FormatInt(int64(request.body.Len()), 10))

	writer := multipart.NewWriter(request.body)
	log.Printf("3 Current Lenght: %v", strconv.FormatInt(int64(request.body.Len()), 10))

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	// add params
	for _, param := range request.params {
		_ = writer.WriteField(param.name, param.value)
	}
	log.Printf("4 Current Lenght: %v", strconv.FormatInt(int64(request.body.Len()), 10))

	req.Header.Add("Content-Length", strconv.FormatInt(int64(request.body.Len()), 10))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("5 Current Lenght: %v", strconv.FormatInt(int64(request.body.Len()), 10))
		// log.Printf("Request failed: %v", req)
		log.Fatal("send request error: \n", err)
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func GetUsers(server, token string) *Request {
	// construct the request
	getUsersRequest := NewRequest("GET", server+"/api/v2/users/", nil, nil, nil)
	// add headers to the request
	getUsersRequest.AddHeader("Authorization", "Token "+token)
	// call the api
	callApi(getUsersRequest)
	return getUsersRequest
}

func uploadReport(server, reportPath, token string) *Request {
	// get the file as a multipart body
	bodyData := buildMultipartBody(reportPath, "file")
	// construct the request
	uploadRequest := NewRequest("POST", server+"/api/v2/import-scan/", bodyData, nil, nil)
	// add headers to the request
	uploadRequest.AddHeader("Authorization", "Token "+token)
	// uploadRequest.AddHeader("Content-Type", "application/json")
	// add params to the request
	uploadRequest.AddParam("scan_type", "Trufflehog Scan")
	uploadRequest.AddParam("minimum_severity", "Info")
	uploadRequest.AddParam("active", "true")
	uploadRequest.AddParam("verified", "true")
	uploadRequest.AddParam("push_to_jira", "false")
	uploadRequest.AddParam("close_old_findings", "false")
	uploadRequest.AddParam("engagement", "1")
	uploadRequest.AddParam("lead", "1")
	uploadRequest.AddParam("tags", "['test']")
	uploadRequest.AddParam("scan_date", "2022-05-09")

	log.Printf("Request type: %T \n", uploadRequest)
	log.Printf("the upload request is a: %v \n", uploadRequest)

	// call the api
	callApi(uploadRequest)
	return uploadRequest
}

func initFlags() {
	flag.StringVar(&tokenFlag, "t", "", "Defect Dojo API token")
	flag.StringVar(&reportPathFlag, "p", "defectDojo.json", "Defect Dojo Report Path")
	flag.StringVar(&serverFlag, "h", "", "DefectDojo Server Url")
	flag.StringVar(&scanTypeFlag, "s", "", "Defect Dojo Scan Type")
	flag.Parse()
}

func NewRequest(method, url string, body *bytes.Buffer, headers []Header, params []Param) *Request {
	if body == nil {
		body = new(bytes.Buffer)
	}

	return &Request{method: method, url: url, body: body, headers: headers, params: params}
}

func main() {
	initFlags()
	fmt.Println("Token is ", tokenFlag)
	fmt.Println("Path is ", reportPathFlag)
	fmt.Println("Server is ", serverFlag)
	fmt.Printf("Scan Type is %s \n\n", scanTypeFlag)

	GetUsers(serverFlag, tokenFlag)
	uploadReport(serverFlag, reportPathFlag, tokenFlag)

}
