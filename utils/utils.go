package utils

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	JWTToken     string
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
