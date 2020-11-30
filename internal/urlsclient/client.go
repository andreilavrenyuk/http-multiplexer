package urlsclient

import (
	"context"
	"encoding/json"
	"http_multiplexer/internal/config"
	"net/http"
	"sync"
)

type urlsResult struct {
	sync.Mutex
	data map[string]interface{}
}

func (r *urlsResult) add(url string, data interface{}) {
	r.Lock()
	r.data[url] = data
	r.Unlock()
}

func Get(ctx context.Context, urls []string, parallelRequests int) (map[string]interface{}, error) {
	requests := make(chan struct{}, parallelRequests)

	result := &urlsResult{
		data: make(map[string]interface{}, len(urls)),
	}

	ctx, cancel := context.WithCancel(ctx)

	var (
		wg  sync.WaitGroup
		err error
	)

	for _, url := range urls {
		if err != nil {
			break
		}

		requests <- struct{}{}

		wg.Add(1)

		go func(url string) {
			reqCtx, reqCancel := context.WithTimeout(ctx, config.RequestTimeout())

			defer func() {
				reqCancel()
				<-requests
				wg.Done()
			}()

			data, e := get(reqCtx, url)
			if e != nil {
				err = e

				cancel()

				return
			}

			result.add(url, data)
		}(url)
	}

	wg.Wait()

	cancel()

	if err != nil {
		return nil, err
	}

	return result.data, nil
}

func get(ctx context.Context, url string) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body := map[string]interface{}{}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
