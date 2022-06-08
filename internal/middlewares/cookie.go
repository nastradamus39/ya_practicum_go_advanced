package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const cookieName = "ya_practicum_uuid"
const cookieMaxAge = 60 * 60 * 24 * 360
const cookieSalt = "salt"
const secret = "secret key"

var UserSignedCookie SignedCookie

func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentCookie, err := r.Cookie(cookieName)
		fmt.Printf("Текущее значение куки %s - %s\n", cookieName, currentCookie)

		// Если куки нет - создаем новую, подписываем, назначаем
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Printf("Нет куки. Создаем новую.\n")
			sc, _ := NewSignedCookie()
			sc.Sign() // подписываем
			http.SetCookie(w, sc.Cookie)
			UserSignedCookie = sc
		} else {
			// cookie есть. проверим подпись
			sc := SignedCookie{}
			sc.Cookie = currentCookie

			fmt.Printf("Валидируем существующую куку\n")

			err := sc.Validate()
			if err != nil {
				// Кука не валидна. Подписываем снова
				sc.Sign()
				http.SetCookie(w, sc.Cookie)
				fmt.Printf("Новая подписанная кука - %s\n", sc.signedValue)
			}
			UserSignedCookie = sc
		}
		next.ServeHTTP(w, r)
	})
}

// NewSignedCookie конструктор для подписанной куки
func NewSignedCookie() (_sc SignedCookie, err error) {
	_sc.salt = cookieSalt
	_sc.key = []byte("key")

	_sc.clearValue = uuid.New().String()
	_sc.Uuid = _sc.clearValue
	_sc.sign = _sc.CalcSign(_sc.clearValue)

	_sc.Cookie = &http.Cookie{
		Name:   cookieName,
		Value:  _sc.clearValue,
		Path:   "/",
		MaxAge: cookieMaxAge,
	}

	_sc.Sign()

	return _sc, nil
}

type SignedCookie struct {
	*http.Cookie
	salt        string
	key         []byte
	signedValue string
	clearValue  string
	Uuid        string
	sign        string
}

// Sign подписывает куку
func (sc *SignedCookie) Sign() (err error) {
	// Если подпись еще не высчитывали - считаем
	if sc.sign == "" {
		sc.sign = sc.CalcSign(sc.clearValue)
	}

	// Подписываем
	sc.Value = fmt.Sprintf("%s|%s", sc.clearValue, sc.sign)
	return nil
}

// CalcSign вычисляет подпись
func (sc *SignedCookie) CalcSign(cookie string) (value string) {
	secretKey := []byte(secret)
	secretKey = append(secretKey, []byte(cookie)[5:10]...)
	sc.key = secretKey

	h := hmac.New(sha256.New, sc.key)
	h.Write([]byte(cookie))

	return hex.EncodeToString(h.Sum(nil))
}

// Validate валидирует подписанную куку
func (sc *SignedCookie) Validate() (err error) {
	cookieParts := strings.Split(sc.Value, "|")

	sc.clearValue = cookieParts[0]
	sc.Uuid = sc.clearValue
	sc.sign = cookieParts[1]

	if sc.sign == sc.CalcSign(sc.clearValue) {
		fmt.Printf("Кука валидна\n")
		return nil
	} else {
		fmt.Printf("Кука не валидна\n")
		return errors.New("кука не валидна")
	}
}
