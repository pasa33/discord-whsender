package discordwhsender

type Message struct {
	Content   string   `json:"content,omitempty"`
	Username  string   `json:"username,omitempty"`
	AvatarURL string   `json:"avatar_url,omitempty"`
	Embeds    []*Embed `json:"embeds,omitempty"`
	Files     []File   `json:"-"`
}

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

type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
}

type Fields struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type File struct {
	Name  string
	Bytes []byte
}
