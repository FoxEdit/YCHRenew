package Models

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type AuthModel struct{}

func NewAuthModel() *AuthModel {
	return &AuthModel{}
}

func (am *AuthModel) CookieLogin() error {
	log.Println("COOKIE LOGIN")
	client := getWebClientInstance()
	if client.isAuthenticated {
		log.Println("ALREADY LOGGED IN")
		return errors.New("already logged in")
	}

	cookies := NewFileModel().ReadAuthCacheFromStorage()
	client.client.Jar.SetCookies(client.ychCommishesURL, cookies)

	if isLoggedIn() {
		client.isAuthenticated = true
		return nil
	}

	return errors.New("cookie authorization failure")
}

func (am *AuthModel) Login(login string, password string) error {
	log.Println("DEFAULT LOGIN: ")
	client := getWebClientInstance()
	if client.isAuthenticated {
		return errors.New("already logged in")
	}

	loginRequest, _ := http.NewRequest("GET", "https://account.commishes.com/user/login/", nil)
	loginResponse, _ := client.Do(loginRequest)
	loginResponseByte, _ := io.ReadAll(loginResponse.Body)
	loginResponseString := string(loginResponseByte)

	xsrfPattern := "name=\"_xsrf_\"   value=\""
	xsrfStart := strings.Index(loginResponseString, xsrfPattern) + len(xsrfPattern)
	xsrfEnd := strings.Index(loginResponseString[xsrfStart:], "\" />")
	xsrf := loginResponseString[xsrfStart : xsrfStart+xsrfEnd]
	defer loginResponse.Body.Close()

	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	defer multipartWriter.Close()

	multipartWriter.WriteField("_xsrf_", xsrf)
	multipartWriter.WriteField("time", "Sun, 19 Jan 2025 19:30:17 GMT")
	multipartWriter.WriteField("username", login)
	multipartWriter.WriteField("password", password)

	loginPostReq, _ := http.NewRequest("POST", "https://account.commishes.com/user/login/", &body)
	loginPostReq.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client.Do(loginPostReq)

	getCookiesRequest, _ := http.NewRequest("GET", "https://ych.commishes.com/account/", &body)
	client.Do(getCookiesRequest)

	if isLoggedIn() {
		client.saveCookies()
		client.isAuthenticated = true
		return nil
	}

	log.Println("DEFAULT LOGIN FAILRUE")
	return errors.New("default authorization failure")
}

func isLoggedIn() bool {
	req, _ := http.NewRequest("GET", COMMISHES_URL+"/account.json", nil)
	client := getWebClientInstance()
	res, _ := client.Do(req)

	finalURL := res.Request.URL.String()

	return finalURL == COMMISHES_URL+"/account.json"

	//return res.StatusCode == 200
}

func (am *AuthModel) Register(login string, password string, email string) {
	panic("implement me")
}
