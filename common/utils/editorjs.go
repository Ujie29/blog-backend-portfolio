package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

type EditorJsBlock struct {
	Type string `json:"type"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

type EditorJsContent struct {
	Blocks []EditorJsBlock `json:"blocks"`
}

// regex: 移除 <b>、<i> 等 HTML tag
var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

func stripHTMLTags(input string) string {
	return htmlTagRegex.ReplaceAllString(input, "")
}

// 從 Editor.js JSON 擷取摘要
func ExtractSummaryFromEditorJS(jsonContent string, maxLength int) string {
	var content EditorJsContent
	if err := json.Unmarshal([]byte(jsonContent), &content); err != nil {
		return ""
	}

	var sb strings.Builder
	currentLen := 0

	for _, block := range content.Blocks {
		if block.Type == "paragraph" || block.Type == "header" {
			text := stripHTMLTags(block.Data.Text) // 移除 HTML 標籤

			for _, r := range text {
				sb.WriteRune(r)
				currentLen++
				if currentLen >= maxLength {
					return sb.String() + "..."
				}
			}
			sb.WriteString(" ")
		}
	}

	return strings.TrimSpace(sb.String())
}
