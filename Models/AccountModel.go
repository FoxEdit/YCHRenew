package Models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type AccountBidJson struct {
	User struct {
		Name    string      `json:"name"`
		Email   interface{} `json:"email"`
		Guest   bool        `json:"guest"`
		Account struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Avatar  string `json:"avatar"`
			Sponsor bool   `json:"sponsor"`
		} `json:"account"`
	} `json:"user"`
	Id        int  `json:"id"`
	Bid       int  `json:"bid"`
	Confirmed bool `json:"confirmed"`
	Disabled  bool `json:"disabled"`
	Ommitted  bool `json:"ommitted"`
	Spam      bool `json:"spam"`
	Created   int  `json:"created"`
}
type AccountPayloadJson struct {
	Id       int             `json:"id"`
	Url      string          `json:"url"`
	User     int             `json:"user"`
	UserId   int             `json:"userId"`
	UserName string          `json:"userName"`
	UserKofi bool            `json:"userKofi"`
	UserImg  string          `json:"userImg"`
	Title    string          `json:"title"`
	Subtitle string          `json:"subtitle"`
	Bid      *AccountBidJson `json:"bid"`
	Startbid int             `json:"startbid"`
	Ends     string          `json:"ends"`
	Endsunix int64           `json:"endsunix"`
	EndsRel  string          `json:"endsRel"`
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

type AccountGeneralJson struct {
	Result  string               `json:"result"`
	Time    int                  `json:"time"`
	Pages   int                  `json:"pages"`
	Payload []AccountPayloadJson `json:"payload"`
}

const COMMISHES_URL = "https://ych.commishes.com"

var (
	accountModelInstance *AccountModel
	onceAccountModel     sync.Once
)

type AccountModel struct {
	dataBuffer *AccountGeneralJson
}

func GetAccountModelInstance() *AccountModel {
	onceAccountModel.Do(func() {
		accountModelInstance = &AccountModel{dataBuffer: nil}
	})

	return accountModelInstance
}

// GetData Get from local cache
func (c *AccountModel) GetData(page uint16) *AccountGeneralJson {
	log.Println("CACHE REQUEST")
	return c.dataBuffer
}

// FetchData Fetch form remote resource to cache
func (c *AccountModel) FetchData(page uint16) error {
	if page == 0 {
		page = 1
	}

	log.Println("FETCH REQUEST: ")
	client := getWebClientInstance()
	if !client.isAuthenticated {
		log.Print("UNSUCCESSFUL (unauthorized)")
		return errors.New("unauthorized")
	}
	log.Print("SUCCESSFUL")

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/account.json?page=%d", COMMISHES_URL, page), nil)
	response, ok := client.Do(req)

	if ok != nil || response.StatusCode != 200 {
		if ok != nil {
			log.Println("Error fetching data from account (request error)")
			return errors.New("request error")
		} else if response.StatusCode != 200 {
			log.Println("Error fetching data from account (response != 200): ")
			log.Println(req.URL.String())
			return errors.New("response error")
		} else {
			log.Println("Error fetching data from account (unknown error)")
			return errors.New("unknown fetching error")
		}
	}

	tableData := AccountGeneralJson{}
	responseData, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	err := json.Unmarshal(responseData, &tableData)

	if err != nil {
		log.Println("UNMARSHAL ERROR:", err)
		c.dataBuffer = &AccountGeneralJson{}
		return err
	}

	c.dataBuffer = &tableData

	return nil
}
