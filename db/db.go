package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sabiasi777/chat-server/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var DB *sql.DB

func ConnectToDB() {
	dsn := "root:skofildi123@tcp(localhost:3306)/chat"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Connected to the database")

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
}

func Login(user utils.User) error {
	stmt, err := DB.Prepare("SELECT id, username, email, password_hash FROM users WHERE email = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var dbUser utils.User
	err = stmt.QueryRow(user.Email).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid email or password")
		}

		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		return fmt.Errorf("Invalid email or password")
	}

	return nil
}

func Register(user utils.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	_, err = DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		user.Username, user.Email, user.PasswordHash)

	if err != nil {
		return err
	}

	return nil
}

func LoadChats() ([]utils.ChatRoom, error) {
	// get data from the database
	rows, err := DB.Query("SELECT id, name FROM chat_rooms")
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
	rows, err := DB.Query("SELECT id, chat_room_id, sender_id, message FROM messages")
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
