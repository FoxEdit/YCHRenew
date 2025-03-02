package ViewModels

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/FoxEdit/YCHRenew/Models"
)

type CardItem struct {
	Data *CleanAccountData
	Card fyne.CanvasObject
}

type FilterViewModel struct {
	filterModel  *Models.FilterModel
	accountModel *Models.AccountModel
	cards        []*CardItem

	onUpdateTable func()

	UISearchBar    binding.String
	UIAuctionsList binding.StringList
}

func (f *FilterViewModel) AddNewCard(card *CardItem) {
	f.cards = append(f.cards, card)
}

func (f *FilterViewModel) SetUIRefreshCallback(callback func()) {
	f.onUpdateTable = callback
}

func NewFilterViewModel(filterModel *Models.FilterModel, accountModel *Models.AccountModel) *FilterViewModel {
	avm := &FilterViewModel{
		filterModel:    filterModel,
		accountModel:   accountModel,
		UIAuctionsList: binding.NewStringList(),
		UISearchBar:    binding.NewString(),
	}

	avm.LoadData()

	return avm
}

func (f *FilterViewModel) UpdateAuctionTableList() {
	log.Println("CLEAN FILTER_CARDS AND UPDATE AUCTION TABLE LIST")
	f.cards = nil // clean
	f.accountModel.FetchData(1)
	f.onUpdateTable()
}

func (f *FilterViewModel) LoadData() {
	//get data from model and set to VM "buffer"
	f.UIAuctionsList.Set([]string{"auc1", "auc2", "auc3", "auc4", "auc5"})
	//avm.UIAuctionsList.Set([]string{ ... отфильтрованные элементы ... })
}

func (f *FilterViewModel) GetFilteredData() {
	//filter, _ := avm.UIsearchBar.Get()
	//avm.auctionTableModel.Search()
}
