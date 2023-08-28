package downloader

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultRetries           = 3
	defaultRetryDelaySeconds = 5
	defaultUserAgent         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
)

type File struct {
	URL        string
	LocalPath  string
	retries    int
	retryDelay time.Duration
}

func NewFile(URL, localPath string) File {
	return File{
		URL:        URL,
		LocalPath:  localPath,
		retries:    defaultRetries,
		retryDelay: defaultRetryDelaySeconds * time.Second,
	}
}

func (fd File) Exists() bool {
	_, err := os.Stat(fd.LocalPath)
	return err == nil || errors.Is(err, os.ErrExist)
}

func (fd File) downloadIncomplete(client *http.Client, incompleteLocalPath string) error {
	var resp *http.Response
	var err error
	req, err := http.NewRequest("GET", fd.URL, nil)
	if err != nil {
		return errors.Join(fmt.Errorf("creating HTTP request %s", fd.URL), err)
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	resp, err = client.Do(req)
	if err != nil {
		return errors.Join(fmt.Errorf("making HTTP request to %s", fd.URL), err)
	}
	defer resp.Body.Close()
	_ = os.Remove(incompleteLocalPath)
	err = os.MkdirAll(filepath.Dir(incompleteLocalPath), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(incompleteLocalPath)
	if err != nil {
		return errors.Join(errors.New("creating file "+incompleteLocalPath), err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Join(errors.New("downloading file "+incompleteLocalPath), err)
	}
	return nil
}
func (fd File) Download(client *http.Client) error {
	var err error
	if fd.Exists() {
		slog.Debug("File exists", "local_path", fd.LocalPath)
		return nil
	}
	incompleteLocalPath := fmt.Sprintf("%s.incomplete", fd.LocalPath)
	defer os.Remove(incompleteLocalPath)
	for tries := 1; tries <= fd.retries; tries++ {
		err = fd.downloadIncomplete(client, incompleteLocalPath)
		if err == nil {
			break
		}
		time.Sleep(fd.retryDelay * time.Duration(tries))
	}
	if err == nil {
		err = os.Rename(incompleteLocalPath, fd.LocalPath)
		if err != nil {
			slog.Debug("Downloaded", "local_path", fd.LocalPath)
		}
	}
	return err
}

type Options struct {
	ProxyURL         string
	Timeout          int
	Threads          int
	WarnAsErrorLimit int
}

type OptionFunc func(*Options)

func WithProxy(proxyURL string) OptionFunc {
	return func(o *Options) {
		o.ProxyURL = proxyURL
	}
}

func WithTimeout(timeout int) OptionFunc {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func WithThreads(threads int) OptionFunc {
	return func(o *Options) {
		o.Threads = threads
	}
}

func WithWarnAsErrorLimit(limit int) OptionFunc {
	return func(o *Options) {
		o.WarnAsErrorLimit = limit
	}

}

type Downloader struct {
	ch               chan File
	wg               *sync.WaitGroup
	client           *http.Client
	threads          int
	warnAsErrorLimit int
}

func New(options ...OptionFunc) *Downloader {
	do := &Options{}
	for _, fo := range options {
		fo(do)
	}

	return &Downloader{
		ch:      make(chan File),
		wg:      new(sync.WaitGroup),
		client:  createHttpClient(do),
		threads: max(1, do.Threads),
	}
}

func createHttpClient(options *Options) *http.Client {
	client := &http.Client{}
	if options.ProxyURL != "" {
		proxyURL, _ := url.Parse(options.ProxyURL)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}
	if options.Timeout > 0 {
		client.Timeout = time.Duration(options.Timeout) * time.Second
	}
	return client
}

func (dl *Downloader) Start() {
	for t := 0; t < dl.threads; t++ {
		dl.wg.Add(1)
		go func() {
			defer dl.wg.Done()
			continuousWarnings := 0
			for di := range dl.ch {
				err := di.Download(dl.client)
				if err != nil {
					slog.Warn("Downloading", "error", err)
					continuousWarnings += 1
					if dl.warnAsErrorLimit > 0 && continuousWarnings >= dl.warnAsErrorLimit {
						return
					}
				} else {
					continuousWarnings = 0
				}
			}
		}()
	}
}

func (dl *Downloader) Send(di File) {
	dl.ch <- di
}

func (dl *Downloader) Close() {
	close(dl.ch)
}

func (dl *Downloader) Wait() {
	dl.Close()
	dl.wg.Wait()
}
