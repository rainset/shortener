// Пакет для работы с базами данных
package storage

// BatchUrls структура для массового добавления ссылок
type BatchUrls struct {
	CorrelationID string
	OriginalURL   string
}

// ResultURL структура ответа из базы данных
type ResultURL struct {
	ID       int
	Hash     string
	Original string
	Deleted  int
}

// ResultBatchUrls структура ответа из базы данных при массовом добавлении ссылок
type ResultBatchUrls struct {
	CorrelationID string
	Hash          string
}

// ResultHistoryURL структура ответа из базы данных истории добюавленных ссылок пользователя
type ResultHistoryURL struct {
	ID       int
	CookieID string
	Hash     string
	Original string
}

// Stats статистика
type Stats struct {
	Urls  int
	Users int
}

// InterfaceStorage интерфейс для взаимодействия с базой данных
type InterfaceStorage interface {
	//AddURL Добавление ссылки
	AddURL(u, original string) (err error)
	//GetURL Получение ссылки
	GetURL(hash string) (resultURL ResultURL, err error)
	//GetByOriginalURL Получение хеша короткой ссылки по ее значению
	GetByOriginalURL(original string) (hash string, err error)
	//AddBatchURL Массовое добавление ссылок
	AddBatchURL([]BatchUrls) ([]ResultBatchUrls, error)
	//DeleteUserBatchURL Удаление ссылки по cookieID пользователя
	DeleteUserBatchURL(cookieID string, hashes []string) error
	//DeleteBatchURL Массовое удаление ссылок
	DeleteBatchURL(hashes []string) error
	//AddUserHistoryURL Добавление ссылки пользователя в историю
	AddUserHistoryURL(cookieID, hash string) (err error)
	//Ping Проверка работы бд
	Ping() error
	//GetListUserHistoryURL Получение списка ссылок пользователя
	GetListUserHistoryURL(cookieID string) ([]ResultHistoryURL, error)
	//GetStats Получение статистики сервиса
	GetStats() (Stats, error)
}
