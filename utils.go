package discordwhsender

import "cmp"

func (e *Embed) AddField(name, value string, inline bool) {
	e.Fields = append(e.Fields, &Fields{
		Name:   cmp.Or(name, "-"),
		Value:  cmp.Or(value, "-"),
		Inline: inline,
	})
}
