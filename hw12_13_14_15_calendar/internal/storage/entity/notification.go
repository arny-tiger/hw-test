package entity

type Notification struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Date   string `db:"date"`
	UserID int    `db:"user_id"`
}
