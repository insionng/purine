package utils

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"time"
)

func Str2HTML(str string) template.HTML {
	return template.HTML(str)
}

func TimeUnixFormat(unixStamp int64, layout string) string {
	return time.Unix(unixStamp, 0).Format(layout)
}

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

func Md2Html(str string) template.HTML {
	return template.HTML(string(blackfriday.MarkdownCommon([]byte(str))))
}

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
