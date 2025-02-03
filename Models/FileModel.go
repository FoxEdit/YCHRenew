package Models

import (
	"bufio"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const appFolder = "YCHRenew"
const appAuthCache = "auth_cache"

const cookieFile = "data.dat"

type FileModel struct {
}

func NewFileModel() *FileModel {
	return new(FileModel)
}

func (f FileModel) WriteAuthCacheToStorage(cookies []*http.Cookie) {
	if !f.isCookieFileExists() {
		err := f.prepareCookieFolder()
		if err != nil {
			log.Println("COOKIE PREPARING FOLDER FAILED: ", err.Error()) // refactor to log + popup
			return
		}
	}

	var cookiesStr string
	for _, cookie := range cookies {
		cookiesStr += cookie.Name + "=" + cookie.Value + "\n"
	}

	appdata, _ := os.UserConfigDir()
	file, err := os.OpenFile(appdata+"\\"+appFolder+"\\"+appAuthCache+"\\"+cookieFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("COOKIE OPEN ERROR: ", err.Error()) // add popup
		return
	}
	defer file.Close()

	err = os.WriteFile(file.Name(), []byte(cookiesStr), 0644)
	if err != nil {
		log.Println("COOKIE WRITE ERROR: ", err.Error()) // add popup
		return
	}
}

func (f FileModel) isCookieFileExists() bool {
	appData, _ := os.UserConfigDir()

	appFolderPath := appData + "\\" + appFolder
	cookiePath := appFolderPath + "\\" + appAuthCache + "\\" + cookieFile

	_, fileErr := os.Stat(cookiePath)
	_, appFolderErr := os.Stat(appFolderPath)

	return fileErr == nil && appFolderErr == nil
}

func (f FileModel) prepareCookieFolder() error {
	if f.isCookieFileExists() {
		return nil
	}

	appData, _ := os.UserConfigDir()
	appFolderPath := appData + "\\" + appFolder
	cookiePath := appFolderPath + "\\" + appAuthCache + "\\" + cookieFile

	_, appFolderErr := os.Stat(appFolderPath)
	_, fileErr := os.Stat(cookiePath)

	if appFolderErr != nil {
		err := os.Mkdir(appFolderPath, 0777)
		if !errors.Is(err, fs.ErrExist) {
			log.Println("DIR ERROR: ", err.Error()) // refactor to log + popup
			return err
		}
	}

	if fileErr != nil {
		file, err := os.OpenFile(cookiePath, os.O_CREATE, 0644)
		defer file.Close()
		if err != nil {
			log.Println("FILE CREATE ERROR: ", err.Error()) // refactor to log + popup
			return err
		}
	}

	return nil
}

func (f FileModel) ReadAuthCacheFromStorage() []*http.Cookie {
	appdataPath, _ := os.UserConfigDir()
	appFolderPath := appdataPath + "\\" + appFolder + "\\" + appAuthCache
	cookiePath := appFolderPath + "\\" + cookieFile

	file, err := os.Open(cookiePath)
	defer file.Close()
	if err != nil {
		log.Println("error opening cookie file: ", err)
	}

	var cookiesArr []*http.Cookie
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cookie := strings.Split(line, "=")
		cookiesArr = append(cookiesArr, &http.Cookie{Name: cookie[0], Value: cookie[1]})
	}

	if err := scanner.Err(); err != nil {
		log.Println("error scanning cookie file: ", err)
	}

	if len(cookiesArr) == 0 {
		log.Println("empty cookie error")
	}

	return cookiesArr
}
