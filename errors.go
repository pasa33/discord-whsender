package discordwhsender

import "encoding/base64"

func sendError(url string, status string, req, res []byte) error {
	s, found := senders.LoadOrStore(url, newSender(url))
	sender := s.(*sender)
	if !found {
		sender.initSender()
	}
	if err := sender.queueAdd(makeErrorMsg(status, req, res), true); err != nil {
		return err
	}
	return nil
}

func makeErrorMsg(status string, req, res []byte) Message {
	return Message{
		Username: "whsender-error",
		Content:  status,
		Files: []File{
			{Name: "ReqPayload.txt", Bytes: []byte(base64.StdEncoding.EncodeToString(req))},
			{Name: "ResPayload.txt", Bytes: []byte(base64.StdEncoding.EncodeToString(res))},
		},
	}
}
