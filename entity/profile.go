package entity

// Profile represents data about an profile.
type Profile struct {
	User    User   `json:"user"`
	ID      int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Block   string `json:"block"`
	ZipCode string `json:"zip_code"`
	City    string `json:"city"`
	State   string `json:"state"`
	UserID  int    `json:"user_id"`
}
