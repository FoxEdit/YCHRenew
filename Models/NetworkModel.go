package Models

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

var (
	instance      *WebClient
	onceNetClient sync.Once
)

type WebClient struct {
	client          http.Client
	ychCommishesURL *url.URL
	isAuthenticated bool
}

func getWebClientInstance() *WebClient {
	onceNetClient.Do(createWebClientInstance)
	return instance
}

func createWebClientInstance() {
	ychCommishesURL, _ := url.Parse(COMMISHES_URL)
	jar, _ := cookiejar.New(nil)
	instance = &WebClient{ychCommishesURL: ychCommishesURL, client: http.Client{Jar: jar}}

	log.Println("CLIENT CREATED")
}

func (n *WebClient) saveCookies() {
	fileModel := NewFileModel()
	cookiesArr := n.client.Jar.Cookies(n.ychCommishesURL)
	fileModel.WriteAuthCacheToStorage(cookiesArr)
}

func (n *WebClient) loadCookies() error {
	fileModel := NewFileModel()

	cookies := fileModel.ReadAuthCacheFromStorage()

	n.client.Jar.SetCookies(n.ychCommishesURL, cookies)
	n.isAuthenticated = true

	return nil
}

func (n *WebClient) PrintCookies() {
	for _, cookie := range n.client.Jar.Cookies(n.ychCommishesURL) {
		fmt.Printf("\n%s = %s", cookie.Name, cookie.Value)
	}
}

func (n *WebClient) Do(req *http.Request) (*http.Response, error) {
	response, err := n.client.Do(req)
	return response, err
}

func (n *WebClient) GetClient() *WebClient {
	return getWebClientInstance()
}
