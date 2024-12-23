package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	const timeOut = 5 * time.Second

	fmt.Println("Введите URL адреса через запятую")
	fmt.Scan(&URL)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	if URL == "" {
		fmt.Println("URL адреса не введены")
		return
	}
	URLList := strings.Split(URL, ",")

	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()
		respList := make(map[string]string)

		for _, URL := range URLList {
			wg.Add(1)
			go func(URL string) {
				defer wg.Done()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
				if err != nil {
					//fmt.Println("Ошибка при создании запроса:", err)
					return
				}
				client := &http.Client{}
				response, err := client.Do(req)
				if err != nil {
					//fmt.Println("Ошибка запроса по адресу:" + URL)
					return
				}
				defer response.Body.Close()
				body, err := ioutil.ReadAll(response.Body)
				respList[URL] = string(body)
			}(URL)
		}
		<-ctx.Done()

		fmt.Println(time.Now())
		for _, val := range respList {
			fmt.Println(val)
		}
		<-stop
		//time.Sleep(timeOut)
	}
}
