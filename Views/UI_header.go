package Views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"image/color"
)

const (
	avatarSizeH = 50
	avatarSizeW = 50
)

type HeaderContent struct {
	headerLinks  *ViewModels.LinkViewModel
	headerPopups *ViewModels.PopupViewModel
	parentWindow *fyne.Window
}

func NewHeaderContent(links *ViewModels.LinkViewModel, popups *ViewModels.PopupViewModel, parentWindow *fyne.Window) *HeaderContent {
	return &HeaderContent{headerLinks: links, headerPopups: popups, parentWindow: parentWindow}
}

func (h *HeaderContent) Build() fyne.CanvasObject {
	headerAvatar := h.createAvatar()

	// separate to "createAuthDialog" func
	loginBtn := widget.NewButton("Новая авторизация", func() { fmt.Print("TODO: login btn") })
	loadBtn := widget.NewButton("Предыдущая сессия", func() { fmt.Print("TODO: load btn") })
	authContainer := container.NewVBox(loginBtn, loadBtn, NewHSpacer(10), NewSeparator(), NewHSpacer(10))
	authDialog := dialog.NewCustom("Выберите метод авторизации", "Отмена", authContainer, *h.parentWindow)
	authBtn := widget.NewButton("Авторизация", func() { authDialog.Show() })
	// ------------------------

	profileLink, _ := h.headerLinks.GetLinkByName("profile")
	profileHyperlink := widget.NewHyperlink("Профиль", h.headerLinks.GetUrlFromRawString(profileLink))

	crmLink, _ := h.headerLinks.GetLinkByName("crm")
	crmHyperlink := widget.NewHyperlink("CRM", h.headerLinks.GetUrlFromRawString(crmLink))

	createAuctionButton := widget.NewButton("Новый аукцион", func() { fmt.Print("TODO: new auc btn") })

	header := container.NewPadded(container.NewHBox(
		headerAvatar,
		NewWSpacer(55),
		authBtn,
		layout.NewSpacer(),
		profileHyperlink,
		crmHyperlink,
		createAuctionButton,
	))

	return header
}

func (h *HeaderContent) createEmptyAvatar() fyne.CanvasObject {
	circle := canvas.NewCircle(color.White)
	circle.Resize(fyne.NewSize(avatarSizeW, avatarSizeH))

	text := canvas.NewText("Not found", color.RGBA{R: 255, G: 0, B: 0, A: 150})
	text.TextSize = 10
	text.Move(fyne.NewPos(0.5, 17))

	content := container.NewWithoutLayout(circle, text)

	return content
}

func (h *HeaderContent) createAvatar() fyne.CanvasObject {
	avatarLink, ok := h.headerLinks.GetLinkByName("avatar")
	if ok != nil {
		fyne.LogError("Hashmap has no avatar", ok)
		return h.createEmptyAvatar()
	}

	URI, ok := h.headerLinks.GetFyneURIFromString(avatarLink)

	if ok != nil {
		fyne.LogError("Cannot get URI from str", ok)
		return h.createEmptyAvatar()
	}

	circle := canvas.NewCircle(color.Transparent)
	circle.Resize(fyne.NewSize(avatarSizeW+40, avatarSizeH+40))
	circle.Move(fyne.NewPos(circle.Position().X-20, circle.Position().Y-20))
	circle.StrokeColor = theme.Color("background")
	circle.StrokeWidth = 20
	circle.FillColor = color.Transparent

	canvasAvatar := canvas.NewImageFromURI(URI)
	canvasAvatar.Resize(fyne.NewSize(avatarSizeW, avatarSizeH))

	avatar := container.NewWithoutLayout(canvasAvatar, circle)
	return avatar
}
