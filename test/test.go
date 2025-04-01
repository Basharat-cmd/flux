package main

import (
	"fmt"
	"fserver"
	"strings"
	"time"
)

func main() {

	routes := `
	/: file:index.html
	/hello: Hello, world!
	/newpage: {{.data}}
	/data: file:data.json
	`
	go fserver.Server("8080", routes, "file:custom404.html")

	var old_url string = ""

	for {
		time.Sleep(2 * time.Second)
		requestedURL := fserver.GetRequestedURL()
		if requestedURL != "" && requestedURL != old_url {
			if requestedURL == "/:submit=true" {
				fmt.Print("sus")
			} else if strings.Contains(requestedURL, "newpage") {
				so, _ := fserver.FetchWebPageContent("http://localhost:8080" + requestedURL)
				fmt.Print(so)
			}
			old_url = requestedURL
		} else {
			fmt.Println(requestedURL)
		}
	}
}
