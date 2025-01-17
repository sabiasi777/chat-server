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

	chatList, err := db.LoadChats()
	if err != nil {
		panic(err)
	}

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

	inputContainer := container.NewHBox(inputField, sendButton)

	listView.OnSelected = func(id widget.ListItemID) {
		contentText.Text = "Selected chat: " + chatList[id].Name
		contentText.Refresh()

		w.SetContent(container.NewBorder(
			nil,                            // no top part
			inputContainer,                 // bottom part with the input and button
			listView,                       // left side with the list view
			nil,                            // no right part
			container.NewVBox(contentText), // the content (chat messages) area
		))
	}

	w.SetContent(container.NewHSplit(
		listView,
		container.NewMax(contentText),
	))

	// split := container.NewHSplit(
	// 	listView,
	// 	container.NewMax(contentText),
	// )
	// split.Offset = 0.2

	// w.SetContent(split)

	w.ShowAndRun()
}
