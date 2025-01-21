package Views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

type AuctionTable struct {
}

func NewAuctionTable() *AuctionTable {
	return &AuctionTable{}
}

func (a *AuctionTable) Build() fyne.CanvasObject {

	return container.NewVScroll(container.NewVBox(
		a.CreateCard(),
		a.CreateCard(),
		a.CreateCard(),
		a.CreateCard(),
		a.CreateCard(),
		a.CreateCard(),
		a.CreateCard(),
		a.CreateCard(),
	))
}

func (a *AuctionTable) CreateCard() fyne.CanvasObject {
	card := container.NewHBox()
	cardBackground := canvas.NewRectangle(color.RGBA{R: 32, G: 32, B: 35, A: 255})
	cardBackground.SetMinSize(fyne.NewSize(600, 90))

	card.Add(cardBackground)
	card.Add(NewWSpacer(15))

	return card
}

func (a *AuctionTable) getCardImage() fyne.CanvasObject {
	return nil
}

func (a *AuctionTable) getCardTitle() fyne.CanvasObject {
	return nil
}

func (a *AuctionTable) getCardSubtitle() fyne.CanvasObject {
	return nil
}

func (a *AuctionTable) getCardBidInfo() fyne.CanvasObject {
	return nil
}

func (a *AuctionTable) getCardTimeInfo() fyne.CanvasObject {
	return nil
}

func (a *AuctionTable) cardFunctionality() fyne.CanvasObject {
	return nil
}
