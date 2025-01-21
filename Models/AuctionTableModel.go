package Models

import (
	"fyne.io/fyne/v2"
)

type auctionElement struct {
	Id       int         `json:"id"`
	Url      string      `json:"url"`
	User     int         `json:"user"`
	UserId   int         `json:"userId"`
	UserName string      `json:"userName"`
	UserKofi bool        `json:"userKofi"`
	UserImg  string      `json:"userImg"`
	Title    string      `json:"title"`
	Subtitle string      `json:"subtitle"`
	Bid      interface{} `json:"bid"`
	Startbid int         `json:"startbid"`
	Ends     string      `json:"ends"`
	Endsunix int         `json:"endsunix"`
	EndsRel  string      `json:"endsRel"`
	EndsRelO struct {
		Amt  int    `json:"amt"`
		Unit string `json:"unit"`
	} `json:"endsRelO"`
	Started  int         `json:"started"`
	Rating   int         `json:"rating"`
	Adult    bool        `json:"adult"`
	Promoted bool        `json:"promoted"`
	Sale     bool        `json:"sale"`
	Heat     float64     `json:"heat"`
	Removed  interface{} `json:"removed"`
	Thumb    string      `json:"thumb"`
	AltThumb string      `json:"altThumb"`
	AvgColor string      `json:"avgColor"`
	Crm      interface{} `json:"crm"`
	Reports  bool        `json:"reports"`
}
type AuctionTableData struct {
	Result  string           `json:"result"`
	Time    int              `json:"time"`
	Pages   int              `json:"pages"`
	Payload []auctionElement `json:"payload"`
}

type AuctionTableModel struct {
	auctions []AuctionTableData
}

func NewAuctionTableModel() *AuctionTableModel {
	return &AuctionTableModel{auctions: make([]AuctionTableData, 0)}
}

func (m AuctionTableModel) GetAll() []AuctionTableData {
	return m.auctions
}

func (m AuctionTableModel) GetByFilter() []AuctionTableData {
	return nil
}

func (m AuctionTableModel) DataRefresh() {
	// get from https://ych.commishes.com/account.json "pages" field -> get all posts from all pages
	// -> store them into file cache and RAM cache. maybe show this process visually for better user experience
}

func (m AuctionTableModel) WriteToFileCache() {

}

func (m AuctionTableModel) ReadFromFileCache() {

}

func (m AuctionTableModel) Search() fyne.CanvasObject {
	return nil
}
