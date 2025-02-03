package Views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views/CustomUITools"
)

type Filter struct {
	navigationFilterViewModel *ViewModels.FilterViewModel
	cardViewModel             *ViewModels.AccountViewModel
	searchBar                 *widget.Entry
}

func NewNavigationFilter(avm *ViewModels.FilterViewModel) *Filter {
	return &Filter{
		navigationFilterViewModel: avm,
		searchBar:                 widget.NewEntryWithData(avm.UISearchBar),
	}
}

func (n *Filter) Build() fyne.CanvasObject {
	background := canvas.NewRectangle(theme.Color("background")) //(color.RGBA{R: 44, G: 44, B: 46, A: 0})
	background.SetMinSize(fyne.NewSize(200, 460))

	n.searchBar.SetPlaceHolder("Фильтр по названиям")

	auctionLabel := widget.NewLabel("Показывать только")
	auctionRadioGroup := widget.NewRadioGroup([]string{"Все аукционы", "Активные аукционы", "Прошедшие аукционы"}, func(s string) {
		fmt.Println(s)
	})
	auctionRadioGroup.SetSelected("Все аукционы")
	auctionRadioGroup.Required = true

	sortLabel := widget.NewLabel("Сортировать как")
	sortRadioGroup := widget.NewRadioGroup([]string{"Сначала новые", "Сначала старые", "Сначала со ставкой", "Сначала без ставки"}, func(s string) {
		fmt.Println(s)
	})
	sortRadioGroup.SetSelected("Сначала новые")
	sortRadioGroup.Required = true

	updateBtn := widget.NewButton("Обновить", func() {
		n.navigationFilterViewModel.UpdateAuctionTableList()
	})

	nextPageBtn := widget.NewButton(">>", func() {})
	prevPageBtn := widget.NewButton("<<", func() {})
	pageNavContainer := container.NewGridWithColumns(2, prevPageBtn, nextPageBtn)

	content := container.NewStack(background, container.NewPadded(container.NewVBox(
		n.searchBar,
		auctionLabel,
		auctionRadioGroup,
		CustomUITools.NewSeparator(),
		sortLabel,
		sortRadioGroup,
		CustomUITools.NewSeparator(),
		updateBtn,
		pageNavContainer,
	)))

	return content
}

func (n *Filter) createSearchbar() {
	n.searchBar.SetPlaceHolder("Фильтр по названию")
}
