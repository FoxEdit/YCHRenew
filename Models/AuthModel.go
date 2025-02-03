package Models

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type AuthModel struct {
	login    string
	password string
}

func NewAuthModel() *AuthModel {
	return &AuthModel{}
}

func (am *AuthModel) CookieLogin() error {
	log.Println("COOKIE LOGIN")
	return getWebClientInstance().loadCookies()
}

func (am *AuthModel) Login(login string, password string) error {
	log.Println("DEFAULT LOGIN")
	client := getWebClientInstance()

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

	client.saveCookies()
	client.isAuthenticated = true

	return nil // todo check if login success or not
}

func (am *AuthModel) Register(login string, password string, email string) {
	panic("implement me")
}
