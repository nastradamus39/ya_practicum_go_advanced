package errors

import "errors"

var ErrURLConflict = errors.New(`url уже существует`)

var ErrURLNotFound = errors.New(`url не существует`)

var ErrURLDeleted = errors.New(`url удален`)

var ErrNoDBConnection = errors.New(`нет подключения к бд`)
