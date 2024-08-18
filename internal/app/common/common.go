// Модуль общих констант.
package common

import "errors"

// ContextValueKey тип ключа контекста.
type ContextValueKey int

// KeyUserID ключ ID пользователя для контекста.
const KeyUserID ContextValueKey = iota

// ErrFetchUserIDFromContext ошибка получеения ID пользователя из контекста
var ErrFetchUserIDFromContext = errors.New("failed to fetch user id from context")

// ErrPermDenied ошибка прав доступа к ссылке.
var ErrPermDenied = errors.New("permission denied for url")

// EncRespErrStr ошибка кодировки ответа.
var EncRespErrStr = "error encoding response"

// ReadReqErrStr ошибка чтения тела запроса.
var ReadReqErrStr = "failed to read request body"

// ContentTypeHeader тип контента.
var ContentTypeHeader = "Content-Type"

// JSONContentType тип контента JSON.
var JSONContentType = "application/json"
