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

// TODO: add popups when error occured
type FileModel struct {
}

func NewFileModel() *FileModel {
	return &FileModel{}
}

func (f FileModel) WriteAuthCacheToStorage(cookies []*http.Cookie) error {
	if !f.isCookieFileExists() {
		err := f.prepareCookieFolder()
		if err != nil {
			log.Println("COOKIE PREPARING FOLDER FAILED: ", err.Error()) // add popup
			return errors.New("cookie preparing folder failed")
		}
	}

	var cookiesStr string
	for _, cookie := range cookies {
		cookiesStr += cookie.Name + "=" + cookie.Value + "\n"
	}

	appdata, _ := os.UserConfigDir()
	file, err := os.OpenFile(appdata+"\\"+appFolder+"\\"+appAuthCache+"\\"+cookieFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("COOKIE OPEN ERROR: ", err.Error())
		return errors.New("cookie open error")
	}
	defer file.Close()

	err = os.WriteFile(file.Name(), []byte(cookiesStr), 0644)
	if err != nil {
		log.Println("COOKIE WRITE ERROR: ", err.Error())
		return errors.New("cookie write error")
	}

	return nil
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
	cookieFolderPath := appFolderPath + "\\" + appAuthCache
	cookiePath := cookieFolderPath + "\\" + cookieFile

	_, appFolderErr := os.Stat(appFolderPath)
	_, cookieFolderErr := os.Stat(cookieFolderPath)
	_, cookieFileErr := os.Stat(cookiePath)

	if appFolderErr != nil {
		err := os.Mkdir(appFolderPath, 0777)
		if !errors.Is(err, fs.ErrExist) {
			log.Println("DIR ERROR: ", err.Error())
			return err
		}
	}

	if cookieFolderErr != nil {
		err := os.Mkdir(appFolderPath, 0777)
		if !errors.Is(err, fs.ErrExist) {
			log.Println("COOKIE DIR ERROR: ", err.Error())
			return err
		}
	}

	if cookieFileErr != nil {
		file, err := os.OpenFile(cookiePath, os.O_CREATE, 0644)
		defer file.Close()
		if err != nil {
			log.Println("FILE CREATE ERROR: ", err.Error())
			return err
		}
	}
	return nil
}

func (f FileModel) ReadAuthCacheFromStorage() ([]*http.Cookie, error) {
	appdataPath, _ := os.UserConfigDir()
	appFolderPath := appdataPath + "\\" + appFolder + "\\" + appAuthCache
	cookiePath := appFolderPath + "\\" + cookieFile

	file, err := os.Open(cookiePath)
	if err != nil {
		log.Println("error opening cookie file: ", err)
		return nil, err
	}
	defer file.Close()

	var cookiesArr []*http.Cookie
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cookie := strings.Split(line, "=")
		cookiesArr = append(cookiesArr, &http.Cookie{Name: cookie[0], Value: cookie[1]})
	}

	if err := scanner.Err(); err != nil {
		log.Println("error scanning cookie file: ", err)
		return nil, err
	}

	if len(cookiesArr) == 0 {
		log.Println("empty cookie error")
		return nil, err
	}

	return cookiesArr, nil
}
