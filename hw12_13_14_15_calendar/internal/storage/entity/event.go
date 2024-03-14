package entity

type Event struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Date        string `db:"date" json:"date"`
	Duration    string `db:"duration" json:"duration"`
	Description string `db:"description" json:"description"`
	OwnerID     int    `db:"owner_id" json:"ownerId"`
}
