package Models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

var (
	instance      *NetworkModel
	onceNetClient sync.Once
)

type NetworkModel struct {
	client          http.Client
	ychCommishesURL *url.URL
	isAuthenticated bool
}

func getWebClientInstance() *NetworkModel {
	onceNetClient.Do(createWebClientInstance)
	return instance
}

func createWebClientInstance() {
	ychCommishesURL, _ := url.Parse(COMMISHES_URL)
	jar, _ := cookiejar.New(nil)
	instance = &NetworkModel{ychCommishesURL: ychCommishesURL, client: http.Client{Jar: jar}}
	instance.isAuthenticated = false

	log.Println("CLIENT CREATED")
}

func (n *NetworkModel) saveCookies() {
	fileModel := NewFileModel()
	cookiesArr := n.client.Jar.Cookies(n.ychCommishesURL)
	fileModel.WriteAuthCacheToStorage(cookiesArr)
}

func (n *NetworkModel) PrintCookies() {
	for _, cookie := range n.client.Jar.Cookies(n.ychCommishesURL) {
		fmt.Printf("\n%s = %s", cookie.Name, cookie.Value)
	}
}

func (n *NetworkModel) Do(req *http.Request) (*http.Response, error) {
	response, err := n.client.Do(req)
	return response, err
}

func (n *NetworkModel) GetXSRFByPattern(endpoint string, xsrfBegin string, xsrfEnd string) string {
	req, _ := http.NewRequest("GET", endpoint, nil)
	res, _ := n.Do(req)
	resByte, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	resString := string(resByte)

	xsrfStartIndex := strings.Index(resString, xsrfBegin) + len(xsrfBegin)
	xsrfEndIndex := strings.Index(resString[xsrfStartIndex:], xsrfEnd)

	xsrf := resString[xsrfStartIndex : xsrfStartIndex+xsrfEndIndex]

	return xsrf
}

func (n *NetworkModel) GetClient() *NetworkModel {
	return getWebClientInstance()
}
