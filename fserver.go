package fserver

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// parseRoutes converts a raw string into a map of paths and messages.
func parseRoutes(data string) map[string]string {
	routes := make(map[string]string)
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			path := strings.TrimSpace(parts[0])
			message := strings.TrimSpace(parts[1])
			routes[path] = message
		}
	}

	return routes
}

// loadFileContent reads a file if the message starts with "file:"
func loadFileContent(message string) (string, string, bool) {
	if strings.HasPrefix(message, "file:") {
		filename := strings.TrimPrefix(message, "file:")
		content, err := ioutil.ReadFile(strings.TrimSpace(filename))
		if err != nil {
			return "Error loading file: " + err.Error(), "text/plain", false
		}

		// Detect file type
		contentType := "text/plain"
		if strings.HasSuffix(filename, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(filename, ".json") {
			contentType = "application/json"
		}

		return string(content), contentType, strings.HasSuffix(filename, ".json")
	}
	return message, "text/plain", false
}

// extractQueryParams replaces {{.param:default}} with actual values
func extractQueryParams(content string, r *http.Request) string {
	re := regexp.MustCompile(`\{\{\.([a-zA-Z0-9_]+)(?::([^}]+))?\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		paramName := match[1]    // Variable name (e.g., "name")
		defaultValue := match[2] // Default value (e.g., "Guest")
		value := r.URL.Query().Get(paramName)

		if value == "" {
			value = defaultValue
		}
		content = strings.ReplaceAll(content, match[0], value)
	}

	return content
}

// Server starts an HTTP server with templating support
func Server(port string, rawRoutes string, notFoundPage string) {
	routes := parseRoutes(rawRoutes)
	custom404Message, _, _ := loadFileContent(notFoundPage)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if message, exists := routes[r.URL.Path]; exists {
			content, contentType, _ := loadFileContent(message)

			// Replace placeholders {{.param:default}} before parsing
			content = extractQueryParams(content, r)

			// If it's an HTML file, use Go templates
			if contentType == "text/html" {
				tmpl, err := template.New("page").Parse(content)
				if err != nil {
					http.Error(w, "Error processing template: "+err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text/html")
				tmpl.Execute(w, nil) // No need to pass params, they are already replaced
			} else {
				w.Header().Set("Content-Type", contentType)
				fmt.Fprint(w, content)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, custom404Message)
		}
	})

	fmt.Println("Server is running on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
