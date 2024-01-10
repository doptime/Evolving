package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func CurlGet(url string) ([]byte, error) {
	var resp *http.Response
	var err error
	if resp, err = http.Get(url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
func CurlGetWithTimeout(url string, timeoutSeconds int) ([]byte, error) {
	var resp *http.Response
	var err error
	var httpClient = &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	if resp, err = httpClient.Get(url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func FakeAIPost(FakeOpenAIUrl string, AccessToken string, XYHelperAuthKey string, post any) (messege string, err error) {
	var (
		resp                *http.Response
		postBytes, revBytes []byte
		req                 *http.Request
		client              *http.Client = &http.Client{Timeout: 360 * time.Second}
	)
	// if strings.Contains(url, "fakeopen.com") {
	// 	today := time.Now().Format("20060102")

	// 	// Generate the string
	// 	url = fmt.Sprintf("ai-%s.fakeopen.com", today)
	// }
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if postBytes, err = json.Marshal(post); err != nil {
		return "", err
	}
	if req, err = http.NewRequest("POST", FakeOpenAIUrl, bytes.NewBuffer(postBytes)); err != nil {
		return messege, err
	}
	req.Header.Add("User-Agent", getUseragent())
	if len(AccessToken) > 0 {
		req.Header.Add("Authorization", "Bearer "+AccessToken)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "text/event-stream; charset=utf-8")
	if len(XYHelperAuthKey) > 0 {
		req.Header.Add("authkey", XYHelperAuthKey)
	}
	//add cookie _puid with value puid
	//req.AddCookie(&http.Cookie{Name: "_puid", Value: puid[AccessTokenIndex%len(puid)]})

	if resp, err = client.Do(req); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	revBytes, err = io.ReadAll(resp.Body)
	return string(revBytes), err
}

var UAs []string = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) ",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Mobile Safari/537.36"}

func getUseragent() string {
	//with weight	0.5, 0.3, 0.2
	randomWeight := rand.Float64()
	if randomWeight < 0.5 {
		return UAs[0]
	} else if randomWeight < 0.8 {
		return UAs[1]
	} else {
		return UAs[2]
	}
}
func FakeAIGet(Url string, AccessToken string, post any) (messege string, err error) {
	var (
		resp                *http.Response
		postBytes, revBytes []byte
		req                 *http.Request
		client              *http.Client = &http.Client{Timeout: 420 * time.Second}
	)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if postBytes, err = json.Marshal(post); err != nil {
		return "", err
	}
	// http://chatgpt.lan:8010/api/conversation/gen_title/53fd1db3-e328-48bc-bd3e-0913daf49f05
	if req, err = http.NewRequest("GET", Url, bytes.NewBuffer(postBytes)); err != nil {
		return messege, err
	}
	req.Header.Add("User-Agent", getUseragent())
	req.Header.Add("Authorization", "Bearer "+AccessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "text/event-stream; charset=utf-8")
	//add cookie _puid with value puid
	//req.AddCookie(&http.Cookie{Name: "_puid", Value: puid[AccessTokenIndex%len(puid)]})

	if resp, err = client.Do(req); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	revBytes, err = io.ReadAll(resp.Body)
	return string(revBytes), err
}
func FakeAIDel(Url string, AccessToken string) (messege string, err error) {
	type In struct {
		Is_visible bool `json:"is_visible"`
	}
	var (
		resp                *http.Response
		bytes               []byte
		visible             = In{Is_visible: false}
		visibleJsonBytes, _ = json.Marshal(visible)
		visibleJsonString   = string(visibleJsonBytes)
	)
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", Url, strings.NewReader(visibleJsonString))
	req.Header.Add("User-Agent", getUseragent())
	req.Header.Add("Authorization", "Bearer "+AccessToken)
	req.Header.Add("Content-Type", "application/json")

	if resp, err = client.Do(req); err != nil {
		return messege, err
	}
	defer resp.Body.Close()
	bytes, err = io.ReadAll(resp.Body)
	return string(bytes), err
}
func HttpPostStructured(url string, in interface{}, out interface{}) (err error) {
	var (
		jsonBytes, ret []byte
		resp           *http.Response
	)
	//create http post request,
	if jsonBytes, err = json.Marshal(in); err != nil {
		return err
	}
	var httpClient = &http.Client{Timeout: 240 * time.Second}
	if resp, err = httpClient.Post(url, "application/json", strings.NewReader(string(jsonBytes))); err != nil {
		return err
	}
	defer resp.Body.Close()
	if ret, err = io.ReadAll(resp.Body); err != nil {
		return err
	}
	if err = json.Unmarshal(ret, &out); err != nil {
		return err
	}

	return
}

func CurlPostFormData(url string, postData url.Values) ([]byte, error) {
	var resp *http.Response
	var err error
	var httpClient = &http.Client{Timeout: 240 * time.Second}
	if resp, err = httpClient.PostForm(url, postData); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
func CurlPost(url string, postJson []byte) ([]byte, error) {
	var resp *http.Response
	var err error
	var httpClient = &http.Client{Timeout: 240 * time.Second}
	if resp, err = httpClient.Post(url, "application/json", strings.NewReader(string(postJson))); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
func CurlDel(url string) (err error) {
	var (
		req    *http.Request
		client *http.Client
		resp   *http.Response
	)

	if req, err = http.NewRequest(http.MethodDelete, url, nil); err != nil {
		return err
	}

	client = &http.Client{}
	if resp, err = client.Do(req); err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	return nil
}

var regexRemoveJsonPadding = regexp.MustCompile(`"|\s`)

func TrimJson(in string) (out string) {
	return regexRemoveJsonPadding.ReplaceAllString(in, "")
}
