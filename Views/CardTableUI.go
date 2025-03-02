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

var imageCache = make(map[string]*canvas.Image)

type CardTable struct {
	parentWindow     *fyne.Window
	links            *ViewModels.LinkViewModel
	filter           *ViewModels.FilterViewModel
	accountVM *ViewModels.AccountViewModel
	auctionVM *ViewModels.AuctionViewModel
}

func NewCardTable(accountVM *ViewModels.AccountViewModel, auctionVM *ViewModels.AuctionViewModel, links *ViewModels.LinkViewModel, filter *ViewModels.FilterViewModel, parentWindow *fyne.Window) *CardTable {
	cardTable := CardTable{accountVM: accountVM, auctionVM: auctionVM, links: links, filter: filter, parentWindow: parentWindow}
	return &cardTable
}

func (a *CardTable) Build() *fyne.Container {
	log.Println("STARTED BUILDING CARD TABLE VIEW")
	contentContainer := container.NewVBox()
	
	log.Println("LOADING CACHE")
	a.accountVM.LoadCachedAccountData()

	if a.accountVM.GetAllRaw() == nil {
		log.Println("CACHE IS NULL, UPDATING CACHE")
		a.accountVM.UpdateDataFromAccount()
	}
	log.Println("CACHE UPDATED")

	cleanedCards := a.accountVM.GetAllCleaned()

	if cleanedCards == nil {
		return contentContainer // return without scroll background
	}

	for _, data := range cleanedCards {
		card := a.CreateCard(&data)

		//a.filter.AddNewCard(&ViewModels.CardItem{Data: data, Card: card})
		contentContainer.Add(card)
	}

	content := container.NewBorder(
		nil,
		nil,
		nil,
		CustomUITools.NewColorWSpacer(12, theme.Color(theme.ColorNameButton)), contentContainer)

	return content
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

	if val, ok := imageCache[cleanedData.ImgUrl]; ok {
		log.Println("IMAGE CACHE HIT")
		return val
	}
	
	log.Println("IMAGE CACHE MISS")
	img := canvas.NewImageFromURI(URI)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(80, 80))

	imageCache[cleanedData.ImgUrl] = img

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
			var selectDialog *dialog.CustomDialog

			var auctionCategory string
			auctionCategorySelect := widget.NewSelect([]string{
				"Фурри",
				"Люди",
				"Адопты",
				"Пони",
				"Самоделки",
			}, func(s string) {
				auctionCategory = s
			})
			auctionCategorySelect.PlaceHolder = "-- Категория --"

			var auctionTime string
			auctionTimeSelect := widget.NewSelect([]string{
				"24 часа",
				"3 дня",
				"7 дней",
			}, func(s string) {
				auctionTime = s
			})
			auctionTimeSelect.PlaceHolder = "-- Длительность аукциона --"

			submitBtn := widget.NewButton("Подтвердить", func() {
				a.auctionVM.RenewAuction(cleanedData.CardUrl, auctionCategory, auctionTime)
				selectDialog.Hide()
			})

			content := container.NewVBox(
				auctionCategorySelect,
				auctionTimeSelect,
				submitBtn)

			selectDialog = dialog.NewCustom("Рестарт аукциона без изменений", "Отмена", content, *a.parentWindow)
			selectDialog.Show()
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
