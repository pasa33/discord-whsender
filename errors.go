package discordwhsender

import (
	"cmp"
	"encoding/base64"
	"fmt"
)

func sendError(url string, status string, req, res []byte) error {
	sender := getSender(cmp.Or(debugUrl, errUrl, url))
	return sender.queueAdd(makeErrorMsg(status, url, req, res), true)
}

func makeErrorMsg(status, url string, req, res []byte) Message {
	c := status
	if len(errUrl) > 0 {
		c += fmt.Sprintf("\n`%s`", url)
	}
	return Message{
		Username: "whsender-error",
		Content:  c,
		Files: []File{
			{Name: "ReqPayload.txt", Bytes: []byte(base64.StdEncoding.EncodeToString(req))},
			{Name: "ResPayload.txt", Bytes: []byte(base64.StdEncoding.EncodeToString(res))},
		},
	}
}
