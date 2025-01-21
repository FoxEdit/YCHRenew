package ViewModels

import "github.com/FoxEdit/YCHRenew/Models"

type PopupViewModel struct {
	model *Models.PopupModel
}

func NewPopupViewModel(model *Models.PopupModel) *PopupViewModel {
	controller := new(PopupViewModel)
	return controller
}
