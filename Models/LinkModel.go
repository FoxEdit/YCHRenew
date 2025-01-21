package Models

type LinkModel struct {
	urls map[string]string
}

func NewLinkModel() *LinkModel {
	model := new(LinkModel)

	model.urls = map[string]string{
		"profile": "https://account.commishes.com/",
		"crm":     "https://crm.commishes.com/",
		"avatar":  "https://account.commishes.com/image/user/260366/64/?t=1736195748",
	}

	return model
}

func (m LinkModel) GetLinks() map[string]string {
	return m.urls
}
