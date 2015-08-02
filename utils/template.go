package utils

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"strings"
	"time"
)

// string to html
func Str2HTML(str string) template.HTML {
	return template.HTML(str)
}

// format time unixstamp
func TimeUnixFormat(unixStamp int64, layout string) string {
	return time.Unix(unixStamp, 0).Format(layout)
}

// format time unixstamp friendly,
// like xxx seconds ago
func FriendTimeUnixFormat(unixStamp int64) string {
	t := time.Unix(unixStamp, 0)
	seconds := int64(time.Since(t).Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%d Seconds Ago", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%d Minutes Ago", seconds/60)
	}
	if seconds < 86400 {
		return fmt.Sprintf("%d Hours Ago", seconds/3600)
	}
	return fmt.Sprintf("%d Days Ago", seconds/86400)
}

// markdown to html
func Md2Html(str string) template.HTML {
	return template.HTML(Md2String(str))
}

// markdown to string
func Md2String(str string) string {
	return string(blackfriday.MarkdownCommon([]byte(str)))
}

// format bytes size friendly,
// like xx.x MB
func FriendBytesSize(size int64) string {
	sFloat := float64(size)
	if sFloat >= 1024*1024 {
		return fmt.Sprintf("%.1f MB", sFloat/1024/1024)
	}
	if sFloat > 1024 {
		return fmt.Sprintf("%.1f KB", sFloat/1024)
	}
	return fmt.Sprintf("%.1f B", sFloat)
}

// newline 2 break
func Nl2Br(str string) template.HTML {
	return template.HTML(strings.Replace(str, "\n", "<br/>", -1))
}
