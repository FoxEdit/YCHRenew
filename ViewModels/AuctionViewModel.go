package ViewModels

import (
	"github.com/FoxEdit/YCHRenew/Models"
)

// TODO fill this
const (
	FURRY  = "furry"
	HUMAN  = "human"
	ADOPT  = "adopt"
	PONY   = "brony"
	CRAFTS = "crafts"
)

const (
	RATING_SAFE         = 0
	RATING_QUESTIONABLE = 1
	RATING_EXPLICIT     = 2
	RATING_SHOCK        = 3
)

const (
	ONE_DAY_DURATION    = "24"
	THREE_DAYS_DURATION = "3d"
	SEVEN_DAYS_DURATION = "7d"
)

type AuctionViewModel struct {
	auctionModel *Models.AuctionModel
}

func NewAuctionViewModel(model *Models.AuctionModel) *AuctionViewModel {
	return &AuctionViewModel{auctionModel: model}
}

func (a *AuctionViewModel) RenewAuction(url string, auctionCategory string, auctionTime string) {
	switch auctionTime {
	case "24 часа":
		auctionTime = ONE_DAY_DURATION
	case "3 дня":
		auctionTime = THREE_DAYS_DURATION
	case "7 дней":
		auctionTime = SEVEN_DAYS_DURATION
	}

	switch auctionCategory {
	case "Фурри":
		auctionCategory = FURRY
	case "Люди":
		auctionCategory = HUMAN
	case "Адопты":
		auctionCategory = ADOPT
	case "Пони":
		auctionCategory = PONY
	case "Самоделки":
		auctionCategory = CRAFTS
	}

	a.auctionModel.RestartAuctionAsIs(url, auctionCategory, auctionTime)
}
