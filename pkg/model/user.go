package model

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	TgChatID *int64 `db:"tg_chat_id"`
	Nickname string `db:"nickname"`
	PassHash string `db:"pass_hash"`
}
