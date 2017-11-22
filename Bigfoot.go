package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sclevine/agouti"
)

// need: web driver multithread error
const acyncWait = 300

// Bigfoot is clawer
type Bigfoot struct {
	multi  bool
	driver *agouti.WebDriver
	urls   []string
	par    int
	sec    time.Duration
}

// NewBigfoot retuns Bigfoot and error
func NewBigfoot(name string, mul bool, par int, sec time.Duration) (*Bigfoot, error) {
	urls, err := getURL(name)
	if err != nil {
		return nil, err
	}

	return &Bigfoot{
		multi:  mul,
		driver: agouti.ChromeDriver(),
		urls:   urls,
		par:    par,
		sec:    sec,
	}, nil
}

func getURL(name string) ([]string, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("No Such File: %s", name)
	}
	defer fp.Close()
	reader := csv.NewReader(fp)

	urls := []string{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(row) != 3 {
			return nil, errors.New("CSV Error format")
		}
		urls = append(urls, row[0])
	}

	return urls, nil
}

// Start invokes web driver
func (it *Bigfoot) Start() error {
	return it.driver.Start()
}

// Stop stops web driver
func (it *Bigfoot) Stop() error {
	return it.driver.Stop()
}

func (it *Bigfoot) groupRun(urls []string) error {
	var wg sync.WaitGroup

	var pages []*agouti.Page
	for _, URL := range urls {
		p, err := it.driver.NewPage()
		pages = append(pages, p)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(p *agouti.Page, URL string) {
			defer wg.Done()
			err := p.Navigate(URL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "groupRun err:%v", err)
			}
		}(p, URL)

		// need : protect web driver multithreading error
		time.Sleep(acyncWait * time.Millisecond)
	}
	wg.Wait()
	time.Sleep(it.sec)

	for _, p := range pages {
		p.Destroy()
	}
	return nil
}

// multiRun open many browser
func (it *Bigfoot) multiRun() error {
	count := 0
	l := len(it.urls)
	urls := it.urls

	low := 0
	high := it.par

	for {
		if low > l {
			break
		}
		if high > l {
			high = l
		}
		it.groupRun(urls[low:high])

		count++
		low = it.par * count
		high = low + it.par
	}

	return nil
}

// Run invoke chrome browser
func (it *Bigfoot) Run() error {
	if it.multi {
		return it.multiRun()
	}

	return it.run()
}

// Run is walked web driver par sec by urls
func (it *Bigfoot) run() error {
	p, err := it.driver.NewPage()
	if err != nil {
		return err
	}
	for _, URL := range it.urls {
		err = p.Navigate(URL)
		if err != nil {
			return err
		}
		time.Sleep(it.sec)
	}

	return nil
}
