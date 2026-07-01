package discordwhsender

import "cmp"

// validate enforces Discord's embed field constraints (non-empty, length
// limits) before a message is sent.
func (msg *Message) validate() {
	for _, emb := range msg.Embeds {
		for _, v := range emb.Fields {
			v.Name = cmp.Or(truncate(v.Name, 256), "-")
			v.Value = cmp.Or(truncate(v.Value, 1024), "-")
		}
	}
}

func truncate(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:max-3] + "..."
}
