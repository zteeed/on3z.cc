package swagger

import "fmt"

type ShortUrlMap struct {
	LongURL  string
	ShortURL string
}

func (s ShortUrlMap) String() string {
	return fmt.Sprintf("Urls<%s>", s.ShortURL)
}
