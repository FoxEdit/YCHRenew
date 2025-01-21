package Models

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	instance *NetClient
	once     sync.Once
)

const appDir = "YCHRenew"
const cookieFile = "data.dat"

type NetClient struct {
	client          http.Client
	ychCommishesURL *url.URL
}

func getClientInstance() *NetClient {
	once.Do(createNetClientInstance)
	return instance
}

func createNetClientInstance() {
	ychCommishesURL, _ := url.Parse("https://ych.commishes.com/")
	jar, _ := cookiejar.New(nil)
	instance = &NetClient{ychCommishesURL: ychCommishesURL, client: http.Client{Jar: jar}}

	fmt.Println("CLIENT CREATED")
}

func (n *NetClient) saveCookies() {
	if !n.isCookieFileExists() {
		err := n.prepareCookieFolder()
		if err != nil {
			println("COOKIE_PREPARE_FOLDER FAILED", err.Error())
			return
		}
	}
	var cookiesString string
	cookiesArr := n.client.Jar.Cookies(n.ychCommishesURL)
	for _, cookie := range cookiesArr {
		cookiesString += cookie.Name + "=" + cookie.Value + "\n"
	}

	appdata, _ := os.UserConfigDir()
	file, err := os.OpenFile(appdata+"\\"+appDir+"\\"+cookieFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		println("COOKIE_OPEN (TODO do POPUP here) Error:", err.Error())
		return
	}

	err = os.WriteFile(file.Name(), []byte(cookiesString), 0644)
	if err != nil {
		println("COOKIE_WRITE (TODO do POPUP here) Error:", err.Error())
		return
	}
}

func (n *NetClient) isCookieFileExists() bool {
	appData, _ := os.UserConfigDir()
	appFolderPath := appData + "\\" + appDir
	cookiePath := appFolderPath + "\\" + cookieFile

	_, fileErr := os.Stat(cookiePath)
	_, appFolderErr := os.Stat(appFolderPath)

	return fileErr == nil && appFolderErr == nil
}

func (n *NetClient) prepareCookieFolder() error {
	appData, _ := os.UserConfigDir()
	appFolderPath := appData + "\\" + appDir
	cookiePath := appFolderPath + "\\" + cookieFile

	_, fileErr := os.Stat(cookiePath)
	_, appFolderErr := os.Stat(appFolderPath)

	if fileErr == nil && appFolderErr == nil {
		return nil
	}

	if appFolderErr != nil {
		err := os.Mkdir(appFolderPath, 0777)
		if !errors.Is(err, fs.ErrExist) {
			println("DIR (TODO do POPUP here) Error:", err.Error())
			return err
		}
	}

	if fileErr != nil {
		file, err := os.OpenFile(cookiePath, os.O_CREATE, 0644)
		defer file.Close()
		if err != nil {
			println("FILE (TODO do POPUP here) Error:", err.Error())
			return err
		}
	}

	return nil
}

func (n *NetClient) loadCookies() error {
	appdataPath, _ := os.UserConfigDir()
	appFolderPath := appdataPath + "\\" + appDir
	cookiePath := appFolderPath + "\\" + cookieFile

	file, err := os.Open(cookiePath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("error opening cookie file: %v", err)
	}

	// Читаем файл и создаем cookies
	var cookiesArr []*http.Cookie
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cookie := strings.Split(line, "=")
		cookiesArr = append(cookiesArr, &http.Cookie{Name: cookie[0], Value: cookie[1]})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning cookie file: %v", err)
	}

	if len(cookiesArr) == 0 {
		return errors.New("empty cookie error")
	}

	n.client.Jar.SetCookies(n.ychCommishesURL, cookiesArr)

	return nil
}

func (n *NetClient) PrintCookies() {
	for _, cookie := range n.client.Jar.Cookies(n.ychCommishesURL) {
		fmt.Printf("\n%s = %s", cookie.Name, cookie.Value)
	}
}

func (n *NetClient) Do(req *http.Request) (*http.Response, error) {
	response, err := n.client.Do(req)
	return response, err
}

//
// =================================================================================================================
//

type AuthModel struct{}

func NewAuthModel() *AuthModel {
	return &AuthModel{}
}

func (am *AuthModel) GetAuthorizedClient() *NetClient {
	return getClientInstance()
}

func (am *AuthModel) CookieLogin() {
	getClientInstance().loadCookies()
}

func (am *AuthModel) Login(login string, password string) {
	client := getClientInstance()

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
	multipartWriter.WriteField("username", os.Args[0]) // debug ver, actual - login from params
	multipartWriter.WriteField("password", os.Args[1]) // debug ver, actual - password from params

	loginPostReq, _ := http.NewRequest("POST", "https://account.commishes.com/user/login/", &body)
	loginPostReq.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client.Do(loginPostReq)

	getCookiesRequest, _ := http.NewRequest("GET", "https://ych.commishes.com/account/", &body)
	client.Do(getCookiesRequest)
	client.saveCookies()
}

func (am *AuthModel) Register(login string, password string, email string) {
	panic("implement me")
}
