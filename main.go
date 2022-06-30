package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

type statusCodeHandler struct{}

func (sch statusCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log
	log.Println("Incoming request on:", r.Method, r.Host, r.URL, "From:", r.RemoteAddr, "Using:", r.UserAgent())

	// get the statuscode from query
	rawStatusCode := r.URL.Query().Get("stat")
	if rawStatusCode == "" {
		log.Println("Returning usage")
		w.Write([]byte("usage: `/?stat=CODE`; e.g.: `/?stat=500`"))
		return
	}

	// check if that statuscode exist and send it back as empty with that statuscode at header
	statusCode, err := CheckQuery(rawStatusCode)
	if err != nil {
		log.Println("Returning error as query was not expected")
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("Returning statuscode", statusCode, http.StatusText(statusCode))
	w.WriteHeader(statusCode)
	return
}

func CheckQuery(query string) (int, error) {
	// if query is not only number
	stat, err := strconv.Atoi(query)
	if err != nil {
		return 0, err
	}

	string := http.StatusText(stat)
	if string == "" {
		return 0, errors.New("code does not exist")
	}

	return stat, nil
}

func main() {
	mux := http.NewServeMux()

	sch := statusCodeHandler{}

	mux.Handle("/", sch)

	log.Print("Listening...")
	http.ListenAndServe(":19011", mux)
}
