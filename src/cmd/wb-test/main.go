package main

import (
	"bufio"
	"bytes"
	"fmt"
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

func processLink(link string, limiter chan struct{}) (int, error) {
	defer func() {
		select {
		case <-limiter:
		}
	}()
	buff := &bytes.Buffer{}
	err := download(link, buff)
	if err != nil {
		return -1, fmt.Errorf("error while downloading data from %q: %s", link, err)
	}
	return bytes.Count(buff.Bytes(), []byte("Go")), nil
}
