package main

import "fserver"

func main() {
	// Define routes as a string (supports text, HTML, and files)
	routes := `
	/: file:index.html
	/hello: Hello, world!
	/about: file:about.html
	/data: file:data.json
	`

	// Start the server with a custom 404 page (text or file)
	fserver.Server("8080", routes, "custom404.html")
}
