package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sabiasi777/chat-server/utils"
)

var db *sql.DB

func ConnectToDB() {
	// temporarily.
	dsn := "root:skofildi123@tcp(localhost:3306)/chat"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the database")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func LoadChats() ([]utils.ChatRoom, error) {
	// get data from the database
	rows, err := db.Query("SELECT id, name from chat_rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []utils.ChatRoom
	for rows.Next() {
		var chat utils.ChatRoom
		err := rows.Scan(&chat.ID, &chat.Name)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
