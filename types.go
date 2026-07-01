package discordwhsender

// Message is a Discord webhook payload. Construct it as a struct literal and
// call Send on it.
type Message struct {
	Content   string   `json:"content,omitempty"`
	Username  string   `json:"username,omitempty"`
	AvatarURL string   `json:"avatar_url,omitempty"`
	Embeds    []*Embed `json:"embeds,omitempty"`
	Files     []File   `json:"-"`
}

// Embed is a Discord rich embed, attached to a Message.
type Embed struct {
	URL         string    `json:"url,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Author      *Author   `json:"author,omitempty"`
	Color       int       `json:"color,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
	Thumbnail   *Image    `json:"thumbnail,omitempty"`
	Image       *Image    `json:"image,omitempty"`
	Fields      []*Fields `json:"fields,omitempty"`
	Footer      *Footer   `json:"footer,omitempty"`
}

// Author is an embed's author block.
type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// Image is a thumbnail or image attached to an embed.
type Image struct {
	URL string `json:"url,omitempty"`
}

// Fields is a single name/value field within an embed. Use Embed.AddField to
// append one with Discord's empty-string validation applied automatically.
type Fields struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// Footer is an embed's footer block.
type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// File is a file attachment sent alongside a Message.
type File struct {
	Name  string
	Bytes []byte
}
