// Модуль моделей сервиса.
package models

// Request модель запроса короткой ссылки для оригинальной.
type Request struct {
	URL string `json:"url"`
}

// Response модель ответа на запрос короткой ссылки.
type Response struct {
	Result string `json:"result"`
}

// URL модель ссылки.
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	ID          uint   `json:"id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// BatchRequest модель запроса можественного получения коротких ссылок.
type BatchRequest []BatchDataRequest

// BatchDataRequest модель конкретного запроса в составе множественного.
type BatchDataRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse модель ответа на множественное получение коротких ссылок.
type BatchResponse []BatchDataResponse

// BatchDataResponse модель конкретного ответа в состааве множественного.
type BatchDataResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURLsResponse модель ответа на получение всех пользовательских ссылок.
type UserURLsResponse []UserURLsDataResponse

// UserURLsDataResponse модель конкретной пользовательской ссылки.
type UserURLsDataResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// StatsResponse модель статистических данных.
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
