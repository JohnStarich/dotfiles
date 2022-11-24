package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	err := run(os.Args[1])
	if err != nil {
		panic(err)
	}
}

func run(url string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var content string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.JavascriptAttribute("html", "innerText", &content),
	)
	if err != nil {
		return err
	}
	fmt.Println(content)
	return nil
}
