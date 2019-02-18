package model

type URL struct {
	ID        int
	URL       string
	Status    string
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
