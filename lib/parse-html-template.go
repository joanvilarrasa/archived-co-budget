package lib

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func ParseHtmlTemplate(path string, data any) string {
	htmlRawTemplate, htmlRawTemplateErr := os.ReadFile(path)
	if htmlRawTemplateErr != nil {
		fmt.Printf("Error reading : %v\n", htmlRawTemplateErr)
		return "¿¿ERROR??"
	}

	htmlTemplate, htmlTemplateErr := template.New("test").Parse(string(htmlRawTemplate))
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
