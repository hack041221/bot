package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func NewDownloader(client *http.Client, threads int) *Downloader {
	return &Downloader{client, threads}
}

type Downloader struct {
	client  *http.Client
	threads int
}

func (dl *Downloader) download(url, filename string) error {
	startTime := time.Now()
	l.Info().Msgf("Going to download %s from %s", filename, url)
	resp, err := dl.client.Head(url)
	if err != nil {
		return err
	}
	contentLength := resp.Header.Get("Content-Length")
	ranges := resp.Header.Get("Accept-Ranges")
	if contentLength == "" {
		l.Info().Msgf("Content length not specified")
		return dl.singleThreaded(url, filename, startTime)
	}
	if ranges != "bytes" {
		l.Info().Msgf("Server does not accept byte ranges")
		return dl.singleThreaded(url, filename, startTime)
	}
	if dl.threads <= 1 {
		return dl.singleThreaded(url, filename, startTime)
	}
	length, err := strconv.Atoi(contentLength)
	if err != nil {
		return err
	}
	return dl.threadedDownload(url, filename, length, startTime)
}

func (dl *Downloader) singleThreaded(url, filename string, startTime time.Time) error {
	l.Info().Msgf("Starting single threaded download...")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := dl.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	l.Info().Msgf("Writing to %s", filename)
	if _, err := file.Write(body); err != nil {
		return err
	}
	l.Info().Msgf("Downloaded %d bytes to %s in %s", len(body), filename, time.Now().Sub(startTime))
	if err := resp.Body.Close(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func (dl *Downloader) threadedDownload(url string, filename string, length int, startTime time.Time) error {
	l.Info().Msgf("Starting threaded download...")
	size := length / dl.threads
	remainder := length % dl.threads
	l.Info().Msgf("Downloading %s on %d threads", filename, dl.threads)
	wg := &sync.WaitGroup{}
	for i := 0; i < dl.threads; i++ {
		wg.Add(1)

		start := i * size
		end := (i + 1) * size

		if i == dl.threads-1 {
			end += remainder
		}

		l.Info().Msgf("Starting thread %d", i)
		go func(start, end, i int) error {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				wg.Done()
				return err
			}
			byteRange := fmt.Sprintf("bytes=%d-%d", start, end-1)
			req.Header.Add("Range", byteRange)
			resp, err := dl.client.Do(req)
			if err != nil {
				wg.Done()
				return err
			}
			defer resp.Body.Close()
			l.Info().Msgf("Thread: %d Reading response body", i)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				wg.Done()
				return err
			}
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				wg.Done()
				return err
			}
			defer file.Close()
			if _, err := io.Copy(file, resp.Body); err != nil {
				wg.Done()
				return err
			}
			l.Info().Msgf("Thread: %d writing bytes %d - %d", i, start, end)
			if _, err := file.WriteAt(body, int64(start)); err != nil {
				wg.Done()
				return err
			}
			wg.Done()
			l.Info().Msgf("Thread: %d done", i)
			return nil
		}(start, end, i)
	}
	wg.Wait()
	l.Info().Msgf("Downloaded %s in %s", filename, time.Now().Sub(startTime))
	return nil
}
