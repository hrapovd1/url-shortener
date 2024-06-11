package storage

type Storage interface {
	GetShort() (string, error)
	GetURL(string) (string, error)
	SaveURL(string, string) error
}
