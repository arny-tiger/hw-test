package entity

type Event struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Date        string `db:"date"`
	Duration    string `db:"duration"`
	Description string `db:"description"`
	OwnerID     int    `db:"owner_id"`
}
