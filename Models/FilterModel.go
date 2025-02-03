package Models

type FilterModel struct {
	accountModel *AccountModel
}

func NewFilterModel() *FilterModel {
	filterModel := &FilterModel{}

	return filterModel
}

func (m *FilterModel) GetAll() *AccountGeneralJson {
	return m.accountModel.GetData(0)
}

func (m *FilterModel) GetByFilter(filter string) *AccountPayloadJson { return nil }

func (m *FilterModel) GetNewestFirst()     {}
func (m *FilterModel) GetOldestFirst()     {}
func (m *FilterModel) GetWithBidFirst()    {}
func (m *FilterModel) GetWithoutBidFirst() {}

func (m *FilterModel) GetActiveOnly()   {}
func (m *FilterModel) GetFinishedOnly() {}

func (m *FilterModel) GetNextPage() {}
func (m *FilterModel) GetPrevPage() {}

//func (m *FilterModel) DataRefresh() {} ?
