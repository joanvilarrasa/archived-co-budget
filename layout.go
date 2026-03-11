package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func navHeader() string {
	return "CO-Budget"
}

type HomeProps struct {
	NavHeader string
}

func home() string {
	homeprops := HomeProps{NavHeader: navHeader()}

	htmlRawTemplate, htmlRawTemplateErr := os.ReadFile("./layout.html")
	if htmlRawTemplateErr != nil {
		fmt.Printf("Error reading app/layout.html: %v\n", htmlRawTemplateErr)
		return "failed to read app/layout.html"
	}

	htmlTemplate, htmlTemplateErr := template.New("test").Parse(string(htmlRawTemplate))
	if htmlTemplateErr != nil {
		fmt.Printf("Error parsing template: %v\n", htmlTemplateErr)
		return "failed to parse the template"
	}

	var sb strings.Builder
	err := htmlTemplate.Execute(&sb, homeprops)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return fmt.Sprintf("failed to execute template: %v", err)
	}

	return sb.String()
}
