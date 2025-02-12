package Views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views/CustomUITools"
	"image/color"
	"log"
	"time"
)

type MainWindow struct {
	window fyne.Window

	accountViewModel     *ViewModels.AccountViewModel
	auctionViewModel     *ViewModels.AuctionViewModel
	linkViewModel        *ViewModels.LinkViewModel
	authViewModel        *ViewModels.AuthViewModel
	filterTableViewModel *ViewModels.FilterViewModel
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
	authViewModel *ViewModels.AuthViewModel,
	auctionTableViewModel *ViewModels.FilterViewModel,
	accountViewModel *ViewModels.AccountViewModel,
	auctionViewModel *ViewModels.AuctionViewModel) {
	mw.auctionViewModel = auctionViewModel
	mw.accountViewModel = accountViewModel
	mw.linkViewModel = linkViewModel
	mw.filterTableViewModel = auctionTableViewModel
	mw.authViewModel = authViewModel
}

func (mw *MainWindow) animateLoading(isDone <-chan bool) {
	loadingText := canvas.NewText("Loading", color.White)
	loadingText.TextSize = 24
	loadingContainer := container.NewCenter(loadingText)

	for i := 1; ; i++ {
		select {
		case stop := <-isDone:
			if stop {
				return
			}
		default:
			time.Sleep(500 * time.Millisecond)
			if i%4 == 0 {
				loadingText.Text = "Loading"
			} else {
				loadingText.Text += "."
			}
			mw.window.SetContent(loadingContainer)
		}
	}
}

func (mw *MainWindow) SetUI() {
	// reset and build new
	mw.window.SetContent(canvas.NewRectangle(color.Transparent))
	isDone := make(chan bool)
	go mw.animateLoading(isDone)
	mw.window.SetContent(mw.buildUI(isDone))
}

func (mw *MainWindow) RefreshHeader() {
}

func (mw *MainWindow) RefreshCardTable() {

}

func (mw *MainWindow) buildUI(isDone chan<- bool) fyne.CanvasObject {
	log.Println("MAIN WINDOW BUILD STARTED")
	log.Println("BUILDING HEADER")
	header := NewHeaderContent(mw.linkViewModel, mw.authViewModel, &mw.window).Build()

	log.Println("BUILDING NAVIGATION FILTER")
	filter := NewNavigationFilter(mw.filterTableViewModel).Build()

	log.Println("BUILDING AUCTION TABLE")
	cardTable := NewCardTable(mw.accountViewModel, mw.auctionViewModel, mw.linkViewModel, mw.filterTableViewModel, &mw.window).Build()

	content := container.NewHBox(filter, CustomUITools.NewSeparator(), cardTable)

	isDone <- true
	return container.NewPadded(container.NewVBox(
		header,
		CustomUITools.NewHSpacer(2.5),
		CustomUITools.NewSeparator(),
		content,
	))
}

func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}
