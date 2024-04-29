package discordwhsender

import "cmp"

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
