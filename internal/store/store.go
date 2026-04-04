package store

type Store interface {
	SaveURL(shortURL, originalURL string) error
	GetOriginalURL(shortURL string) (string, error)
	IncrementCount(shortURL string) error
	GetCount(shortURL string) (int, error)
	DeleteURL(shortURL string) error
	Close() error
}
