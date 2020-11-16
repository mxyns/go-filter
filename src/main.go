package main

import (
	filters "filter"
	filfuncs "filter/funcs"
	"flag"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"pouf"
)

// TODO gotcp ajouter parametre pour ignore le Await de wantsUserResponse dans SendRequestOn

// TODO : image mutable
// TODO filtres : gris, reduction bruit (moyenne pix alentours), bords (diff), code barre
func main() {

	runServer := flag.Bool("s", false, "run in server mode")
	address := flag.String("a", "127.0.0.1", "address to host on / connect to")
	proto := flag.String("P", "tcp", "protocol")
	port := flag.Uint("p", 8887, "port")
	timeout := flag.String("t", "10s", "client connection timeout")

	flag.Parse()

	registerFilters()

	// TODO image => filtres => response until closed
	// TODO image + filtre(s) => response => close
	if *runServer {
		pouf.MainServer(address, proto, port)
	} else {
		pouf.MainClient(address, proto, port, timeout)
	}
}

func registerFilters() {

	filters.RegisterFilter(&filters.Filter{Name: "invert", Apply: filfuncs.InvertColor})
	filters.RegisterFilter(&filters.Filter{Name: "nullify", Apply: filfuncs.Nullify})
	filters.RegisterFilter(&filters.Filter{Name: "copy", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "identity", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "print", Apply: filfuncs.Print})
}
