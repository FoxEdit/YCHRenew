package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/FoxEdit/YCHRenew/Models"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views"
)

func main() {
	// application
	mainApp := app.New()
	//Models.NewAuthModel().CookieLogin()
	//Models.NewCardModel().GetAllCardsFromAccount()

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
