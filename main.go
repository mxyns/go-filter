package main

import (
	"bufio"
	"flag"
	"fmt"
	filters "github.com/mxyns/go-filter/filter"
	filfuncs "github.com/mxyns/go-filter/filter/funcs"
	"github.com/mxyns/go-filter/io"
	"github.com/mxyns/go-tcp/fileio"
	log "github.com/sirupsen/logrus"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

// TODO gotcp ajouter parametre pour ignore le Await de wantsUserResponse dans SendRequestOn
// TODO parametres pour les filtres
// TODO filtres : gris, reduction bruit (moyenne pix alentours), bords (diff), code barre
func main() {

	// define general flags (some are defined in other packages' init functions, e.g: routines)
	runServer := flag.Bool("s", false, "run in server mode")
	address := flag.String("a", "127.0.0.1", "address to host on / connect to")
	proto := flag.String("P", "tcp", "protocol")
	port := flag.Uint("p", 8887, "port")
	timeout := flag.String("t", "10s", "client connection timeout")
	debugLevel := flag.String("l", "panic", "client connection timeout")
	customFormatter := flag.Bool("f", true, "use custom formatter for network log")

	flag.Parse()

	// apply custom log formatter for go-tcp logs
	if *customFormatter {
		log.SetFormatter(&io.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[go-tcp][%lvl%]: %time% - %msg% %fields%\n",
		})
	}

	// apply custom debug level for go-tcp logs
	level, err := log.ParseLevel(*debugLevel)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		log.SetLevel(level)
	}

	registerFilters()

	// free disk space on close
	defer fileio.ClearDir("./dl")
	defer fileio.ClearDir("./out")

	if *runServer {
		startServer(address, proto, port)
	} else {
		startClient(address, proto, port, timeout)
	}
}

func registerFilters() {

	filters.RegisterFilter(&filters.Filter{Name: "invert", Usage: "no args needed", Apply: filfuncs.InvertColor})
	filters.RegisterFilter(&filters.Filter{Name: "edges", Usage: "no args needed", Apply: filfuncs.FindEdges, Parser: filfuncs.ParseFindEdgesArgs})

	// useless filters
	filters.RegisterFilter(&filters.Filter{Name: "nullify", Usage: "no args needed", Apply: filfuncs.Nullify})
	filters.RegisterFilter(&filters.Filter{Name: "copy", Usage: "no args needed", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "identity", Usage: "no args needed", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "print", Usage: "no args needed", Apply: filfuncs.Print})
}

func terminalInput() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		line := scanner.Text()
		fmt.Printf("got > %v\n", line)
		if line == "stop" {
			break
		}
	}
}
