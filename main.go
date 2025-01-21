package main

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/FoxEdit/YCHRenew/Models"
	"github.com/FoxEdit/YCHRenew/Utility"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func main() {
	// application
	mainApp := app.New()
	Utility.GetClientInstance().LoadCookies()

	// models
	linkModel := Models.NewLinkModel()
	popupModel := Models.NewPopupModel()
	auctionTableModel := Models.NewAuctionTableModel()

	// viewmodels
	linkViewModel := ViewModels.NewLinkViewModel(linkModel)
	popupViewModel := ViewModels.NewPopupViewModel(popupModel)
	auctionTableViewModel := ViewModels.NewAuctionTableViewModel(auctionTableModel)

	// setup
	v := Views.NewMainWindow(
		mainApp.NewWindow("YCHRenew"),
		fyne.Size{Width: 850, Height: 500},
	)
	v.SetupViewModels(
		linkViewModel,
		popupViewModel,
		auctionTableViewModel,
	)
	v.SetUI()

	// main loop
	v.ShowAndRun()
}

func auth() {
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

	// Буфер для тела запроса
	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	defer multipartWriter.Close()

	multipartWriter.WriteField("_xsrf_", xsrf)
	multipartWriter.WriteField("time", "Sun, 19 Jan 2025 19:30:17 GMT")
	multipartWriter.WriteField("username", "CatEdit")
	multipartWriter.WriteField("password", "ZZxa098um7")

	loginPostReq, _ := http.NewRequest("POST", "https://account.commishes.com/user/login/", &body)
	loginPostReq.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client.Do(loginPostReq)
}
