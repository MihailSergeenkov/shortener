package common

import "errors"

type ContextValueKey int

const KeyUserID ContextValueKey = iota

var ErrFetchUserIDFromContext = errors.New("failed to fetch user id from context")
var ErrPermDenied = errors.New("permission denied for url")
var EncRespErrStr = "error encoding response"
var ReadReqErrStr = "failed to read request body"
var ContentTypeHeader = "Content-Type"
var JSONContentType = "application/json"
