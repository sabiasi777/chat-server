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
	myWindow.Resize(fyne.NewSize(1920, 1080))
	myWindow.Show()
}

func LoadMainPage(w fyne.Window) {
	var chatID int
	var currentChat utils.ChatRoom
	chatList, err := db.LoadChats(db.CurrentUser.ID)
	if err != nil {
		fmt.Errorf("couldn't load chats: %v", err)
		return
	}

	messages, err := db.LoadMessages()
	if err != nil {
		fmt.Errorf("couldn't load messages: %v", err)
		return
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

	secondContainer := container.NewVBox(contentText)

	messagesContainer := container.NewVBox()

	messagesContainer.Resize(fyne.NewSize(600, 800))
	messagesContainer.Layout = layout.NewVBoxLayout()

	scrollContainer := container.NewScroll(messagesContainer)
	scrollContainer.SetMinSize(fyne.NewSize(600, 800))

	sendButton := widget.NewButton("Send", func() {
		if inputField.Text != "" {
			db.SendMessage(utils.Message{ChatRoomID: chatID, SenderID: db.CurrentUser.ID, Message: inputField.Text})
			inputField.SetText("")

			messages, err := db.LoadMessages()
			if err != nil {
				fmt.Errorf("couldn't load messages: %v", err)
				return
			}

			// Clear current messages
			messagesContainer.Objects = nil

			// Add the messages to the container
			for _, msg := range messages {
				if msg.ChatRoomID == chatID {
					senderUsername, err := db.LoadUserByID(msg.SenderID)
					if err != nil {
						fmt.Errorf("couldn't fetch user by id")
					}
					messageLabel := widget.NewLabel(fmt.Sprintf("%s: %s", senderUsername, msg.Message))
					messagesContainer.Add(messageLabel)
				}
			}

			// Refresh to update the container with new messages
			messagesContainer.Refresh()
			scrollContainer.Refresh()

			// Scroll to the bottom
			scrollContainer.ScrollToBottom()
		}
	})

	listView.OnSelected = func(id widget.ListItemID) {
		secondContainer.Objects = nil
		messagesContainer.Objects = nil
		contentText.Text = "Selected chat: " + chatList[id].Name
		currentChat, err = db.LoadChatByName(chatList[id].Name)
		if err != nil {
			fmt.Errorf("couldn't fetch chat by name: %v", err)
		}
		chatID = currentChat.ID
		contentText.Refresh()

		for _, msg := range messages {
			if msg.ChatRoomID == chatID {
				senderUsername, err := db.LoadUserByID(msg.SenderID)
				if err != nil {
					fmt.Errorf("couldn't fetch user by id")
				}
				messageLabel := widget.NewLabel(fmt.Sprintf("%s: %s", senderUsername, msg.Message))
				messagesContainer.Add(messageLabel)
			}
		}

		contentContainer := container.NewBorder(
			contentText,
			container.NewVBox(inputField, sendButton),
			nil,
			nil,
			scrollContainer, // Use the scrollable container for messages
		)

		secondContainer.Add(contentContainer)
		secondContainer.Refresh()

		// Scroll to the bottom after loading messages
		scrollContainer.ScrollToBottom()
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
	w.Resize(fyne.NewSize(1920, 1080))
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
	var chatWindow fyne.Window // Declare the window here so it's accessible in the callback
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

		// Show success message and close the window
		dialog.ShowInformation("Success", "Chat created successfully!", w)
		chatWindow.Close() // Close the window
		LoadMainPage(w)    // Reload the main page
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

	// Initialize the "Create New Chat" window
	chatWindow = fyne.CurrentApp().NewWindow("Create New Chat") // Assign the window instance here
	chatWindow.SetContent(split)
	chatWindow.Resize(fyne.NewSize(1920, 1080)) // Adjust the size as needed
	chatWindow.Show()
}
