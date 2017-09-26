package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
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

// Other funcs, if any.

func main() {
	kworkers := flag.Uint("k", 5, "limits the number of workers, > 0")
	flag.Parse()
	if *kworkers == 0 {
		log.Printf("number of workers must be greater than 0!")
		flag.PrintDefaults()
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	limit := newLimiter(*kworkers)
	var total uint64 = 0

	for scanner.Scan() {
		go processLink(scanner.Text(), &total, limit.Start())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	limit.Wait()
	log.Printf("Total: %d", total)
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

func processLink(link string, total *uint64, done func()) (int, error) {
	defer done()
	buff := &bytes.Buffer{}
	err := download(link, buff)
	if err != nil {
		return -1, fmt.Errorf("error while downloading data from %q: %s", link, err)
	}
	n := bytes.Count(buff.Bytes(), []byte("Go"))
	log.Printf("Count for %s: %d\n", link, n)
	atomic.AddUint64(total, uint64(n))
	return n, nil
}
