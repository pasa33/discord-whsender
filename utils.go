package discordwhsender

import "cmp"

// AddField appends a field to the embed and returns it for chaining.
// Discord rejects empty field name/value, so blanks are replaced with "-".
func (e *Embed) AddField(name, value string, inline bool) *Embed {
	e.Fields = append(e.Fields, &Fields{
		Name:   cmp.Or(name, "-"),
		Value:  cmp.Or(value, "-"),
		Inline: inline,
	})
	return e
}
