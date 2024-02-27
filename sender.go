package discordwhsender

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	senders sync.Map
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
)

type sender struct {
	WhUrl   string
	ErrUrl  string
	Queue   []msgPayload
	Mu      *sync.Mutex
	Waiter  *sync.WaitGroup
	Waiting bool
}

type msgPayload struct {
	Bytes       []byte
	ContentType string
	IsError     bool
}

// Send a message to a specific discord webhook url
func (msg Message) Send(url string, mergeEmbeds ...bool) error {
	s, found := senders.LoadOrStore(url, newSender(url))
	sender := s.(*sender)
	if !found {
		sender.initSender()
	}
	if err := sender.queueAdd(msg, false); err != nil {
		return err
	}
	return nil
}

func newSender(url string) *sender {
	return &sender{
		WhUrl:   url,
		Queue:   []msgPayload{},
		Mu:      &sync.Mutex{},
		Waiter:  &sync.WaitGroup{},
		Waiting: false,
	}
}

func (s *sender) initSender() {
	go func() {
		for {
			s.Waiter.Wait()
			if p := s.queueGet(); len(p.Bytes) > 0 {
				retry := true
				for retry {
					res, err := http.Post(s.WhUrl, p.ContentType, bytes.NewBuffer(p.Bytes))
					if err != nil {
						continue
					}

					switch res.StatusCode {
					case 204:
						rtRemaining, _ := strconv.Atoi(res.Header.Get("x-ratelimit-remaining"))
						if rtRemaining < 3 {
							time.Sleep(300 * time.Millisecond)
						}
						retry = false
					case 429:
						ratelimitDelay, _ := strconv.Atoi(res.Header.Get("retry-after"))
						fmt.Println("WH Ratelimited for: ", ratelimitDelay)
						time.Sleep(time.Duration(ratelimitDelay) * time.Millisecond)
						retry = true
					default:
						if !p.IsError {
							bbody, _ := io.ReadAll(res.Body)
							sendError(s.WhUrl, res.Status, p.Bytes, bbody)
						}
						retry = false
					}
					res.Body.Close()
				}
			}
		}
	}()
}
