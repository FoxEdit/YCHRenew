package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/FoxEdit/YCHRenew/Models"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal("Pprof server error:", err)
		}
	}()

	// application
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Cannot open/create log file; console output only")
	} else {
		multiWriter := io.MultiWriter(os.Stdout, file)
		log.SetOutput(multiWriter)
		log.SetFlags(log.Ltime | log.Lshortfile)
		defer file.Close()
	}

	log.Print("======================================== Application start ========================================")
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("YCHRenew")

	// models
	linkModel := Models.NewLinkModel()
	authModel := Models.NewAuthModel()
	filterModel := Models.NewFilterModel()
	accountModel := Models.GetAccountModelInstance() // singleton model
	auctionModel := Models.NewAuctionModel()

	// viewmodels
	linkViewModel := ViewModels.NewLinkViewModel(linkModel)
	authViewModel := ViewModels.NewAuthViewModel(authModel)
	accountViewModel := ViewModels.NewAccountViewModel(accountModel)
	auctionViewModel := ViewModels.NewAuctionViewModel(auctionModel)
	filterViewModel := ViewModels.NewFilterViewModel(filterModel, accountModel)

	// setup viewmodels
	v := Views.NewMainWindow(
		mainWindow,
		fyne.Size{Width: 1347, Height: 540},
	)
	v.SetupViewModels(
		linkViewModel,
		authViewModel,
		filterViewModel,
		accountViewModel,
		auctionViewModel,
	)
	go v.SetUI()

	// setup callbacks to rerender whole UI
	authViewModel.SetUIRefreshCallback(v.SetUI)
	filterViewModel.SetUIRefreshCallback(v.SetUI)

	// main loop
	v.ShowAndRun()
}
