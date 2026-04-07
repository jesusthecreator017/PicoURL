package api

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

type statsResponse struct {
	ShortURL   string `json:"short_url"`
	ClickCount int    `json:"click_count"`
}

type errorResponse struct {
	Error string `json:"error"`
}
