package http_mng

import (
	service "dynamic-dirb/internal/services"
	"dynamic-dirb/internal/string_mng"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type ResponseInfo struct {
	Url                string
	Domain             string
	Response           string
	ResponseCode       int
	ResponseBodyString string
	ResponseHeader     map[string][]string
}

func RequestGet(url string) ResponseInfo {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Add("If-None-Match", `W/"wyzzy"`)
	if req != nil {
		var headerParam = service.GetParameters().GetHeaders()
		if headerParam != "" {
			headerSplit := strings.Split(headerParam, ",")
			for _, singleHeader := range headerSplit {
				headerDetail := strings.Split(singleHeader, ":")
				req.Header.Set(headerDetail[0], headerDetail[1])
			}
		}
	} else {
		string_mng.PrintError("[-] Nil request client generated, no problem, keep going")
	}

	if err != nil {
		service.GetParameters().PrintDebug(err.Error())
		string_mng.PrintError("[-] Error on url: " + url + ". Error details: " + err.Error())
		var response ResponseInfo
		response.ResponseCode = 0
		response.Domain = ""
		response.Url = ""
		response.ResponseHeader = nil
		response.ResponseBodyString = ""
		return response
	}

	resp, err := client.Do(req)
	// if error occurred set all struct variable to empty
	if err != nil {
		service.GetParameters().PrintDebug(err.Error())
		string_mng.PrintError("[-] Error on url: " + url + ". Error details: " + err.Error())
		var response ResponseInfo
		response.ResponseCode = 0
		response.Domain = ""
		response.Url = ""
		response.ResponseHeader = nil
		response.ResponseBodyString = ""
		return response

	}
	var response ResponseInfo
	response.ResponseCode = resp.StatusCode
	response.Domain = GetDomainName(url)
	response.Url = url
	response.ResponseHeader = resp.Header
	data, _ := ioutil.ReadAll(resp.Body)
	response.ResponseBodyString = string(data)
	resp.Body.Close()

	return response
}

//url(string): url to test
//return (ResponseInfo,bool): true if response code = 200 else false (both case return also response)
func RequestResourceExist(url string) (ResponseInfo, bool) {
	response := RequestGet(url)
	if response.ResponseBodyString != "" {
		if service.IsResourceExist(service.GetParameters().GetResourceExistMapCode(), response.ResponseCode) {
			return response, true
		}
	}
	return response, false
}

// url(string): url string
// result(string): return only the domain name starting from a url
func GetDomainName(url string) string {
	//https://admin:admin@www.test.com
	var regex = regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+)`)
	var result = regex.FindAllStringSubmatch(url, -1)
	//[[https://example.com example.com]]
	if len(result) == 0 {
		return ""
	}
	return result[0][1]
}
