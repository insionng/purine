package model

import (
	"bytes"
	"time"
)

type Rss struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Items       []*Rss
}

func (rss *Rss) toStringItem() string {
	var buf bytes.Buffer
	buf.WriteString("<item>")

	buf.WriteString("<title>")
	buf.WriteString(rss.Title)
	buf.WriteString("</title>")

	buf.WriteString("<link>")
	buf.WriteString(rss.Link)
	buf.WriteString("</link>")

	buf.WriteString("<description><![CDATA[")
	buf.WriteString(rss.Description)
	buf.WriteString("]]></description>")

	buf.WriteString("<pubDate>")
	buf.WriteString(rss.PubDate.Format(time.RFC822))
	buf.WriteString("</pubDate>")

	buf.WriteString("</item>")
	return buf.String()
}

func (rss *Rss) ToString() string {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>`)

	buf.WriteString("<title>")
	buf.WriteString(rss.Title)
	buf.WriteString("</title>")

	buf.WriteString("<link>")
	buf.WriteString(rss.Link)
	buf.WriteString("</link>")

	buf.WriteString("<description><![CDATA[")
	buf.WriteString(rss.Description)
	buf.WriteString("]]></description>")

	buf.WriteString("<pubDate>")
	buf.WriteString(rss.PubDate.Format(time.RFC822))
	buf.WriteString("</pubDate>")

	for _, r := range rss.Items {
		buf.WriteString(r.toStringItem())
	}

	buf.WriteString(`</channel></rss>`)
	return buf.String()
}
