package discordwhsender

import (
	"bytes"
	"cmp"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	senders  cmap.ConcurrentMap[string, *sender]
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
	errUrl   string
	debugUrl string
)

type sender struct {
	WhUrl   string
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
// TODO: implement mergeEmbeds for reduce ratelimit
func (msg Message) Send(url string, mergeEmbeds ...bool) error {
	sender := getSender(cmp.Or(debugUrl, url))
	msg.validate()
	return sender.queueAdd(msg, false)
}

// Set global error webhook url
// for unset, just set to empty string
func SetErrorWh(url string) {
	errUrl = url
}

// Set debug webhook
// that override every whs
func SetDebugWh(url string) {
	debugUrl = url
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

func getSender(url string) *sender {
	s, found := senders.Get(url)
	if s == nil {
		s = newSender(url)
		senders.Set(url, s)
	}
	if !found {
		s.initSender()
	}
	return s
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
					case 200, 204:
						rtRemaining, _ := strconv.Atoi(res.Header.Get("x-ratelimit-remaining"))
						if rtRemaining < 3 {
							time.Sleep(300 * time.Millisecond)
						}
						retry = false
					case 429:
						ratelimitDelay, _ := strconv.Atoi(res.Header.Get("retry-after"))
						log.Printf("[discord-whsender] Ratelimited: %dms (%s)\n", ratelimitDelay, s.WhUrl)
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
