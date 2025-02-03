package Views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views/CustomUITools"
	"image/color"
	"log"
	"strconv"
	"time"
)

type CardTable struct {
	parentWindow     *fyne.Window
	links            *ViewModels.LinkViewModel
	filter           *ViewModels.FilterViewModel
	accountViewModel *ViewModels.AccountViewModel
}

func NewCardTable(accountVM *ViewModels.AccountViewModel, links *ViewModels.LinkViewModel, filter *ViewModels.FilterViewModel, parentWindow *fyne.Window) *CardTable {
	return &CardTable{accountViewModel: accountVM, links: links, filter: filter, parentWindow: parentWindow}
}

func (a *CardTable) Build() fyne.CanvasObject {
	log.Println("CARD TABLE BUILD START")
	contentContainer := container.NewVBox()
	log.Println("LOADING CACHE")
	a.accountViewModel.LoadCachedAccountData()

	data := a.accountViewModel.GetAllRaw()
	if data == nil {
		log.Println("CACHE IS NULL, UPDATING CACHE")
		a.accountViewModel.UpdateDataFromAccount()
	}
	log.Println("CACHE UPDATED")

	cleanedCards := a.accountViewModel.GetAllCleaned()

	if len(cleanedCards) == 0 {
		return contentContainer // return without scroll background
	}

	for i := range cleanedCards {
		data := &cleanedCards[i]
		card := a.CreateCard(data)

		a.filter.AddNewCard(&ViewModels.CardItem{Data: data, Card: card})
		contentContainer.Add(card)
	}

	content := container.NewBorder(nil, nil, nil,
		CustomUITools.NewColorWSpacer(12, theme.Color(theme.ColorNameButton)), contentContainer)

	return container.NewVScroll(content)
}

func (a *CardTable) CreateCard(cleanedData *ViewModels.CleanAccountData) fyne.CanvasObject {
	card := container.NewHBox()
	card.Add(CustomUITools.NewWSpacer(5))

	img := a.getCardImage(cleanedData)
	card.Add(img)

	card.Add(CustomUITools.NewSeparator())

	title := a.getCardTitle(cleanedData)
	subtitle := a.getCardSubtitle(cleanedData)
	heat := a.getCardHeat(cleanedData)
	isOver := canvas.NewText(fmt.Sprintf("🕒 Идёт (%s)", cleanedData.EndsIn), color.RGBA{R: 0, G: 255, B: 0, A: 255})
	isOver.TextSize = 16
	if cleanedData.EndsUnix < time.Now().Unix() {
		isOver.Text = "🕒 Закончен"
		isOver.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}

	infoBlockWidth := canvas.NewRectangle(color.Transparent)
	infoBlockWidth.SetMinSize(fyne.NewSize(600, 0))
	infoBlock := container.NewStack(infoBlockWidth, container.NewPadded(container.NewVBox(subtitle, title, heat, isOver)))
	card.Add(infoBlock)

	card.Add(CustomUITools.NewSeparator())

	bid := a.getCardBid(cleanedData)
	bidBlockWidth := canvas.NewRectangle(color.Transparent)
	bidBlockWidth.SetMinSize(fyne.NewSize(150, 0))
	bidBlock := container.NewStack(bidBlockWidth, container.NewPadded(bid))
	card.Add(bidBlock)

	card.Add(CustomUITools.NewSeparator())

	auctionActionBtn := a.cardFunctionality(cleanedData)
	auctionActionWidth := canvas.NewRectangle(color.Transparent)
	auctionActionWidth.SetMinSize(fyne.NewSize(172.5, 0))
	auctionBlock := container.NewCenter(container.NewStack(auctionActionWidth, container.NewVBox(auctionActionBtn)))
	card.Add(auctionBlock)

	cardBackground := canvas.NewRectangle(color.RGBA{R: 32, G: 32, B: 35, A: 255})
	cardBackground.SetMinSize(fyne.NewSize(1100, 90))
	cardContainer := container.NewStack(cardBackground, card)

	return cardContainer
}

func (a *CardTable) getCardImage(cleanedData *ViewModels.CleanAccountData) fyne.CanvasObject {
	URI, _ := a.links.GetFyneURIFromString(cleanedData.ImgUrl)
	img := canvas.NewImageFromURI(URI)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(80, 80))

	c := container.NewStack(img)
	c.Add(img)

	return img
}

func (a *CardTable) getCardHeat(cleanedData *ViewModels.CleanAccountData) fyne.CanvasObject {
	rating := canvas.NewText("🔥 РЕЙТИНГ: "+strconv.FormatFloat(cleanedData.Heat*100, 'f', 2, 64)+"%", color.RGBA{R: 255, G: 165, B: 0, A: 255})
	rating.TextSize = 16
	return rating
}

func (a *CardTable) getCardTitle(cleanedData *ViewModels.CleanAccountData) fyne.CanvasObject {
	title := canvas.NewText(cleanedData.Title, color.White)
	title.TextSize = 14
	if cleanedData.Subtitle == "" {
		title.Text = "БЕЗ НАЗВАНИЯ"
		title.TextStyle = fyne.TextStyle{Italic: true}
	}

	return title
}

func (a *CardTable) getCardSubtitle(cleanedData *ViewModels.CleanAccountData) fyne.CanvasObject {
	subtitle := canvas.NewText(cleanedData.Subtitle, color.White)
	subtitle.TextSize = 12
	if cleanedData.Subtitle == "" {
		subtitle.Text = "БЕЗ ЗАГОЛОВКА"
		subtitle.TextStyle = fyne.TextStyle{Italic: true}
	}

	return subtitle
}

func (a *CardTable) getCardBid(cleanedData *ViewModels.CleanAccountData) fyne.CanvasObject {
	bidBox := container.NewVBox()

	bidBoxBackground := canvas.NewRectangle(color.Transparent)
	bidBoxBackground.SetMinSize(fyne.NewSize(200, 0))

	if cleanedData.Bid != nil {
		bid := canvas.NewText(fmt.Sprintf("Ставка: %d$", cleanedData.Bid.Bid), color.White)
		bidUserName := canvas.NewText(fmt.Sprintf("Юзер: %s", cleanedData.Bid.User.Name), color.White)

		bidBox.Add(bid)
		bidBox.Add(bidUserName)
	} else {
		bidBox.Add(canvas.NewText("Ставок нет.", color.White))
	}

	return container.NewStack(bidBoxBackground, bidBox)
}

func (a *CardTable) cardFunctionality(cleanedData *ViewModels.CleanAccountData) *widget.Button {
	var auctionActionBtn *widget.Button

	menuOptions := []*fyne.MenuItem{
		fyne.NewMenuItem("Рестарт без изменений", func() {
			dialog.NewInformation("Ошибка!", "Функция пока не реализована.", *a.parentWindow).Show()
		}),
		fyne.NewMenuItem("Рестарт c изменениями", func() {
			dialog.NewInformation("Ошибка!", "Функция пока не реализована.", *a.parentWindow).Show()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Посмотреть в браузере", func() {
			link := a.links.GetUrlFromRawString("https://ych.commishes.com" + cleanedData.CardUrl)
			err := fyne.CurrentApp().OpenURL(link)
			if err != nil {
				print("error when open card url: ", err)
			}
		}),
	}

	auctionActionBtn = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		popUp := widget.NewPopUpMenu(
			fyne.NewMenu("", menuOptions...),
			(*a.parentWindow).Canvas(),
		)

		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(auctionActionBtn)
		popUp.ShowAtPosition(pos.Add(fyne.NewPos(0, auctionActionBtn.Size().Height)))
	})

	return auctionActionBtn
}
