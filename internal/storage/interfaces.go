package storage

type BatchUrls struct {
	CorrelationID string
	OriginalURL   string
}

type ResultBatchUrls struct {
	CorrelationID string
	Hash          string
}

type ResultHistoryUrl struct {
	ID       int
	CookieID string
	Hash     string
	Original string
}

type InterfaceStorage interface {
	AddURL(u, original string) (err error)
	GetURL(hash string) (originalURL string, err error)
	GetByOriginalURL(original string) (hash string, err error)
	AddBatchURL([]BatchUrls) ([]ResultBatchUrls, error)
	AddUserHistoryURL(cookieID, hash string) (err error)
	Ping() error
	GetListUserHistoryURL(cookieID string) ([]ResultHistoryUrl, error)
}
