package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	maxRoutine      = 256
	DefaultRoutines = 128
)

// {'success': True, 'msg': '登录成功', 'obj': None}
type XuiCheckResult struct {
	Success bool   `json:"success" xml:"success"`
	Msg     string `json:"msg" xml:"msg"`
	Obj     string `json:"obj" xml:"obj"`
}

var (
	loginData = map[string]string{
		"username": "admin",
		"password": "admin",
	}

	loginJSON []byte
)

func init() {
	var err error
	loginJSON, err = json.Marshal(loginData)
	if err != nil {
		log.Fatalf("Failed to marshal login data: %v", err)
	}
}

func ScanLinks(links []string) []string {
	// return success links
	res := []string{}
	total := len(links)

	controlChan := make(chan struct{}, DefaultRoutines)

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	sum := 0
	width := len(fmt.Sprint(total))

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			controlChan <- struct{}{}
			defer func() { <-controlChan }()
			ok := check_xui(link)
			mu.Lock()
			if ok {
				res = append(res, link)
			}
			sum++
			mu.Unlock()
			fmt.Printf("\r\033[KProgress scan xui: %*d/%*d | \033[32m%*d\033[0m", width, sum, width, total, width, len(res))
			// fmt.Printf("\r\033[KProgress scan xui: %*d/%*d: \033[32m%*d\033[0m", width, total, width, len(links), width, len(res))
		}(link)
	}

	wg.Wait()
	close(controlChan)
	fmt.Println()

	return res
}

func check_xui(url string) bool {

	target := url + "/login"

	// send login request
	hc := &http.Client{
		Timeout: 2 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodPost, target, bytes.NewReader(loginJSON))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := hc.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var xuiRes XuiCheckResult
	if err := json.NewDecoder(resp.Body).Decode(&xuiRes); err != nil {
		return false
	}

	if xuiRes.Success {
		return true
	}

	return false
}
