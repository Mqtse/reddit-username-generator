package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)
type ResponseWrapper struct {
	Data Response `json:"data"`
}
type Response struct {
	GeneratedUsernames []string `json:"generatedUsernames"`
	Duration float32 `json:"duration"`
	Errors []string `json:"errors"`
}
func main() {
	httpClient := http.Client{}
	request, error := http.NewRequest(http.MethodPost, "https://old.reddit.com/svc/shreddit/graphql", nil)
	request.Header.Set("User-Agent", "Mozilla")
	request.Header.Set("Cookie", "csrf_token=c")
	request.Header.Set("Content-Type", "application/json")
	mux := http.NewServeMux()
	if error != nil {
		fmt.Println(error)
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})
	mux.HandleFunc("/username", func(w http.ResponseWriter, r *http.Request) {
		request.Body = io.NopCloser(bytes.NewBufferString(`{"operation":"GeneratedUsernames","variables":{"count":1},"csrf_token":"c"}`))
		response, error := httpClient.Do(request)
		if error != nil {
			fmt.Println(error)
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		var data ResponseWrapper
		json.Unmarshal([]byte(body), &data)
		io.WriteString(w, data.Data.GeneratedUsernames[0])
	})
	server := http.Server {
		Addr: ":5050",
		Handler: mux,
	}
	log := server.ListenAndServe()
	fmt.Print(log)
}