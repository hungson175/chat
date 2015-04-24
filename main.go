package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/hungson175/chat/trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		fullPath := filepath.Join("templates", t.filename)
		t.templ = template.Must(template.ParseFiles(fullPath))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//For assets Bootstrap
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("path/to/assets")))) //change the path later on

	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting webserver on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen and server: ", err)
	}
}
