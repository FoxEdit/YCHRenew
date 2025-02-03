package ViewModels

import (
	"github.com/FoxEdit/YCHRenew/Models"
	"log"
)

type AuthViewModel struct {
	authModel                       *Models.AuthModel
	onLoginSuccessUIRefreshCallback func()
}

func NewAuthViewModel(authModel *Models.AuthModel) *AuthViewModel {
	authVM := &AuthViewModel{}
	authVM.authModel = authModel

	return authVM
}

func (a *AuthViewModel) SetUIRefreshCallback(callback func()) {
	a.onLoginSuccessUIRefreshCallback = callback
}

func (a *AuthViewModel) LoginButtonFunctional(login string, password string) error {
	loginErr := a.authModel.Login(login, password)
	if loginErr == nil {
		if a.onLoginSuccessUIRefreshCallback != nil {
			log.Println("SUCCESS LOGIN")
			go a.onLoginSuccessUIRefreshCallback()
		}
	}

	log.Println("DEFAULT LOGIN ERROR: ", loginErr)
	return loginErr
}

func (a *AuthViewModel) LoadSessionButtonFunctional() error {
	loginErr := a.authModel.CookieLogin()
	if loginErr == nil {
		if a.onLoginSuccessUIRefreshCallback != nil {
			log.Println("SUCCESS LOGIN")
			go a.onLoginSuccessUIRefreshCallback()
			return nil
		}
	}

	log.Println("COOKIE LOGIN ERROR: ", loginErr)
	return loginErr
}
