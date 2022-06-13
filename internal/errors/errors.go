package errors

import "errors"

var UrlConflict = errors.New(`url уже существует`)

var NoDbConnection = errors.New(`нет подключения к бд`)
