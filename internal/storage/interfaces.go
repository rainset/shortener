package storage

type BatchUrls struct {
	CorrelationID string
	OriginalURL   string
}

type ResultURL struct {
	ID       int
	Hash     string
	Original string
	Deleted  int
}

type ResultBatchUrls struct {
	CorrelationID string
	Hash          string
}

type ResultHistoryURL struct {
	ID       int
	CookieID string
	Hash     string
	Original string
}

type InterfaceStorage interface {
	AddURL(u, original string) (err error)
	GetURL(hash string) (resultURL ResultURL, err error)
	GetByOriginalURL(original string) (hash string, err error)
	AddBatchURL([]BatchUrls) ([]ResultBatchUrls, error)
	DeleteUserBatchURL(cookieID string, hashes []string) error
	DeleteBatchURL(hashes []string) error
	AddUserHistoryURL(cookieID, hash string) (err error)
	Ping() error
	GetListUserHistoryURL(cookieID string) ([]ResultHistoryURL, error)
}
