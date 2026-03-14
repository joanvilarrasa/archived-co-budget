package lib

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"
)

var (
	htmlTemplateCache   = make(map[string]string)
	htmlTemplateCacheMu sync.RWMutex
)

func ParseHtmlTemplate(path string, data any, cacheString bool) string {
	htmlTemplateCacheMu.RLock()
	cachedTemplate, exists := htmlTemplateCache[path]
	htmlTemplateCacheMu.RUnlock()

	htmlRawTemplate := cachedTemplate
	if !exists {
		htmlRawTemplateBytes, htmlRawTemplateErr := os.ReadFile(path)
		if htmlRawTemplateErr != nil {
			fmt.Printf("Error reading : %v\n", htmlRawTemplateErr)
			return "¿¿ERROR??"
		}

		htmlRawTemplate = string(htmlRawTemplateBytes)
		htmlTemplateCacheMu.Lock()
		htmlTemplateCache[path] = htmlRawTemplate
		htmlTemplateCacheMu.Unlock()
	}

	htmlTemplate, htmlTemplateErr := template.New("test").Parse(htmlRawTemplate)
	if htmlTemplateErr != nil {
		fmt.Printf("Error parsing template: %v\n", htmlTemplateErr)
		return "¿¿ERROR??"
	}

	var sb strings.Builder
	err := htmlTemplate.Execute(&sb, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return "¿¿ERROR??"
	}

	return sb.String()
}
