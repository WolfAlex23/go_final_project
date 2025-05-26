package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type PasswordInput struct {
	Password string `json:"password"`
}

func makeHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

var secret = []byte("my_secret_key")

func authHandler(w http.ResponseWriter, r *http.Request) {

	var input PasswordInput
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка десериализации JSON: %v", err)})
		return
	}
	inputPassword := input.Password

	if len(inputPassword) == 0 {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Пароль не введен"})
		return

	}

	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Аутентификация отключена"})
		return
	}
	if inputPassword != envPassword {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	} else {

		hashString := makeHash(inputPassword)

		// создаём payload
		claims := jwt.MapClaims{
			"hash": hashString,
		}

		// создаём jwt и указываем payload
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// получаем подписанный токен
		signedToken, err := jwtToken.SignedString(secret)
		if err != nil {
			writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Ошибка создания токена"})
			return
		}
		writeJson(w, http.StatusOK, map[string]string{"token": signedToken})
	}

}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var token string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}

			// здесь код для валидации и проверки JWT-токена
			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				http.Error(w, "failed to parse token", http.StatusInternalServerError)
				return
			}

			if !jwtToken.Valid {

				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			// обязательно используем второе возвращаемое значение ok и проверяем его, потому что
			// если Сlaims вдруг окажется другого типа, мы получим панику
			if !ok {
				http.Error(w, "Failed to extract claims from the token", http.StatusInternalServerError)
				return
			}

			hash := claims["hash"]

			hashStr, ok := hash.(string)
			if !ok {
				http.Error(w, "failed to typecast to string", http.StatusInternalServerError)
				return
			}

			checkHash := makeHash(pass)

			if checkHash != hashStr {

				http.Error(w, "missmatch hash in token", http.StatusUnauthorized)
				return

			}
		}
		next(w, r)
	})
}
