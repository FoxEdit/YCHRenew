package Models

import (
	"bytes"
	"github.com/FoxEdit/YCHRenew/Utility"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type AuthModel struct {
	login    string
	password string
	email    string
}

func NewAuthModel(login string, password string, email string) *AuthModel {
	return &AuthModel{login, password, email}
}

func (am *AuthModel) Login() {
	client := Utility.GetClientInstance()

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
	multipartWriter.WriteField("username", "")
	multipartWriter.WriteField("password", "")

	loginPostReq, _ := http.NewRequest("POST", "https://account.commishes.com/user/login/", &body)
	loginPostReq.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client.Do(loginPostReq)

	getCookiesRequest, _ := http.NewRequest("GET", "https://ych.commishes.com/account/", &body)
	client.Do(getCookiesRequest)
	client.SaveCookies()
}

func (am *AuthModel) Register() {
	panic("implement me")
}
