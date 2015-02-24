package main

import (
	"flag"
	"fmt"
	"github.com/asarcar/go_test/reddit"
	"log"
)

func main() {
	sPtr := flag.String("subreddit", "golang", "subreddit string on which to query Reddit API")
	flag.Parse()

	items, err := reddit.Get(*sPtr)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fmt.Println(item)
	}
}
