package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFetchURL(t *testing.T) {
	var wg sync.WaitGroup
	respList := make(map[string]string)

	// Создаем тестовый HTTP-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Тест"))
	}))
	defer ts.Close()

	// Определяем URL для тестирования
	testURL := ts.URL

	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем тестируемую функцию
	go fetchURL(ctx, &wg, testURL, respList)
	wg.Wait()

	// Проверяем результат
	expected := "Тест"
	if respList[testURL] != expected {
		t.Errorf("Ожидал ответ %q, но получен: %q", expected, respList[testURL])
	}
}

func TestFetchURL_Error(t *testing.T) {
	var wg sync.WaitGroup
	respList := make(map[string]string)

	// Не существующий URL для проверки обработки ошибок
	testURL := "http://aaasssdddsssaassss.url"

	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go fetchURL(ctx, &wg, testURL, respList)
	wg.Wait()

	// Проверяем, что в списке ответов есть ошибка
	if !strings.Contains(respList[testURL], "Ошибка запроса по адресу") {
		t.Errorf("Ожидалась ошибка %q но получен: %q", testURL, respList[testURL])
	}
}
