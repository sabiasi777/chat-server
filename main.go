package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/sabiasi777/chat-server/auth"
	"github.com/sabiasi777/chat-server/db"
)

func main() {
	fmt.Println("main function...")
	db.ConnectToDB()

	defer func() {
		if db.DB != nil {
			db.DB.Close()
		}
	}()

	a := app.New()
	w := a.NewWindow("Chat-Server")
	w.Resize(fyne.NewSize(500, 500))

	auth.LoadAuthPage(w)
	a.Run()
}
