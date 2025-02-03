package Views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FoxEdit/YCHRenew/ViewModels"
	"github.com/FoxEdit/YCHRenew/Views/CustomUITools"
	"image/color"
	"log"
)

const (
	avatarSizeH = 50
	avatarSizeW = 50
)

type Header struct {
	linkViewModel *ViewModels.LinkViewModel
	authViewModel *ViewModels.AuthViewModel

	parentWindow *fyne.Window
}

func NewHeaderContent(links *ViewModels.LinkViewModel, auth *ViewModels.AuthViewModel, parentWindow *fyne.Window) *Header {
	return &Header{linkViewModel: links, authViewModel: auth, parentWindow: parentWindow}
}

func (h *Header) Build() fyne.CanvasObject {
	headerAvatar := h.createAvatar()

	authDialog := h.createAuthChoicePopup()
	var authBtn *widget.Button
	authBtn = widget.NewButton("Вход", func() {
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(authBtn)
		authDialog.ShowAtPosition(pos.Add(fyne.NewPos(0, authBtn.Size().Height)))
	})

	profileLink, _ := h.linkViewModel.GetLinkByName("profile")
	profileHyperlink := widget.NewHyperlink("Профиль", h.linkViewModel.GetUrlFromRawString(profileLink))

	crmLink, _ := h.linkViewModel.GetLinkByName("crm")
	crmHyperlink := widget.NewHyperlink("CRM", h.linkViewModel.GetUrlFromRawString(crmLink))

	createAuctionButton := widget.NewButton("Новый аукцион", func() { h.createNewAuctionDialog().Show() })

	header := container.NewPadded(container.NewHBox(
		headerAvatar,
		CustomUITools.NewWSpacer(55),
		authBtn,
		layout.NewSpacer(),
		profileHyperlink,
		crmHyperlink,
		createAuctionButton,
	))

	return header
}

func (h *Header) afterAuthCallback() {

}

func (h *Header) createAuthChoicePopup() *widget.PopUpMenu {
	authMenuItems := []*fyne.MenuItem{
		fyne.NewMenuItem("Загрузить предыдущую сессию", func() {
			err := h.authViewModel.LoadSessionButtonFunctional()
			if err != nil {
				dialog.NewInformation("Ошибка!", err.Error(), *h.parentWindow).Show()
			}
		}),
		fyne.NewMenuItem("Новый вход", func() { h.createLoginDialog().Show() }),
	}

	authChoicePopupMenu := widget.NewPopUpMenu(fyne.NewMenu("Вариант входа", authMenuItems...), (*h.parentWindow).Canvas())

	return authChoicePopupMenu
}

// TODO use enums? separate arrays?
func (h *Header) createNewAuctionDialog() *dialog.CustomDialog {
	var auctionDialog *dialog.CustomDialog

	categoryLabel := widget.NewLabel("Категория работы")
	category := widget.NewRadioGroup([]string{
		"Адопт",
		"Фурри",
		"Пони",
		"Человек",
		"Самоделка",
	}, func(s string) {
		log.Println(s, "CLICKED")
	})
	category.Required = true
	category.Selected = "Адопт"

	subtitleLabel := widget.NewLabel("Подзаголовок")
	subtitle := widget.NewEntry()

	titleLabel := widget.NewLabel("Заголовок")
	title := widget.NewEntry()

	descriptionLabel := widget.NewLabel("Описание")
	descriptionEntry := widget.NewEntry()
	description := container.NewStack(CustomUITools.NewHSpacer(200), descriptionEntry)

	descriptionEntry.MultiLine = true
	descriptionEntry.Wrapping = fyne.TextWrapWord

	ageRestrictionsLabel := widget.NewLabel("Возрастные ограничения")
	ageRestrictions := widget.NewRadioGroup([]string{
		"Безопасный контент (S)",
		"Сомнительный контент (Q)",
		"Взрослый контент (E)",
		"Шок-контент",
	}, func(s string) {
		log.Println(s, "CLICKED")
	})
	ageRestrictions.Required = true
	ageRestrictions.Selected = "Безопасный контент (S)"

	additionalOptionsLabel := widget.NewLabel("Дополнительные опции")
	additionalOptions := widget.NewCheckGroup([]string{
		"Доступен дополнительный NSFW контент",
		"Запретить снайпинг ставки",
	}, func(options []string) {
		log.Println(options, "CLICKED")
	})
	additionalOptions.Horizontal = true

	auctionDurationLabel := widget.NewLabel("Длительность аукциона")
	auctionDuration := widget.NewRadioGroup([]string{
		"24 часа",
		"3 дня",
		"7 дней",
	}, func(s string) {
		log.Println(s, "CLICKED")
	})
	auctionDuration.Required = true
	auctionDuration.Selected = "24 часа"

	postNewAuction := widget.NewButton("Запустить аукцион", func() {
		log.Println("AUCTION RUN CLICKED")
	})

	auctionDialogContainer := container.NewStack(CustomUITools.NewHWSpacer(420, 625),
		container.NewVScroll(
			container.NewBorder(
				nil,
				nil,
				nil,
				CustomUITools.NewWSpacer(12),
				container.NewVBox(
					categoryLabel,
					category,
					CustomUITools.NewSeparator(),
					subtitleLabel,
					subtitle,
					titleLabel,
					title,
					descriptionLabel,
					description,
					ageRestrictionsLabel,
					ageRestrictions,
					CustomUITools.NewSeparator(),
					additionalOptionsLabel,
					additionalOptions,
					CustomUITools.NewSeparator(),
					auctionDurationLabel,
					auctionDuration,
					CustomUITools.NewSeparator(),
					postNewAuction,
				))))

	auctionDialog = dialog.NewCustom("Создание аукциона", "Отмена", auctionDialogContainer, *h.parentWindow)

	return auctionDialog
}

func (h *Header) createLoginDialog() *dialog.CustomDialog {
	var authChoiceDialog *dialog.CustomDialog

	login := widget.NewEntry()
	password := widget.NewPasswordEntry()

	dialogContent := container.NewVBox(
		login,
		password,
		CustomUITools.NewHSpacer(10.0),
		CustomUITools.NewSeparator(),
		CustomUITools.NewHSpacer(10.0),
		widget.NewButton("Войти", func() {
			err := h.authViewModel.LoginButtonFunctional(login.Text, password.Text)
			authChoiceDialog.Hide()
			if err != nil {
				dialog.NewInformation("Ошибка!", err.Error(), *h.parentWindow).Show()
			}
		}),
	)

	authChoiceDialog = dialog.NewCustom("Выберите метод авторизации", "Отмена", dialogContent, *h.parentWindow)

	return authChoiceDialog
}

func (h *Header) createEmptyAvatar() fyne.CanvasObject {
	avatar := canvas.NewCircle(color.RGBA{R: 110, G: 110, B: 110, A: 255})
	avatar.Resize(fyne.NewSize(avatarSizeW, avatarSizeH))
	avatar.Move(fyne.NewPos(0, -5))

	avatarContainer := container.NewWithoutLayout(avatar)

	return avatarContainer
}

func (h *Header) createAvatar() fyne.CanvasObject {
	avatarLink, ok := h.linkViewModel.GetLinkByName("avatar")
	if ok != nil {
		return h.createEmptyAvatar()
	}

	URI, ok := h.linkViewModel.GetFyneURIFromString(avatarLink)

	if ok != nil {
		fyne.LogError("Cannot get URI from str", ok)
		return h.createEmptyAvatar()
	}

	circle := canvas.NewCircle(color.Transparent)
	circle.Resize(fyne.NewSize(avatarSizeW+40, avatarSizeH+40))
	circle.Move(fyne.NewPos(circle.Position().X-20, circle.Position().Y-25))
	circle.StrokeColor = theme.Color("background")
	circle.StrokeWidth = 20
	circle.FillColor = color.Transparent

	canvasAvatar := canvas.NewImageFromURI(URI)
	canvasAvatar.Resize(fyne.NewSize(avatarSizeW, avatarSizeH))
	canvasAvatar.Move(fyne.NewPos(0, -5))

	avatar := container.NewWithoutLayout(canvasAvatar, circle)
	return avatar
}
