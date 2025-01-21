package ViewModels

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/FoxEdit/YCHRenew/Models"
)

type AuctionTableViewModel struct {
	auctionTableModel *Models.AuctionTableModel

	UIAuctionsList binding.StringList
	UIsearchBar    binding.String
}

func NewAuctionTableViewModel(am *Models.AuctionTableModel) *AuctionTableViewModel {
	avm := &AuctionTableViewModel{
		auctionTableModel: am,

		UIAuctionsList: binding.NewStringList(),
		UIsearchBar:    binding.NewString(),
	}

	avm.LoadData()

	return avm
}

func (avm *AuctionTableViewModel) LoadData() {
	//get data from model and set to VM "buffer"
	avm.UIAuctionsList.Set([]string{"auc1", "auc2", "auc3", "auc4", "auc5"})
	//avm.UIAuctionsList.Set([]string{ /* ... отфильтрованные элементы ... */ })
}

func (avm *AuctionTableViewModel) GetFilteredData() {
	//filter, _ := avm.UIsearchBar.Get()
	//avm.auctionTableModel.Search()
}
