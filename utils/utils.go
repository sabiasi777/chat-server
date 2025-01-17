package utils

// type ChatList struct {
// 	List []Chat
// }

type Chat struct {
	Title string
}

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
}

type ChatRoom struct {
	ID   int
	Name string
}

type ChatRoomMember struct {
	ChatRoomID int
	UserID     int
}

type Message struct {
	ID         int
	ChatRoomID int
	SenderID   int
	Message    string
}
