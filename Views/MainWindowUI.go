package Views

import (
	"image/color"
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views/CustomUITools"
)

type MainWindow struct {
	window fyne.Window

	headerContent *fyne.Container
	filterContent *fyne.Container
	auctionContent *fyne.Container
	windowContent *fyne.Container

	accountViewModel     *ViewModels.AccountViewModel
	auctionViewModel     *ViewModels.AuctionViewModel
	linkViewModel        *ViewModels.LinkViewModel
	authViewModel        *ViewModels.AuthViewModel
	filterTableViewModel *ViewModels.FilterViewModel
}

func NewMainWindow(fyneWindow fyne.Window, size fyne.Size) *MainWindow {
	main := new(MainWindow)
	main.window = fyneWindow

	main.headerContent = container.NewStack()
	main.filterContent = container.NewStack()
	main.auctionContent = container.NewStack()
	main.windowContent = container.NewStack()

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

// Temporarily disabled
/* func (mw *MainWindow) animateLoading(isDone <-chan bool) {
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
} */

func (mw *MainWindow) ProcessNewUI() {
	mw.window.SetContent(canvas.NewRectangle(color.Transparent))
	mw.buildAndSetAllUI()
}

func(mw *MainWindow) RebuildAuctionContent() {
	mw.auctionContent.RemoveAll()
	mw.auctionContent.Add(NewCardTable(
		mw.accountViewModel, 
		mw.auctionViewModel, 
		mw.linkViewModel, 
		mw.filterTableViewModel, 
		&mw.window).Build())

	mw.auctionContent.Refresh()
	mw.windowContent.Refresh()
}

func(mw *MainWindow) RebuildHeader() {
	mw.headerContent.RemoveAll()
	mw.headerContent.Add(NewHeaderContent(mw.linkViewModel, mw.authViewModel, &mw.window).Build())
	mw.headerContent.Refresh()
	mw.windowContent.Refresh()
}

func(mw *MainWindow) RebuildFilter() {
	mw.filterContent.RemoveAll()
	mw.filterContent.Add(NewNavigationFilter(mw.filterTableViewModel).Build())
	mw.filterContent.Refresh()
	mw.windowContent.Refresh()
}

func (mw *MainWindow) buildAndSetAllUI() {
	log.Println("MAIN WINDOW BUILD STARTED")

	mw.headerContent.Add(NewHeaderContent(mw.linkViewModel, mw.authViewModel, &mw.window).Build())
	mw.filterContent.Add(NewNavigationFilter(mw.filterTableViewModel).Build())

	mw.auctionContent.Add(NewCardTable(
		mw.accountViewModel, 
		mw.auctionViewModel, 
		mw.linkViewModel, 
		mw.filterTableViewModel, 
		&mw.window).Build())

	mw.windowContent.Add(container.NewPadded(container.NewVBox(
		mw.headerContent,
		CustomUITools.NewHSpacer(2.5),
		CustomUITools.NewSeparator(),
		container.NewHBox(mw.filterContent, CustomUITools.NewSeparator(), container.NewVScroll(mw.auctionContent)),
	)))

	mw.window.SetContent(mw.windowContent)
}

func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}
