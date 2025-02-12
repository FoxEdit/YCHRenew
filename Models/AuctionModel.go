package Models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

type Media struct {
	Preview  string `json:"preview"`
	Thumb    string `json:"thumb"`
	Original string `json:"original"`
}

type Slot struct {
	Name        string `json:"name"`
	Startbid    int    `json:"startbid"`
	Autobuy     int    `json:"autobuy"`
	Minincrease int    `json:"minincrease"`
}

// cutted ver
type Payload struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Subtitle    string          `json:"subtitle"`
	Description string          `json:"description"`
	Rating      int             `json:"rating"`
	Adult       bool            `json:"adult"`
	Media       Media           `json:"media"`
	Slots       map[string]Slot `json:"slots"`
}

// cutted ver
type AuctionCardGeneralJson struct {
	Payload Payload `json:"payload"`
}

type AuctionModel struct {
	//auction *AuctionCardGeneralJson
}

func NewAuctionModel() *AuctionModel {
	return &AuctionModel{}
}

func (a *AuctionModel) StartNewAuction() {

}

func (a *AuctionModel) RestartAuctionAsIs(cardUrl string, auctionCategory string, auctionTime string) {
	client := getWebClientInstance()

	data, err := a.getAuctionCardData(cardUrl)
	if err != nil {
		log.Println("[PROBABLY CANCELLED AUCTION] ATTEMPTING TO RESTART A NON-EXISTENT AUCTION!")
		return
	}

	auctionPostXsrf := client.GetXSRFByPattern(
		COMMISHES_URL+"/auction/create/",
		"\"csrf\":\"",
		"\",\"")

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	auctionFileRequest, _ := http.NewRequest("GET", data.Payload.Media.Original, nil)
	auctionFileResponse, _ := client.Do(auctionFileRequest)
	auctionFile, err := io.ReadAll(auctionFileResponse.Body)
	if err != nil {
		log.Println("Error when reading IO")
		return
	}
	defer auctionFileResponse.Body.Close()

	multipartWriter.WriteField("csrf", auctionPostXsrf)

	fileField := make(textproto.MIMEHeader)
	fileField.Set("Content-Disposition", "form-data; name=\"file\"; filename=\"Illustration.json\"")
	fileField.Set("Content-Type", "image/jpeg")
	filePart, _ := multipartWriter.CreatePart(fileField)
	filePart.Write(auctionFile)

	multipartWriter.WriteField("category", auctionCategory)
	multipartWriter.WriteField("subtitle", data.Payload.Subtitle)
	multipartWriter.WriteField("title", data.Payload.Title)
	multipartWriter.WriteField("description", data.Payload.Description)
	multipartWriter.WriteField("rating", strconv.Itoa(data.Payload.Rating))

	multipartWriter.Close()

	newAuctionSubmitMainDataRequest, _ := http.NewRequest("POST", COMMISHES_URL+"/auction/create.json", &requestBody)
	newAuctionSubmitMainDataRequest.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	newAuctionSubmitMainDataResponse, _ := client.Do(newAuctionSubmitMainDataRequest)

	if newAuctionSubmitMainDataResponse.StatusCode != http.StatusOK {
		log.Println("FIRST STEP OF POST ERROR")
		return
	}

	processAuctionPayload := struct {
		Payload struct {
			Processed bool        `json:"processed"`
			Message   interface{} `json:"message"`
			Redirect  string      `json:"redirect"`
		} `json:"payload"`
	}{}
	newAuctionSubmitMainDataProcessResponseStr, _ := io.ReadAll(newAuctionSubmitMainDataResponse.Body)
	err = json.Unmarshal(newAuctionSubmitMainDataProcessResponseStr, &processAuctionPayload)

	if err != nil {
		log.Println("SECOND STEP OF POST ERROR:", err.Error(), processAuctionPayload.Payload.Message)
		return
	}

	redirectParts := strings.Split(processAuctionPayload.Payload.Redirect, "/")
	var newAuctionId string
	for _, val := range redirectParts {
		if _, err := strconv.Atoi(val); err == nil {
			newAuctionId = val
		}
	}

	newAuctionSubmitDurationFormData := url.Values{
		"duration": {auctionTime},
	}
	newAuctionSubmitDurationRequest, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/auction/start/%s/?result=ok/", COMMISHES_URL, newAuctionId),
		strings.NewReader(newAuctionSubmitDurationFormData.Encode()),
	)
	newAuctionSubmitDurationRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	durationRequest, _ := client.Do(newAuctionSubmitDurationRequest)
	fmt.Println(durationRequest.StatusCode)
	reqAnsw, _ := io.ReadAll(durationRequest.Body)
	fmt.Println(string(reqAnsw))

	newAuctionSubmitPricesFormData := url.Values{
		"startingbid":    {strconv.Itoa(data.Payload.Slots[strconv.Itoa(data.Payload.Id)].Startbid)},
		"minincrease":    {strconv.Itoa(data.Payload.Slots[strconv.Itoa(data.Payload.Id)].Minincrease)},
		"autobuyenabled": {"on"},
		"autobuy":        {strconv.Itoa(data.Payload.Slots[strconv.Itoa(data.Payload.Id)].Autobuy)},
	}

	newAuctionSubmitPriceRequest, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/auction/ready/%s/", COMMISHES_URL, newAuctionId),
		strings.NewReader(newAuctionSubmitPricesFormData.Encode()),
	)
	newAuctionSubmitPriceRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	priceRequest, _ := client.Do(newAuctionSubmitPriceRequest)
	fmt.Println(priceRequest.StatusCode)
	priceAnsw, _ := io.ReadAll(priceRequest.Body)
	fmt.Println(string(priceAnsw))
}

func (a *AuctionModel) RestartAuctionWithChanges() {

}

func (a *AuctionModel) getAuctionCardData(cardUrl string) (*AuctionCardGeneralJson, error) {
	client := getWebClientInstance()
	if !client.isAuthenticated {
		log.Println("CARD DATA GET ERROR: UNAUTHORIZED")
		return nil, errors.New("unauthorized")
	}
	log.Println("SUCCESS GETTING CARD DATA")

	req, _ := http.NewRequest("GET", COMMISHES_URL+cardUrl+".json", nil)
	res, _ := client.Do(req)

	if res.StatusCode != 200 {
		log.Printf("\nCARD DATA GET NETWORK ERROR: %s (%s)", res.Status, req.URL.String())
		return nil, errors.New(res.Status)
	}

	resBody, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	cardData := AuctionCardGeneralJson{}
	err := json.Unmarshal(resBody, &cardData)

	if err != nil {
		log.Println("CARD DESERIALIZING ERROR")
		return nil, errors.New("deserialize error")
	}

	return &cardData, nil
}

/*func (a *AuctionModel) getOriginalIllustrationExtension(directUrl string) {
	req, _ := http.NewRequest("GET", directUrl, nil)
	res, _ := getWebClientInstance().client.Do(req)
	res.
}*/

func (a *AuctionModel) postAuction() {

}
