// Package discordwhsender sends messages to Discord webhooks asynchronously,
// with automatic rate-limit handling and per-webhook delivery queues.
package discordwhsender

import (
	"cmp"
	"sync"
)

var (
	senders   sync.Map // string(webhook URL) -> *sender
	sendersMu sync.Mutex

	errURL   string
	debugURL string
	muted    bool
)

// Send queues msg for asynchronous delivery to the given Discord webhook URL.
// Delivery happens on a per-webhook background worker that retries automatically
// on rate limits (HTTP 429).
func (msg Message) Send(url string) error {
	if muted {
		return nil
	}
	msg.validate()
	return getSender(cmp.Or(debugURL, url)).enqueue(msg, false)
}

// SetErrorWh sets the webhook URL used to report delivery failures (request and
// response payloads are attached as files). Pass an empty string to disable.
func SetErrorWh(url string) {
	errURL = url
}

// SetDebugWh redirects every Send call to the given webhook URL, regardless of
// the URL passed to Send. Pass an empty string to clear it.
func SetDebugWh(url string) {
	debugURL = url
}

// SetMuted disables (or re-enables) all sends. Useful in tests or local dev so
// messages aren't posted to real webhooks.
func SetMuted(m bool) {
	muted = m
}

func getSender(url string) *sender {
	if v, ok := senders.Load(url); ok {
		return v.(*sender)
	}
	sendersMu.Lock()
	defer sendersMu.Unlock()
	if v, ok := senders.Load(url); ok {
		return v.(*sender)
	}
	s := &sender{
		url:      url,
		queue:    make(chan payload, queueCapacity),
		errQueue: make(chan payload, queueCapacity),
	}
	senders.Store(url, s)
	go s.run()
	return s
}
