package Views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/FoxEdit/YCHRenew/ViewModels"
)

type MainWindow struct {
	window fyne.Window

	linkViewModel         *ViewModels.LinkViewModel
	popupViewModel        *ViewModels.PopupViewModel
	auctionTableViewModel *ViewModels.AuctionTableViewModel
}

func NewMainWindow(fyneWindow fyne.Window, size fyne.Size) *MainWindow {
	main := new(MainWindow)
	main.window = fyneWindow

	main.window.SetMaster()
	main.window.Resize(size)
	main.window.SetFixedSize(true)
	main.window.CenterOnScreen()

	return main
}

func (mw *MainWindow) SetupViewModels(
	linkViewModel *ViewModels.LinkViewModel,
	popupViewModel *ViewModels.PopupViewModel,
	auctionTableViewModel *ViewModels.AuctionTableViewModel) {
	mw.linkViewModel = linkViewModel
	mw.popupViewModel = popupViewModel
	mw.auctionTableViewModel = auctionTableViewModel
}

func (mw *MainWindow) SetUI() {
	mw.window.SetContent(mw.buildUI())
}

func (mw *MainWindow) buildUI() fyne.CanvasObject {
	header := NewHeaderContent(mw.linkViewModel, mw.popupViewModel).Build()

	navFilter := NewNavigationFilter(mw.auctionTableViewModel).Build()
	auctionTable := NewAuctionTable().Build()

	content := container.NewHBox(navFilter, NewSeparator(), auctionTable) // ADD: navigation filter + ychTable

	return container.NewPadded(container.NewVBox(
		header,
		layout.NewSpacer(),
		NewSeparator(),
		content,
	))
}

func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}
