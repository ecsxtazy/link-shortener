package handler

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Short string `json:"short"`
}
