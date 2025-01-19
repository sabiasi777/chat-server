package auth

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sabiasi777/chat-server/db"
	"github.com/sabiasi777/chat-server/utils"
)

func LoadAuthPage(myWindow fyne.Window) {
	// Email and Password for Login
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	// Username for Registration
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")
	registerEmailEntry := widget.NewEntry()
	registerEmailEntry.SetPlaceHolder("Email")
	registerPasswordEntry := widget.NewPasswordEntry()
	registerPasswordEntry.SetPlaceHolder("Password")

	// Label to show error messages or login/register status
	statusLabel := widget.NewLabel("")

	// Login button click handler
	loginButton := widget.NewButton("Login", func() {
		credentials := utils.User{Email: emailEntry.Text, PasswordHash: passwordEntry.Text}
		err := db.Login(credentials)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Login failed: %v", err))
		} else {
			statusLabel.SetText("Login successful!")
			// Here you can transition to the main page or show main content
			LoadMainPage(myWindow)
		}
	})

	// Register button click handler
	registerButton := widget.NewButton("Register", func() {
		credentials := utils.User{Username: usernameEntry.Text, Email: registerEmailEntry.Text, PasswordHash: registerPasswordEntry.Text}
		err := db.Register(credentials)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Registration failed: %v", err))
		} else {
			statusLabel.SetText("Registration successful!")
			// Here you can transition to the login page or show main content
		}
	})

	// Create the layout for the window
	loginPage := container.NewVBox(
		// Centered Login label
		widget.NewLabelWithStyle("Login", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		// Space between the Login label and the email/password fields
		layout.NewSpacer(),
		// Login fields (email and password)
		emailEntry,
		passwordEntry,
		loginButton,
		// Space between Login button and Registration label
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Registration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		// Space between Registration label and the registration fields
		layout.NewSpacer(),
		// Registration fields (username, email, password)
		container.NewVBox(
			usernameEntry,
			registerEmailEntry,
			registerPasswordEntry,
			registerButton,
		),
		// Status label for success/failure messages
		statusLabel,
	)

	// Ensure fields are visible and not constrained
	emailEntry.Resize(fyne.NewSize(300, 40))
	passwordEntry.Resize(fyne.NewSize(300, 40))
	usernameEntry.Resize(fyne.NewSize(300, 40))
	registerEmailEntry.Resize(fyne.NewSize(300, 40))
	registerPasswordEntry.Resize(fyne.NewSize(300, 40))

	// Set content of the window and ensure it is centered
	myWindow.SetContent(container.NewVBox(loginPage))

	// Show the window
	myWindow.Show()
}

func LoadMainPage(w fyne.Window) {

	chatList, err := db.LoadChats()
	if err != nil {
		panic(err)
	}

	messages, err := db.LoadMessages()
	if err != nil {
		panic(err)
	}

	listView := widget.NewList(func() int {
		return len(chatList) // length of list slice
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(chatList[id].Name) // ID chatResults.Results[id].Name
	})

	contentText := widget.NewLabel("Please select a chat")
	contentText.Wrapping = fyne.TextWrapWord

	inputField := widget.NewEntry()
	inputField.SetPlaceHolder("Type your message...")

	sendButton := widget.NewButton("Send", func() {
		if inputField.Text != "" {
			// handle message send and store to the database

			fmt.Println("Message sent:", inputField.Text)
			inputField.SetText("")
		}
	})

	secondContainer := container.NewVBox(contentText)
	messagesContainer := container.NewVBox()

	listView.OnSelected = func(id widget.ListItemID) {
		secondContainer.Objects = nil
		messagesContainer.Objects = nil
		contentText.Text = "Selected chat: " + chatList[id].Name
		contentText.Refresh()

		//content := container.NewVBox(inputField, sendButton)

		for _, msg := range messages {
			if msg.ChatRoomID-1 == id { // -1 because an indexation of a database and an array doesn's match
				messageLabel := widget.NewLabel(fmt.Sprintf("Sender %d: %s", msg.SenderID, msg.Message))
				messagesContainer.Add(messageLabel)
			}
		}

		contentContainer := container.NewBorder(
			contentText, // Top (chat content header),
			container.NewVBox(inputField, sendButton), // Bottom (input bar and button),
			messagesContainer,
			nil,
		)

		//secondContainer.Add(contentText)
		secondContainer.Add(contentContainer)
		secondContainer.Refresh()

	}

	split := container.NewHSplit(
		listView,
		secondContainer,
	)
	split.Offset = 0.2

	w.SetContent(split)

	w.Show()
}
