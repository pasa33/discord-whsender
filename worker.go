package discordwhsender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// queueCapacity is the maximum number of pending messages per webhook queue
// (normal and error queues each get their own).
const queueCapacity = 200

var httpClient = &http.Client{Timeout: 15 * time.Second}

// sender delivers messages to a single webhook URL through a background worker.
// Error reports (see errors.go) are queued separately so they aren't stuck
// behind a backlog of regular messages.
type sender struct {
	url      string
	queue    chan payload
	errQueue chan payload
}

type payload struct {
	body        []byte
	contentType string
	isError     bool
}

func (s *sender) enqueue(msg Message, isError bool) error {
	p, err := buildPayload(msg)
	if err != nil {
		return err
	}
	p.isError = isError

	ch := s.queue
	if isError {
		ch = s.errQueue
	}
	select {
	case ch <- p:
		return nil
	default:
		return fmt.Errorf("discord-whsender: queue full (%d messages pending) for %s", queueCapacity, s.url)
	}
}

func (s *sender) run() {
	for {
		var p payload
		select {
		case p = <-s.errQueue:
		default:
			select {
			case p = <-s.errQueue:
			case p = <-s.queue:
			}
		}
		s.deliver(p)
	}
}

func (s *sender) deliver(p payload) {
	for {
		res, err := httpClient.Post(s.url, p.contentType, bytes.NewReader(p.body))
		if err != nil {
			log.Printf("[discord-whsender] http error: %v — retrying in 1s (%s)", err, s.url)
			time.Sleep(time.Second)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK, http.StatusNoContent:
			remaining := rateLimitRemaining(res.Header)
			res.Body.Close()
			if remaining < 3 {
				time.Sleep(300 * time.Millisecond)
			}
			return
		case http.StatusTooManyRequests:
			delay := rateLimitDelay(res.Header)
			res.Body.Close()
			log.Printf("[discord-whsender] rate limited: %s (%s)", delay, s.url)
			time.Sleep(delay)
		default:
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()
			if !p.isError {
				sendError(s.url, res.Status, p.body, body)
			}
			return
		}
	}
}

func rateLimitRemaining(h http.Header) int {
	n, _ := strconv.Atoi(h.Get("x-ratelimit-remaining"))
	return n
}

func rateLimitDelay(h http.Header) time.Duration {
	ms, _ := strconv.Atoi(h.Get("retry-after"))
	return time.Duration(ms) * time.Millisecond
}

func buildPayload(msg Message) (payload, error) {
	if len(msg.Files) > 0 {
		return buildMultipart(msg)
	}
	return buildJSON(msg)
}

func buildJSON(msg Message) (payload, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return payload{}, err
	}
	return payload{body: data, contentType: "application/json"}, nil
}

func buildMultipart(msg Message) (payload, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	for i, f := range msg.Files {
		fw, err := w.CreateFormFile("file"+strconv.Itoa(i), f.Name)
		if err != nil {
			return payload{}, err
		}
		if _, err := fw.Write(f.Bytes); err != nil {
			return payload{}, err
		}
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return payload{}, err
	}
	if err := w.WriteField("payload_json", string(data)); err != nil {
		return payload{}, err
	}
	w.Close()

	return payload{body: buf.Bytes(), contentType: w.FormDataContentType()}, nil
}
