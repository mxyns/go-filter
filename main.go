package main

import (
	"bufio"
	"flag"
	"fmt"
	filters "github.com/mxyns/go-filter/filter"
	filfuncs "github.com/mxyns/go-filter/filter/funcs"
	"github.com/mxyns/go-filter/io"
	"github.com/mxyns/go-tcp/fileio"
	"github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	log "github.com/sirupsen/logrus"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func init() {

	fileio.MAX_FILE_BUFFER_SIZE *= 2

	log.SetLevel(log.PanicLevel)
}

// todo gotcp ajouter parametre pour ignore le Await de wantsUserResponse dans SendRequestOn
// todo filtres : possibilité d'élargir en x ou y
func main() {

	// define general flags (some are defined in other packages' init functions, e.g: routines)
	runServer := flag.Bool("s", false, "run in server mode")
	address := flag.String("a", "127.0.0.1", "address to host on / connect to")
	proto := flag.String("P", "tcp", "protocol")
	port := flag.Uint("p", 8887, "port")
	timeout := flag.String("t", "10s", "client connection timeout")
	debugLevel := flag.String("l", "panic", "debug level")
	deleteFiles := flag.Bool("d", true, fmt.Sprintf("clear directory %v and %v on close", "./"+defaultRequests.TARGET_DIRECTORY, io.OutDir))
	goTCPOutDir := flag.String("o", io.OutDir, fmt.Sprintf("change %v temporary output directory. be sure it exists. won't be created automatically", "./"+defaultRequests.TARGET_DIRECTORY))
	customFormatter := flag.Bool("f", true, "use custom formatter for network log")

	flag.Parse()

	defaultRequests.TARGET_DIRECTORY = *goTCPOutDir

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
	if *deleteFiles {
		defer fileio.ClearDir("./" + defaultRequests.TARGET_DIRECTORY)
		defer fileio.ClearDir(io.OutDir) // only for server
	}

	if *runServer {
		startServer(address, proto, port)
	} else {
		startClient(address, proto, port, timeout)
	}
}

func registerFilters() {

	filters.RegisterFilter(&filters.Filter{Name: "invert", Usage: "no args needed", Apply: filfuncs.InvertColor})
	filters.RegisterFilter(&filters.Filter{Name: "grayScaleAverage", Usage: "no args needed", Apply: filfuncs.GrayScaleAverage})
	filters.RegisterFilter(&filters.Filter{Name: "grayScaleLuminosity", Usage: "no args needed", Apply: filfuncs.GrayScaleLuminosity})
	filters.RegisterFilter(&filters.Filter{Name: "grayScaleDesaturation", Usage: "no args needed", Apply: filfuncs.GrayScaleDesaturation})

	filters.RegisterFilter(&filters.Filter{Name: "edges", Usage: "radius=(int) threshold=(int) dist=[euclidean|norm1|inf] edge_color=[original|(r,g,b,a)|[r,g,b,a]] gap_color=(see edge_color)", Apply: filfuncs.FindEdges, Parser: filfuncs.ParseFindEdgesArgs})
	filters.RegisterFilter(&filters.Filter{Name: "noiseReduction", Usage: "radius(int)", Apply: filfuncs.NoiseReduction, Parser: filfuncs.ParseNoiseReductionArgs})

	// useless filters
	filters.RegisterFilter(&filters.Filter{Name: "nullify", Usage: "no args needed", Apply: filfuncs.Nullify})
	filters.RegisterFilter(&filters.Filter{Name: "copy", Usage: "no args needed", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "identity", Usage: "no args needed", Apply: filfuncs.Identity})
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
