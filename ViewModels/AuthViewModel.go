package ViewModels

import (
	"github.com/FoxEdit/YCHRenew/Models"
	"log"
)

type AuthViewModel struct {
	authModel      *Models.AuthModel
	onLoginSuccess func()
}

func NewAuthViewModel(authModel *Models.AuthModel) *AuthViewModel {
	authVM := &AuthViewModel{}
	authVM.authModel = authModel

	return authVM
}

func (a *AuthViewModel) SetUIRefreshCallback(callback func()) {
	a.onLoginSuccess = callback
}

func (a *AuthViewModel) GetLoginButtonFunctional() func(login string, password string) {
	return func(login string, password string) {
		loginErr := a.authModel.Login(login, password)
		if loginErr == nil {
			if a.onLoginSuccess != nil {
				log.Println("SUCCESS LOGIN")
				go a.onLoginSuccess()
			}
		} else {
			log.Println("DEFAULT LOGIN ERROR: ", loginErr)
		}
	}
}

func (a *AuthViewModel) GetLoadSessionButtonFunctional() func() {
	return func() {
		loginErr := a.authModel.CookieLogin()
		if loginErr == nil {
			if a.onLoginSuccess != nil {
				log.Println("SUCCESS LOGIN")
				go a.onLoginSuccess()
			}
		} else {
			log.Println("COOKIE LOGIN ERROR: ", loginErr)
		}
	}
}
