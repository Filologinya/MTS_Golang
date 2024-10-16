package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Структура для ответа /decode
type DecodeResponse struct {
	OutputString string `json:"outputString"`
}

// Структура для ответа /version
type VersionResponse struct {
  Version string `json:"version"`
}

func main() {
	baseURL := "http://localhost:8080"

	callVersion(baseURL + "/version")

	callDecode(baseURL + "/decode", "SGVsbG8gd29ybGQh")

	callHardOpWithTimeout(baseURL + "/hard-op", 15*time.Second)
}

func callVersion(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при вызове /version:", err)
		return
	}
	defer resp.Body.Close()

	var versionResponse VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResponse); err != nil {
		fmt.Println("Ошибка при обработке ответа /version:", err)
		return
	}
	fmt.Printf("%s\n", versionResponse.Version)
}

func callDecode(url, input string) {
  // Создаём JSON ответ с исходной строкой
	reqBody, _ := json.Marshal(map[string]string{"inputString": input})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Ошибка при вызове /decode:", err)
		return
	}
	defer resp.Body.Close()

	var decodeResp DecodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&decodeResp); err != nil {
		fmt.Println("Ошибка при обработке ответа /decode:", err)
		return
	}

	fmt.Printf("%s\n", decodeResp.OutputString)
}

func callHardOpWithTimeout(url string, timeout time.Duration) {
	// Создаём контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса /hard-op:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		// Проверяем связана ли ошибка с таймаутом
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Ошибка: время ожидания ответа истекло /hard-op")
		} else {
      fmt.Println("Ошибка при вызове /hard-op:", err)
		}
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s, %d\n", body, resp.StatusCode)
}

