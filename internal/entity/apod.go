package entity

type APOD struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Date        string `json:"date"`
	Url         string `json:"url,omitempty"`
	ImageB64    string `json:"image_b64"`
}
