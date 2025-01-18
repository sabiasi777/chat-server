package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sabiasi777/chat-server/utils"
	"golang.org/x/crypto/bcrypt"
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

func CheckAuth() {

}

func Login(User utils.User) error {
	stmt, err := db.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := db.Query(User.Email)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		fmt.Println("Invalid email")
	}

	var user utils.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
		if err != nil {
			return err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(User.PasswordHash)) // nill means it is a match

	if err != nil {
		fmt.Println("Invalid password")
	}

	return nil
}

func Register(User utils.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	User.PasswordHash = string(hashedPassword)
	_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		User.Username, User.Email, User.PasswordHash)

	if err != nil {
		return err
	}

	return nil
}

func LoadChats() ([]utils.ChatRoom, error) {
	// get data from the database
	rows, err := db.Query("SELECT id, name FROM chat_rooms")
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

func LoadMessages() ([]utils.Message, error) {
	rows, err := db.Query("SELECT id, chat_room_id, sender_id, message FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []utils.Message
	for rows.Next() {
		var message utils.Message
		err := rows.Scan(&message.ID, &message.ChatRoomID, &message.SenderID, &message.Message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
