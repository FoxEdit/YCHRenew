package Utility

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
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
	// TODO move this to main? create config logic?
	appDir     = "YCHRenew"
	cookieFile = "data.dat"
)

type NetClient struct {
	client          http.Client
	ychCommishesURL *url.URL
}

func GetClientInstance() *NetClient {
	once.Do(createInstance)
	return instance
}

func createInstance() {
	ychCommishesURL, _ := url.Parse("https://ych.commishes.com/")
	jar, _ := cookiejar.New(nil)
	instance = &NetClient{ychCommishesURL: ychCommishesURL, client: http.Client{Jar: jar}}

	fmt.Println("CLIENT CREATED")
}

func (n *NetClient) SaveCookies() {
	cookiesArr := n.client.Jar.Cookies(n.ychCommishesURL)

	// ----------------- TODO separate to another func
	appdata, _ := os.UserConfigDir()

	err := os.Mkdir(appdata+"\\"+appDir, 0777)
	if !errors.Is(err, fs.ErrExist) {
		println("DIR (TODO do POPUP here) Error:", err.Error())
		return
	}

	file, err := os.OpenFile(appdata+"\\"+appDir+"\\"+cookieFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		println("FILE (TODO do POPUP here) Error:", err.Error())
		return
	}
	// -------------------

	var cookiesString string
	for _, cookie := range cookiesArr {
		cookiesString += cookie.Name + "=" + cookie.Value + "\n"
	}

	os.WriteFile(file.Name(), []byte(cookiesString), 0644)
}

func (n *NetClient) LoadCookies() error {
	appdata, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config dir: %v", err)
	}

	// Проверяем существование папки приложения и файла
	cookiePath := appdata + "\\" + appDir + "\\" + cookieFile
	_, fileErr := os.Stat(cookiePath)
	_, appFolderErr := os.Stat(appdata + "\\" + appDir)

	if appFolderErr != nil || fileErr != nil {
		return fmt.Errorf("error detecting app/cookie folder: %v", appFolderErr)
	}

	// Открываем файл
	file, err := os.Open(cookiePath)
	if err != nil {
		return fmt.Errorf("error opening cookie file: %v", err)
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// Читаем файл и создаем cookies
	cookiesArr := []*http.Cookie{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cookie := strings.SplitN(line, "=", 2) // Безопасный Split с ограничением на 2 части
		if len(cookie) < 2 {
			continue // Пропускаем некорректные строки
		}
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

/*func (n *NetClient) createCookiesFile() error {
	TODO implement this
}*/

func (n *NetClient) PrintCookies() {
	for _, cookie := range n.client.Jar.Cookies(n.ychCommishesURL) {
		fmt.Printf("\n%s = %s", cookie.Name, cookie.Value)
	}
}

func (n *NetClient) Do(req *http.Request) (*http.Response, error) {
	response, err := n.client.Do(req)
	return response, err
}
