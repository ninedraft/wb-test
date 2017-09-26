package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFlags(0)
}

// Other funcs, if any.

func main() {

	// Initialization, if any.

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		// Do something with strings here.

	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// Other code, if any.

	total := 0

	// Other code, if any.

	log.Printf("Total: %v", total)
}

func download(link string, wr io.Writer) error {
	resp, err := (&http.Client{
		Timeout: 4 * time.Second,
	}).Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(wr, resp.Body)
	return err
}
