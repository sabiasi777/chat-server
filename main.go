package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sabiasi777/chat-server/db"
)

func main() {
	fmt.Println("main function...")
	db.ConnectToDB()

	// FETCHES

	chatList, err := db.LoadChats()
	if err != nil {
		panic(err)
	}

	messages, err := db.LoadMessages()
	if err != nil {
		panic(err)
	}

	// ---------

	a := app.New()
	w := a.NewWindow("Chat-Server")
	w.Resize(fyne.NewSize(1200, 800))

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

	w.ShowAndRun()
}
