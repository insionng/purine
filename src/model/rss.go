package model

import (
	"bytes"
	"github.com/fuxiaohei/purine/src/utils"
	"path"
	"time"
)

// rss struct
type Rss struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Items       []*Rss
}

// convert rss item to string
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

// rss to xml string
func (rss *Rss) String() string {
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

// convert articles to rss struct
func Articles2Rss(articles []*Article, link string) *Rss {
	rss := &Rss{
		PubDate: time.Now(),
		Items:   make([]*Rss, len(articles)),
		Link:    link,
	}
	for i, a := range articles {
		r := &Rss{
			Title:       a.Title,
			Link:        link + path.Join("article", a.Href()),
			Description: utils.Md2String(a.Body),
			PubDate:     time.Unix(a.CreateTime, 0),
			Items:       nil,
		}
		if i == 0 {
			rss.PubDate = r.PubDate
		}
		rss.Items[i] = r
	}
	return rss
}

// sitemap url struct
type SiteMapUrl struct {
	Loc             string
	LastMod         time.Time
	ChangeFrequency string
	Priority        string
}

// convert sitemap url to string
func (s *SiteMapUrl) String() string {
	var buf bytes.Buffer
	buf.WriteString("<loc>")
	buf.WriteString(s.Loc)
	buf.WriteString("</loc>")

	buf.WriteString("<lastmod>")
	buf.WriteString(s.LastMod.Format(time.RFC3339))
	buf.WriteString("</lastmod>")

	buf.WriteString("<changefreq>")
	buf.WriteString(s.ChangeFrequency)
	buf.WriteString("</changefreq>")

	buf.WriteString("<priority>")
	buf.WriteString(s.Priority)
	buf.WriteString("</priority>")

	return buf.String()
}

// sitemap group struct,
// contains main url and site urls
type SiteMapGroup struct {
	Url     []*SiteMapUrl
	Loc     string
	LastMod time.Time
}

// convert sitemap group to string
func (s *SiteMapGroup) String() string {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?><?xml-stylesheet type="text/xsl" href="/static/sitemap.xsl"?>`)
	buf.WriteString(`<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	buf.WriteString("<url><loc>")
	buf.WriteString(s.Loc)
	buf.WriteString("</loc>")

	buf.WriteString("<lastmod>")
	buf.WriteString(s.LastMod.Format(time.RFC3339))
	buf.WriteString("</lastmod>")

	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>1.0</priority></url>")

	for _, s := range s.Url {
		buf.WriteString("<url>")
		buf.WriteString(s.String())
		buf.WriteString("</url>")
	}

	buf.WriteString(`</urlset>`)
	return buf.String()
}

// convert articles to sitemap group
func Articles2SiteMap(articles []*Article, link string) *SiteMapGroup {
	group := &SiteMapGroup{
		Url: make([]*SiteMapUrl, len(articles)),
		Loc: link,
	}
	for i, a := range articles {
		r := &SiteMapUrl{
			Loc:             link + path.Join("article", a.Href()),
			LastMod:         time.Unix(a.CreateTime, 0),
			ChangeFrequency: "weekly",
			Priority:        "0.6",
		}
		if i == 0 {
			group.LastMod = r.LastMod
		}
		group.Url[i] = r
	}
	return group
}
