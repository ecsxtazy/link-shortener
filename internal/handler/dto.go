package handler

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type ResolveResponse struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
