package Views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FoxEdit/YCHRenew/ViewModels"
)

type NavigationFilter struct {
	auctionTableVM *ViewModels.AuctionTableViewModel
	searchBar      *widget.Entry
}

func NewNavigationFilter(avm *ViewModels.AuctionTableViewModel) *NavigationFilter {
	return &NavigationFilter{
		auctionTableVM: avm,
		searchBar:      widget.NewEntryWithData(avm.UIsearchBar),
	}
}

func (n *NavigationFilter) Build() fyne.CanvasObject {
	background := canvas.NewRectangle(theme.Color("background")) //(color.RGBA{R: 44, G: 44, B: 46, A: 0})
	background.SetMinSize(fyne.NewSize(200, 420))

	n.searchBar.SetPlaceHolder("Фильтр по названиям")

	auctionLabel := widget.NewLabel("Показывать только:")
	auctionRadioGroup := widget.NewRadioGroup([]string{"Все аукционы", "Активные аукционы", "Прошедшие аукционы"}, func(s string) {
		fmt.Println(s)
	})
	auctionRadioGroup.SetSelected("Все аукционы")
	auctionRadioGroup.Required = true

	sortLabel := widget.NewLabel("Сортировать как:")
	sortRadioGroup := widget.NewRadioGroup([]string{"Сначала новые", "Сначала старые", "Сначала со ставкой", "Сначала без ставки"}, func(s string) {
		fmt.Println(s)
	})
	sortRadioGroup.SetSelected("Сначала новые")
	sortRadioGroup.Required = true

	btnTest := widget.NewButton("Обновить", func() {})

	content := container.NewStack(background, container.NewPadded(container.NewVBox(
		n.searchBar,
		auctionLabel,
		auctionRadioGroup,
		NewSeparator(),
		sortLabel,
		sortRadioGroup,
		btnTest,
	)))

	return content
}

func (n *NavigationFilter) createSearchbar() {
	n.searchBar.SetPlaceHolder("Фильтр по названию")
}
