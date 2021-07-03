package http

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"
)

const defaultContentType = "application/json"

var httpClient = &http.Client{}

type HeaderOption struct {
    Name  string
    Value string
}

type QueryParameter struct {
    Key   string
    Value interface{}
}

func init() {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    httpClient = &http.Client{
        Transport: tr,
        Timeout:   time.Duration(30) * time.Second,
    }
}

func BuildAuthorizationHeaderOptions(tenantAccessToken string) HeaderOption {
    return HeaderOption{
        Name:  "Authorization",
        Value: "Bearer " + tenantAccessToken,
    }
}

func PostRequest(url string, body string, headerOptions ...HeaderOption) (string, error) {
    req, err := http.NewRequest("POST", url, strings.NewReader(body))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", defaultContentType)
    for _, headerOption := range headerOptions {
        req.Header.Set(headerOption.Name, headerOption.Value)
    }
    resp, err := httpClient.Do(req)
    defer func() {
        if resp != nil {
            if e := resp.Body.Close(); e != nil {
                fmt.Println(e)
            }
        }
    }()
    return responseHandle(resp, err)
}

func Post(url string, params map[string]interface{}, body string, headerOptions ...HeaderOption) (string, error) {
    log.Printf("请求body %s\n", body)
    return PostRequest(url + convertToQueryParams(params), body, headerOptions...)
}

func PostRepetition(url string, params []QueryParameter, body string, headerOptions ...HeaderOption) (string, error) {
    log.Printf("请求body %s\n", body)
    return PostRequest(url + convertToQueryParamsMulti(params), body, headerOptions...)
}

func GetRequest(url string, headerOptions ...HeaderOption) (string, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }
    for _, headerOption := range headerOptions {
        req.Header.Set(headerOption.Name, headerOption.Value)
    }
    resp, err := httpClient.Do(req)
    defer func() {
        if resp != nil {
            if e := resp.Body.Close(); e != nil {
                log.Println(err)
            }
        }
    }()
    return responseHandle(resp, err)
}

func Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) (string, error) {
    return GetRequest(url + convertToQueryParams(params), headerOptions...)
}

func GetMulti(url string, params []QueryParameter, headerOptions ...HeaderOption) (string, error) {
    return GetRequest(url + convertToQueryParamsMulti(params), headerOptions...)
}

func responseHandle(resp *http.Response, err error) (string, error) {
    if err != nil {
        log.Println(err)
        return "", err
    }
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        return "", err
    }
    respBody := string(b)
    log.Printf("api %s 响应结果: %s\n", resp.Request.URL, respBody)
    return respBody, nil
}

func convertToQueryParams(params map[string]interface{}) string {
    if len(params) == 0 {
        return ""
    }
    var buffer bytes.Buffer
    buffer.WriteString("?")
    for k, v := range params {
        if v == nil {
            continue
        }
        buffer.WriteString(fmt.Sprintf("%s=%v&", k, url.QueryEscape(fmt.Sprintf("%v", v))))
    }
    buffer.Truncate(buffer.Len() - 1)
    return buffer.String()
}

func convertToQueryParamsMulti(params []QueryParameter) string {
    var buffer bytes.Buffer
    buffer.WriteString("?")
    for _, v := range params {
        if v.Value == nil {
            continue
        }
        buffer.WriteString(fmt.Sprintf("%s=%v&", v.Key, url.QueryEscape(fmt.Sprintf("%v", v.Value))))
    }
    buffer.Truncate(buffer.Len() - 1)
    return buffer.String()
}