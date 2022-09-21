package app

type POSTPayload struct {
	LongURL  string `json:"longURL"`
	Auth0Sub string
}

type PUTPayload struct {
	ShortURL string `json:"shortURL"`
	LongURL  string `json:"longURL"`
}
