package main

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
)

// get comment from JSONPlaceholder
type placeholder interface {
	getComment() (string, error)
}

type JSONPlaceHolder struct{}

func main() {
	port := ":9008"
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", ping)
	mux.HandleFunc("/comment", commentHandler)
	fmt.Println("server is running on port" + fmt.Sprintf("%v", port))
	err := http.ListenAndServe(fmt.Sprintf("%v", port), mux)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (h *JSONPlaceHolder) getComment() (string, error) {
	postId := rand.IntN(10)
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", postId))
	if err != nil {
		slog.Error("could not complete request", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response body", "error", err)
		return "", err
	}

	return string(body), nil
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	h := &JSONPlaceHolder{}
	comment, err := h.getComment()
	if err != nil {
		http.Error(w, "Failed to get comment", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, comment)
}
