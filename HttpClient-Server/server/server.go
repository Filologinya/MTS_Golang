package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const version = "v1.0.0" // Версия API

// Структуры для обработки JSON в /decode
type DecodeRequest struct {
	InputString string `json:"inputString"`
}

type DecodeResponse struct {
	OutputString string `json:"outputString"`
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/decode", decodeHandler)
	http.HandleFunc("/hard-op", hardOpHandler)

	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

// GET /version — возвращает версию API
func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"version": version})
}

// POST /decode — принимает base64 и декодирует его
func decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		return
	}

	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	decoded, err := base64.StdEncoding.DecodeString(req.InputString)
	if err != nil {
		http.Error(w, "Неверная base64 строка", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(DecodeResponse{OutputString: string(decoded)})
}

// Генерирует случайную ошибку 5xx
func randomServerError(w http.ResponseWriter) {
	errors := map[int]string{
		500: "Internal Server Error",
		501: "Not Implemented",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Timeout",
		505: "HTTP Version Not Supported",
	}
	codes := []int{500, 501, 502, 503, 504, 505}
	code := codes[rand.Intn(len(codes))]

	w.WriteHeader(code)
	w.Write([]byte(errors[code]))
}

// GET /hard-op — случайная задержка и код ответа
func hardOpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
		return
	}

	time.Sleep(time.Duration(rand.Intn(11)+10) * time.Second) // Сон 10-20 секунд

	if rand.Intn(2) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
	} else {
		randomServerError(w) // Возвращает случайную 5xx ошибку
	}
}
