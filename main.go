package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var URL string
	var timeOutInSecond int

	fmt.Println("Введите URL адреса через запятую:")
	fmt.Scan(&URL)

	if URL == "" {
		fmt.Println("URL адреса не введены")
		return
	}
	URLList := strings.Split(URL, ",")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Введите период опроса:")
	fmt.Scan(&timeOutInSecond)
	timeOut := time.Duration(timeOutInSecond) * time.Second
	ticker := time.NewTicker(timeOut)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeOut)
			defer cancel()

			respList := make(map[string]string)

			for _, URL := range URLList {
				wg.Add(1)
				go fetchURL(ctx, &wg, URL, respList)
			}

			wg.Wait()

			fmt.Println(time.Now())
			for _, val := range respList {
				fmt.Println(val)
			}

		case <-stop:
			fmt.Println("Отмена")
			return
		}
	}
}

// Выполняет HTTP-запрос по указанному URL и сохраняет результат в respList.
func fetchURL(ctx context.Context, wg *sync.WaitGroup, URL string, respList map[string]string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		respList[URL] = "Ошибка при создании запроса: " + err.Error()
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			respList[URL] = "Запрос по адресу был отменён из-за тайм-аута: " + URL
		} else {
			respList[URL] = "Ошибка запроса по адресу: " + err.Error()
		}
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		respList[URL] = "Ошибка чтения ответа с адреса: " + err.Error()
		return
	}

	respList[URL] = string(body)
}
