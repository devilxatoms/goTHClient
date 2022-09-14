package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
	writer  *multipart.Writer
	headers []Header
	params  []Param
}

var tokenFlag string
var reportPathFlag string
var serverFlag string
var scanTypeFlag string

func formData(reportPath string) (*multipart.Writer, *bytes.Buffer) {
	file, _ := os.Open(reportPath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", file.Name())
	io.Copy(part, file)
	// writer.Close()

	return writer, body
}

func (r *Request) AddHeader(name, value string) {
	r.headers = append(r.headers, Header{name: name, value: value})
}

func (r *Request) AddParam(name, value string) {
	r.params = append(r.params, Param{name: name, value: value})
}

func callApi(request *Request) (*http.Response, error) {
	client := &http.Client{}
	writer := request.writer

	// add params
	for _, param := range request.params {
		// log.Printf("Adding param %s with value %s", param.name, param.value)
		writer.WriteField(param.name, param.value)
	}

	err := writer.Close()
	if err != nil {
		return nil, err
		// log.Fatal("error in params: \n", err)
	}

	req, err := http.NewRequest(request.method, request.url, request.body)
	if err != nil {
		return nil, err
		// log.Fatal("error in req: \n", err)
	}

	// add headers
	for _, header := range request.headers {
		req.Header.Add(header.name, header.value)
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		// log.Fatal("send request error: \n", err)
		return nil, err
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

	return resp, nil
}

func GetUsers(server, token string) (*http.Response, error) {
	// construct the request
	getUsersRequest := NewRequest("GET", server+"/api/v2/users/", nil, nil, nil, nil)
	// add headers to the request
	getUsersRequest.AddHeader("Authorization", "Token "+token)
	getUsersRequest.AddHeader("Content-Type", "application/json")

	// call the api
	resp, err := callApi(getUsersRequest)
	return resp, err
}

func uploadReport(server, reportPath, token string) (*http.Response, error) {

	// get the file as a multipart body
	wr, dt := formData(reportPath)
	// construct the request
	uploadRequest := NewRequest("POST", server+"/api/v2/import-scan/", dt, wr, nil, nil)
	// add headers to the request
	uploadRequest.AddHeader("Authorization", "Token "+token)
	// add params to the request
	uploadRequest.AddParam("scan_date", "2022-05-09")
	uploadRequest.AddParam("minimum_severity", "Info")
	uploadRequest.AddParam("active", "true")
	uploadRequest.AddParam("verified", "true")
	uploadRequest.AddParam("scan_type", "Trufflehog3 Scan")
	uploadRequest.AddParam("product_name", "set-notes-ea")
	uploadRequest.AddParam("engagement_name", "test import set-notes-ea Dependency Check")
	uploadRequest.AddParam("environment", "Development")
	uploadRequest.AddParam("tags", "[\"Test\"]")

	// call the api
	resp, err := callApi(uploadRequest)
	return resp, err
}

func initFlags() {
	flag.StringVar(&tokenFlag, "t", "", "Defect Dojo API token")
	flag.StringVar(&reportPathFlag, "p", "defectDojo.json", "Defect Dojo Report Path")
	flag.StringVar(&serverFlag, "e", "", "DefectDojo Server Url")
	flag.StringVar(&scanTypeFlag, "s", "", "Defect Dojo Scan Type")
	flag.Parse()

}

func NewRequest(method, url string, body *bytes.Buffer, writer *multipart.Writer, headers []Header, params []Param) *Request {
	if body == nil {
		body = &bytes.Buffer{}
	}
	if writer == nil {
		writer = multipart.NewWriter(body)
	}

	return &Request{method: method, url: url, body: body, writer: writer, headers: headers, params: params}
}

func main() {
	initFlags()
	fmt.Println("Token is ", tokenFlag)
	fmt.Println("Path is ", reportPathFlag)
	fmt.Println("Server is ", serverFlag)
	fmt.Printf("Scan Type is %s \n\n", scanTypeFlag)

	GetUsers(serverFlag, tokenFlag)
	resp, err := uploadReport(serverFlag, reportPathFlag, tokenFlag)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Response: ", resp)
	}

}
