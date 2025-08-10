package auth

import (
	"fmt"
	"regexp"

	"github.com/rivo/tview"
)

func ShowModal(appUI *tview.Application, message string, onClose func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if onClose != nil {
				onClose()
			}
		})
	appUI.SetRoot(modal, true)
}

func CreateLoginForm(appUI *tview.Application, serverURL string, onSubmit func(email, password string)) *tview.Form {
	form := tview.NewForm()

	form.AddTextView("Login to your account at", serverURL, 0, 1, false, false)
	form.AddInputField("Email", "", 20, nil, nil)
	form.AddPasswordField("Password", "", 20, '*', nil)
	form.AddButton("Login", func() {
		email := form.GetFormItemByLabel("Email").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

		if err := validateLoginInputs(email, password); err != nil {
			ShowModal(appUI, err.Error(), func() {
				appUI.SetRoot(form, true)
			})
			return
		}
		onSubmit(email, password)
	})

	return form
}

func CreateRegisterForm(appUI *tview.Application, serverURL string, onSubmit func(email, password string)) *tview.Form {
	form := tview.NewForm()

	form.AddTextView("Register an account at", serverURL, 0, 1, false, false)
	form.AddInputField("Email", "", 20, nil, nil)
	form.AddInputField("Confirm Email", "", 20, nil, nil)
	form.AddPasswordField("Password", "", 20, '*', nil)
	form.AddPasswordField("Confirm Password", "", 20, '*', nil)
	form.AddButton("Register", func() {
		email := form.GetFormItemByLabel("Email").(*tview.InputField).GetText()
		confirmEmail := form.GetFormItemByLabel("Confirm Email").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		confirmPassword := form.GetFormItemByLabel("Confirm Password").(*tview.InputField).GetText()

		if err := validateRegisterInputs(email, confirmEmail, password, confirmPassword); err != nil {
			ShowModal(appUI, err.Error(), func() {
				appUI.SetRoot(form, true)
			})
			return
		}
		onSubmit(email, password)
	})

	return form
}

func validateLoginInputs(email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("Email and Password cannot be empty")
	}
	if !isValidEmail(email) {
		return fmt.Errorf("Invalid email format")
	}
	return nil
}

func validateRegisterInputs(email, confirmEmail, password, confirmPassword string) error {
	if email == "" || confirmEmail == "" || password == "" || confirmPassword == "" {
		return fmt.Errorf("All fields are required")
	}
	if email != confirmEmail {
		return fmt.Errorf("Emails do not match")
	}
	if password != confirmPassword {
		return fmt.Errorf("Passwords do not match")
	}
	if !isValidEmail(email) {
		return fmt.Errorf("Invalid email format")
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
