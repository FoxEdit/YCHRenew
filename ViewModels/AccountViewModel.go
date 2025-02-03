package ViewModels

import (
	"github.com/FoxEdit/YCHRenew/Models"
	"log"
)

type CleanAccountData struct {
	ImgUrl   string
	Title    string
	Subtitle string
	Ends     string
	EndsUnix int64
	EndsIn   string
	Heat     float64
	Bid      *Models.AccountBidJson
	CardUrl  string
}

type AccountViewModel struct {
	accountModel *Models.AccountModel
	data         *Models.AccountGeneralJson
}

func NewAccountViewModel(accountModel *Models.AccountModel) *AccountViewModel {
	return &AccountViewModel{accountModel: accountModel, data: nil}
}

func (c *AccountViewModel) UpdateDataFromAccount() {
	err := c.accountModel.FetchData(1)
	if err != nil {
		log.Println("FETCH DATA ERROR")
	}

	c.data = c.accountModel.GetData(1)
}

func (c *AccountViewModel) LoadCachedAccountData() {
	c.data = c.accountModel.GetData(1)
}

func (c *AccountViewModel) GetAllRaw() *Models.AccountGeneralJson {
	return c.data
}

func (c *AccountViewModel) GetAllCleaned() []CleanAccountData {
	cleaned := make([]CleanAccountData, 0)

	if c.data == nil {
		return cleaned
	}

	for _, i := range c.data.Payload {
		cleaned = append(cleaned, CleanAccountData{
			i.AltThumb,
			i.Title,
			i.Subtitle,
			i.Ends,
			i.Endsunix,
			i.EndsRel,
			i.Heat,
			i.Bid,
			i.Url,
		})
	}

	return cleaned
}
