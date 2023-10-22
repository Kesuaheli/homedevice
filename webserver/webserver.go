package webserver

import (
	"html/template"
	"io"
	logger "log"
	"net/http"

	"github.com/gorilla/mux"
)

var log logger.Logger = *logger.New(logger.Writer(), "[WEB] ", logger.LstdFlags|logger.Lmsgprefix)

func Run() {
	addr := ":8080"
	go run(setupHandler(), addr)
	log.Printf("Server running on '%s'!", addr)
}

func setupHandler() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handle404)

	fs := http.FileServer(http.Dir("files/"))
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", fs))
	r.Path("/files").Handler(http.RedirectHandler("/files/", http.StatusMovedPermanently))

	r.Path("/favicon.ico").Handler(http.RedirectHandler("/files/favicon.ico", http.StatusMovedPermanently))

	r.HandleFunc("/api", handle404)

	r.Use(mwfLogRequest)
	return r
}

func run(handler http.Handler, addr string) {
	log.Fatal(http.ListenAndServe(addr, handler))
}

func mwfLogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest("handle ", r)
		handler.ServeHTTP(w, r)
	})
}
func logRequest(prefix string, r *http.Request) {
	log.Printf("%s: %s%s %s", r.RemoteAddr, prefix, r.Method, r.URL)
}

func logReqestData(r *http.Request) {
	if r.Body == nil {
		return
	}
	buf := make([]byte, 1000)
	n, err := r.Body.Read(buf)
	if err != io.EOF {
		log.Printf("ERROR: could not read body: %+v\n\tbuf: %+v\n\tstrbuf: %s", err, buf, string(buf))
		return
	}
	if n == 0 {
		return
	}
	log.Print(string(buf[:n]))
}

func handle404(w http.ResponseWriter, r *http.Request) {
	logRequest("404 -> ", r)
	logReqestData(r)
	data := struct {
		Path string
	}{
		Path: r.RequestURI,
	}
	w.WriteHeader(http.StatusNotFound)
	if err := template.Must(template.ParseFiles("template/web/404_not_found.html")).Execute(w, data); err != nil {
		log.Printf("ERROR could not parse 404 template: %+v", err)
	}
}
