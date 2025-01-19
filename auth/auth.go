package auth

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	chatList, err := db.LoadChats(db.CurrentUser.ID)
	if err != nil {
		panic(err)
	}

	messages, err := db.LoadMessages()
	if err != nil {
		panic(err)
	}

	// Chat List
	listView := widget.NewList(func() int {
		return len(chatList)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(chatList[id].Name)
	})

	// Chat Content
	contentText := widget.NewLabel("Please select a chat")
	contentText.Wrapping = fyne.TextWrapWord

	inputField := widget.NewEntry()
	inputField.SetPlaceHolder("Type your message...")

	sendButton := widget.NewButton("Send", func() {
		if inputField.Text != "" {
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

		for _, msg := range messages {
			if msg.ChatRoomID-1 == id {
				messageLabel := widget.NewLabel(fmt.Sprintf("Sender %d: %s", msg.SenderID, msg.Message))
				messagesContainer.Add(messageLabel)
			}
		}

		contentContainer := container.NewBorder(
			contentText,
			container.NewVBox(inputField, sendButton),
			nil,
			nil,
			messagesContainer,
		)

		secondContainer.Add(contentContainer)
		secondContainer.Refresh()
	}

	// Create New Chat Button
	createChatButton := widget.NewButton("Create New Chat", func() {
		createNewChat(w)
	})

	// Logout Button
	logoutButton := widget.NewButton("Logout", func() {
		LoadAuthPage(w)
	})

	// Main Layout
	topBar := container.NewHBox(createChatButton, widget.NewSeparator(), logoutButton)
	split := container.NewHSplit(
		container.NewBorder(
			widget.NewLabelWithStyle("Chats", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			nil,
			nil,
			nil,
			listView,
		),
		secondContainer,
	)
	split.Offset = 0.3 // Ensure the left panel has sufficient width

	mainContent := container.NewBorder(
		topBar,
		nil,
		nil,
		nil,
		split,
	)

	w.SetContent(mainContent)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
}

func createNewChat(w fyne.Window) {
	// Fetch all existing users from the database
	users, err := db.LoadUsers()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	// Entry for the chat name
	chatNameEntry := widget.NewEntry()
	chatNameEntry.SetPlaceHolder("Enter chat name...")

	// Container for user selection (checkboxes)
	selectedUsers := map[int]bool{} // Map to track selected user IDs
	userCheckboxes := container.NewVBox()

	// Create checkboxes for each user
	for _, user := range users {
		userID := user.ID // Save user ID to use in callback
		checkBox := widget.NewCheck(user.Username, func(checked bool) {
			selectedUsers[userID] = checked // Track user selection
		})
		userCheckboxes.Add(checkBox)
	}

	// Create and add chat button
	createChatButton := widget.NewButton("Create Chat", func() {
		chatName := chatNameEntry.Text
		if chatName == "" {
			dialog.ShowError(fmt.Errorf("Chat name cannot be empty"), w)
			return
		}

		// Gather selected users
		var memberIDs []int
		for userID, selected := range selectedUsers {
			if selected {
				memberIDs = append(memberIDs, userID)
			}
		}

		if len(memberIDs) == 0 {
			dialog.ShowError(fmt.Errorf("Please select at least one member"), w)
			return
		}

		// Save the new chat to the database
		err := db.CreateChat(chatName, memberIDs)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("Success", "Chat created successfully!", w)
		LoadMainPage(w)
	})

	bottomSection := container.NewBorder(
		widget.NewLabelWithStyle("Select Members", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewVScroll(userCheckboxes),
	)

	topSection := container.NewVBox(
		widget.NewLabel("Create New Chat"),
		chatNameEntry,
		createChatButton,
	)

	// Split the layout into two parts: left and right sections
	split := container.NewVSplit(topSection, bottomSection)
	split.Offset = 0.3 // Control the width ratio of the left section

	// Set the final layout of the window
	chatWindow := fyne.CurrentApp().NewWindow("Create New Chat")
	chatWindow.SetContent(split)
	chatWindow.Resize(fyne.NewSize(600, 400)) // Adjust the size as needed
	chatWindow.Show()
}
