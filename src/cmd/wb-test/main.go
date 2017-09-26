package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func init() {
	log.SetFlags(0)
}

func main() {
	kworkers := flag.Uint("k", 5, "limits the number of workers, > 0")
	tlinks := flag.String("e", "", "list of links, separated by \\n")
	flag.Parse()
	if *kworkers == 0 {
		log.Printf("number of workers must be greater than 0!")
		flag.PrintDefaults()
		os.Exit(1)
	}
	limit := newLimiter(*kworkers)
	var scanner *bufio.Scanner
	// if list of links is not provided as flags,
	// scanner reads from stdio
	if *tlinks == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		scanner = bufio.NewScanner(bytes.NewReader([]byte(*tlinks)))
	}
	var total uint64 = 0
	for scanner.Scan() {
		go processLink(scanner.Text(), &total, limit.Start())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error while scanning text: %s\n", err)
	}
	limit.Wait()
	log.Printf("Total: %d", total)
}

func download(link string, wr io.Writer) error {
	resp, err := (&http.Client{
		Timeout: 10 * time.Second,
	}).Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(wr, resp.Body)
	return err
}

func processLink(link string, total *uint64, done func()) {
	defer done()
	buff := &bytes.Buffer{}
	err := download(link, buff)
	if err != nil {
		log.Printf("error while downloading data from %q: %s\n", link, err)
		return
	}
	n := bytes.Count(buff.Bytes(), []byte("Go"))
	log.Printf("Count for %s: %d\n", link, n)
	atomic.AddUint64(total, uint64(n))
}
