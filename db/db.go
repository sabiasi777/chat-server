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
var CurrentUser utils.User

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
			fmt.Println("Login failed: no rows found for the provided email")
			return fmt.Errorf("invalid email or password")
		}

		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		return fmt.Errorf("Invalid email or password")
	}

	CurrentUser = dbUser

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

func LoadChats(ID int) ([]utils.ChatRoom, error) {
	stmt, err := DB.Prepare("SELECT chat_room_id FROM chat_room_members WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var chatRoomIDs []int
	rows, err := stmt.Query(ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chatRoomID int
		if err := rows.Scan(&chatRoomID); err != nil {
			return nil, err
		}
		chatRoomIDs = append(chatRoomIDs, chatRoomID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	stmt1, err := DB.Prepare("SELECT id, name FROM chat_rooms WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	var chatRooms []utils.ChatRoom
	for _, chatRoomID := range chatRoomIDs {
		var chatRoom utils.ChatRoom
		err = stmt1.QueryRow(chatRoomID).Scan(&chatRoom.ID, &chatRoom.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Load chats failed: no rows found for the provided chat IDs")
				return nil, fmt.Errorf("invalid chat_room_id")
			}
			return nil, err
		}
		chatRooms = append(chatRooms, chatRoom)
	}

	return chatRooms, nil
}

func LoadChatByName(chatName string) (utils.ChatRoom, error) {
	stmt, err := DB.Prepare("SELECT id, name FROM chat_rooms WHERE name = ?")
	if err != nil {
		return utils.ChatRoom{}, err
	}
	defer stmt.Close()

	var chatRoom utils.ChatRoom
	err = stmt.QueryRow(chatName).Scan(&chatRoom.ID, &chatRoom.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ChatRoom{}, fmt.Errorf("invalid chat name")
		}
		return utils.ChatRoom{}, err
	}

	return chatRoom, nil
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

func LoadUsers() ([]utils.User, error) {
	rows, err := DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []utils.User
	for rows.Next() {
		var user utils.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateChat(chatName string, memberIDs []int) error {
	fmt.Println("in the create chat function")
	_, err := DB.Exec("INSERT INTO chat_rooms (name) VALUES (?)", chatName)
	if err != nil {
		return err
	}

	fmt.Println("chat created")

	stmt, err := DB.Prepare("SELECT id FROM chat_rooms WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var chatID int
	err = stmt.QueryRow(chatName).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid chat name")
		}
		return err
	}

	fmt.Println("chat room id fetched")

	fmt.Println("Here before memberID for cycle")
	for _, memberID := range memberIDs {
		_, err := DB.Exec("INSERT INTO chat_room_members (chat_room_id, user_id) VALUES (?, ?)", chatID, memberID)
		if err != nil {
			// Handle the error
			fmt.Println("Error inserting member:", err)
			return err
		}
	}

	return nil
}

func SendMessage(message utils.Message) error {
	_, err := DB.Exec("INSERT INTO messages (chat_room_id, sender_id, message) VALUES (?, ?, ?)",
		message.ChatRoomID, message.SenderID, message.Message)

	if err != nil {
		return err
	}

	return nil
}

func LoadUserByID(userID int) (string, error) {
	stmt, err := DB.Prepare("SELECT username FROM users WHERE id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var username string
	err = stmt.QueryRow(userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid userID")
		}
		return "", err
	}

	return username, nil
}
