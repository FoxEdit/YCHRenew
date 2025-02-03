package ViewModels

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/FoxEdit/YCHRenew/Models"
	"log"
	"net/url"
)

type LinkViewModel struct {
	model *Models.LinkModel
}

func NewLinkViewModel(model *Models.LinkModel) *LinkViewModel {
	return &LinkViewModel{model}
}

func (l *LinkViewModel) GetLinkByName(key string) (string, error) {
	link, ok := l.model.GetLinks()[key]
	if !ok {
		log.Println("LINK ERROR: LINK FOR", key, "NOT FOUND")
		return "", errors.New("link not found")
	}
	return link, nil
}

func (l *LinkViewModel) GetFyneURIFromString(strURL string) (fyne.URI, error) {
	parsedURL, ok := url.Parse(strURL)

	if ok != nil {
		return nil, ok
	}

	parsedFyneURI, ok := storage.ParseURI(parsedURL.String())

	if ok != nil {
		return nil, ok
	}

	return parsedFyneURI, nil
}

func (l *LinkViewModel) GetUrlFromRawString(strURL string) *url.URL {
	parsedURL, ok := url.Parse(strURL)
	if ok != nil {
		return &url.URL{}
	}

	return parsedURL
}
