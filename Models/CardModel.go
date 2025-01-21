package Models

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"
	"time"
)

type actions struct {
	renew     widget.Button
	fastRenew widget.Button
}

type CardModel struct {
	image      fyne.URI
	title      string
	subtitle   string
	lastBid    time.Time
	endTime    time.Time
	endTimeStr string
	actions    actions
}

func NewCardModel() *CardModel {
	return &CardModel{}

}

func (c *CardModel) GetAllCardsFromAccount() *AuctionTableData {
	client := NewAuthModel().GetAuthorizedClient()
	req, _ := http.NewRequest("GET", "https://ych.commishes.com/account/index.json?page=1", nil)
	response, ok := client.Do(req)

	if ok != nil {
		fyne.LogError("Error fetching cards", ok)
	}
	defer response.Body.Close()

	tableData := &AuctionTableData{}
	responseData, _ := io.ReadAll(response.Body)
	json.Unmarshal(responseData, tableData)

	println(tableData.Pages)

	return tableData
}
