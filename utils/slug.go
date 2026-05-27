package utils

import (
	"fmt"
	"strings"
	"time"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	return fmt.Sprintf("%s-%d", slug, time.Now().Unix())
}
