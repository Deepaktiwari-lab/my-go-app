package main

import (
	"html/template"
	"net/http"
	"sync"
)

var (
	indiaVotes    int
	pakistanVotes int
	voteMutex     sync.Mutex
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/results", resultsHandler)

	http.ListenAndServe(":8081", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Deepak Tiwari Election</title>
	</head>
	<body>
		<h1>Hello, welcome to Deepak Tiwari's election!</h1>
		<p>Which team do you like more?</p>
		<form action="/vote" method="post">
			<input type="radio" name="team" value="india" id="india">
			<label for="india">India</label><br>
			<input type="radio" name="team" value="pakistan" id="pakistan">
			<label for="pakistan">Pakistan</label><br><br>
			<input type="submit" value="Vote">
		</form>
		<p><a href="/results">See Results</a></p>
	</body>
	</html>`
	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		team := r.FormValue("team")

		voteMutex.Lock()
		if team == "india" {
			indiaVotes++
		} else if team == "pakistan" {
			pakistanVotes++
		}
		voteMutex.Unlock()
		http.Redirect(w, r, "/results", http.StatusSeeOther)
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Election Results</title>
	</head>
	<body>
		<h1>Election Results</h1>
		<p>India: {{.IndiaVotes}} votes</p>
		<p>Pakistan: {{.PakistanVotes}} votes</p>
		{{if gt .IndiaVotes .PakistanVotes}}
			<p>🎉 Hurray! India won! 😄</p>
		{{else if gt .PakistanVotes .IndiaVotes}}
			<p>Pakistan won this time. 🙁</p>
		{{else}}
			<p>The election is a tie! 🤔</p>
		{{end}}
		<p><a href="/">Go Back</a></p>
	</body>
	</html>`
	t, _ := template.New("results").Parse(tmpl)
	data := struct {
		IndiaVotes    int
		PakistanVotes int
	}{
		IndiaVotes:    indiaVotes,
		PakistanVotes: pakistanVotes,
	}
	t.Execute(w, data)
}

