package entity

type User struct {
	ID             int      `json:"id"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	Name           string   `json:"name"`
	Gender         string   `json:"gender"`
	Age            int      `json:"age"`
	Location       Location `json:"location"`
	DistanceFromMe float64  `json:"distanceFromMe,omitempty"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
