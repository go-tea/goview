package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/go-tea/tea"
	"github.com/go-tea/tea/serve"
)

var (
	serveAddr = flag.String("addr", "localhost", "Address to bind server to")
	servePort = flag.Int("port", 3344, "Port to serve from")
)

func main() {

	flag.Parse()

	addr := net.JoinHostPort(*serveAddr, strconv.Itoa(*servePort))

	r := tea.New(serve.RequestID, serve.Recoverer, serve.Logger)

	r.Get("/", home)
	r.Get("/nocontent", nocontent)
	r.Get("/text", text)
	r.Get("/json", json)
	r.Get("/jsonp", jsonp)
	r.Get("/yaml", yaml)
	r.Get("/xml", lmx)
	r.Get("/binary", binary)
	r.Get("/fileinline", fileinline)
	r.Get("/filedownload", filedownload)

	err := r.ListenAndServe(addr)
	if err != nil {
		fmt.Println("err:: %s", err)
	}

	os.Exit(0)

}
