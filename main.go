package main

import (
	"flag"
	"log"
)

func main() {
	token := mustToken()

	//tgclient = telegram.New(token)
	//fetcher = fetcher.New()
	//processor = processor.New()

	//consumer.Start(fetcher, processor) consumer -получает, и обрабатывает события делают это fetcher(получает события), processor(обрабатывает)
}

func mustToken() string { // Обычно возвращают error - (string, error) , но здесь это бесссымсленно + если такого нет то добавляется префикс must
	token := flag.String(
		"tg-bot-token",
		"",
		"token to access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token isn`t specified")
	}

	return *token
}
