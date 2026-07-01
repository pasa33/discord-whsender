# discord-whsender

Minimal Go package for sending messages to Discord webhooks from scrapers and monitors.

Messages are delivered asynchronously via a per-webhook queue. Rate limits (HTTP 429) are
handled automatically using Discord's `retry-after` header.

## Install

```
go get github.com/pasa33/discord-whsender
```

## Usage

```go
import wh "github.com/pasa33/discord-whsender"

msg := wh.Message{
    Username: "my-bot",
    Content:  "New item found",
    Embeds: []*wh.Embed{
        {
            Title: "Item",
            Color: 0x00ff00,
        },
    },
}
msg.Embeds[0].AddField("Name", itemName, true).AddField("Price", "$99", true)

if err := msg.Send(webhookURL); err != nil {
    log.Println(err)
}
```

`Send` returns an error only if the message fails to encode or its per-webhook queue is
full — actual delivery (including retries) happens in the background.

### Attach a file

```go
msg := wh.Message{
    Content: "log attached",
    Files: []wh.File{
        {Name: "log.txt", Bytes: logBytes},
    },
}
msg.Send(webhookURL)
```

### Error reporting

Redirect failed deliveries (non-2xx/429 responses) to a dedicated webhook. The request and
response bodies are attached as files. If unset, error reports fall back to the same webhook
that failed.

```go
wh.SetErrorWh(errorWebhookURL)

// disable
wh.SetErrorWh("")
```

### Debug mode

Redirect every `Send` call to a single webhook — useful during development so you don't spam
production channels.

```go
wh.SetDebugWh(debugWebhookURL)

// disable
wh.SetDebugWh("")
```

### Muting

Disable all sends, e.g. in tests or local dev:

```go
wh.SetMuted(true)
```

## Notes

- Each webhook URL gets its own queue (capacity 200) and background worker; `Send` returns an
  error if the queue is full.
- Embed field names/values are validated automatically: empty strings become `"-"`, and values
  exceeding Discord's limits (256/1024 chars) are truncated.
