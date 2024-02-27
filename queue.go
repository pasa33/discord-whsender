package discordwhsender

import (
	"bytes"
	"mime/multipart"
	"strconv"
)

func (s *Sender) queueGet() (p MsgPayload) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if len(s.Queue) > 0 {
		p = s.Queue[0]
		if len(s.Queue) > 1 {
			s.Queue = s.Queue[1:]
		} else {
			s.Queue = []MsgPayload{}
		}
		return
	}
	s.Waiting = true
	s.Waiter.Add(1)
	return
}

func (s *Sender) queueAdd(wh Message, isErr bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	p := MsgPayload{
		IsError: isErr,
	}

	if len(wh.Files) > 0 {
		buffer := new(bytes.Buffer)
		writer := multipart.NewWriter(buffer)

		for i, f := range wh.Files {
			fw, err := writer.CreateFormFile("file"+strconv.Itoa(i), f.Name)
			if err != nil {
				return err
			}

			if _, err := fw.Write(f.Bytes); err != nil {
				return err
			}
		}

		data, err := json.Marshal(wh)
		if err != nil {
			return err
		}

		if err := writer.WriteField("payload_json", string(data)); err != nil {
			return err
		}
		writer.Close()

		p.Bytes = buffer.Bytes()
		p.ContentType = writer.FormDataContentType()

	} else {
		data, err := json.Marshal(wh)
		if err != nil {
			return err
		}

		p.Bytes = data
		p.ContentType = "application/json"
	}

	if isErr {
		s.Queue = append([]MsgPayload{p}, s.Queue...)
	} else {
		s.Queue = append(s.Queue, p)
	}
	if s.Waiting {
		s.Waiting = false
		s.Waiter.Done()
	}
	return nil
}
