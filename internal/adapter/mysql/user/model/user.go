package model

type User struct {
	ID       int     `db:"id"`
	Email    string  `db:"email"`
	Password string  `db:"password"`
	Name     string  `db:"name"`
	Gender   string  `db:"gender"`
	Age      int     `db:"age"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
}
