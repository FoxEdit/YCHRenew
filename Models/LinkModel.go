package Models

type LinkModel struct {
	defaultUrls map[string]string
}

func NewLinkModel() *LinkModel {
	model := new(LinkModel)

	model.defaultUrls = map[string]string{
		"profile": "https://ych.commishes.com/account/",
		"crm":     "https://crm.commishes.com/",
	}

	model.loadAdditionalUrls()

	return model
}

func (m *LinkModel) loadAdditionalUrls() {
	client := getWebClientInstance()
	if client.isAuthenticated {
		m.loadAvatarLink()
	}

}

func (m *LinkModel) loadAvatarLink() {
	accountModel := GetAccountModelInstance()
	data := accountModel.GetData(0)
	if data == nil {
		accountModel.FetchData(1)
		data = accountModel.GetData(0)
	}

	if data != nil {
		m.defaultUrls["avatar"] = data.Payload[0].UserImg
	}
}

func (m *LinkModel) GetLinks() map[string]string {
	m.loadAdditionalUrls()
	return m.defaultUrls
}
