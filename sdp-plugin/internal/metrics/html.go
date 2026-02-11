package metrics

import (
	"strings"
)

// GenerateHTML creates HTML benchmark report (AC5).
func (r *Reporter) GenerateHTML() (string, error) {
	markdown, err := r.GenerateMarkdown()
	if err != nil {
		return "", err
	}

	// Simple markdown to HTML conversion
	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>AI Code Quality Benchmark</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
		table { border-collapse: collapse; width: 100%%; }
		th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
		th { background-color: #4CAF50; color: white; }
		h1 { color: #333; }
		h2 { color: #666; border-bottom: 1px solid #eee; padding-bottom: 10px; }
	</style>
</head>
<body>
` + markdownToHTML(markdown) + `
</body>
</html>`

	return html, nil
}

// markdownToHTML converts basic markdown to HTML (simplified).
func markdownToHTML(md string) string {
	// Very basic conversion - for production, use a proper library
	html := md
	// Headers
	replacer := strings.NewReplacer("# ", "<h1>", "## ", "<h2>", "### ", "<h3>")
	html = replacer.Replace(html)
	// Bold
	boldReplacer := strings.NewReplacer("**", "<strong>", "***", "</strong>")
	html = boldReplacer.Replace(html)
	// Line breaks
	brReplacer := strings.NewReplacer("\n\n", "</p><p>", "\n", "<br>")
	html = brReplacer.Replace(html)

	return html
}
