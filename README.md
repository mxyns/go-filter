# go-filter

## What is it?
The project allows clients to send files to a server which process them (by applying filters) and then send them back to the clients.

## Filter List
* `nullify`
* `copy` or `identity`
* `invert`
* grayscale filters :
    * `grayScaleAverage`
    * `grayScaleLuminosity`
    * `grayScaleDesaturation`
    * (to learn more about grayScaleFilters click [here](https://tannerhelland.com/2011/10/01/grayscale-image-algorithm-vb6.html)
* `noiseReduction` (a.k.a blur), parameters :
    - `radius=`int (default=1) : larger is blurrier (e.g.: `noiseReduction radius`)
* `edges` return an image with different color for edges and gap :
    - `radius=`int (default=1) : smaller gives sharper edges (e.g.: `radius=2`)
    - `threshold=`float64 (default=1000) : threshold for the detection of edges (e.g.: `threshold=500` )
    - `dist=`string (default=euclidean) : function used to compute color distance in RGBA color space
        available functions :
        - `dist=euclidean`
        - `dist=norm1`
    - `edge_color=`(R,G,B,A) [warning : no spaces] (default=BLACK=[0,0,0,255])
    - `gap color=`(R,G,B,A) [warning : no spaces] (default=TRANSPARENT_WHITE=[255,255,255,0])
    when the color is not well defined, it will keep the originals colors

to use a filter with arguments, simply give `filter arg1=value arg2=value ...`. Every argument has a default value that will be used if custom value is not provided by the user
e.g.: `edges radius=5 threshold=1000 dist=euclidean edge_color=(100,100,100,75) gap_color=(10,10,10,100)` will run an edge detection filter with a radius of 5, a threshold of 1000, the norm euclidean, and a chosen edge color & gap color

## Building project
`git clone https://github.com/mxyns/go-filter`
`cd go-filter`
`go build`

## Launching
If you want to launch the server mode you have to use the `-s` (or `-s=true`) flag, otherwise the client mode is launched by default.

### Server Mode
It should be noted :

Common flags :
* `-a`  (string, default = 127.0.0.1) : the address of the server
* `-P`  (string, default = tcp) : the protocol you want to use
* `-p`  (int, default = 8887) : the port of the server
* `-t`  (string, default = 10s) : the time after client connection crash
* `-l`  (string, default = "panic") : the level of debug, [logrus](https://godoc.org/github.com/sirupsen/logrus) has seven log levels: Trace, Debug, Info, Warn, Error, Fatal, and Panic
* `-f`  (boolean, default = true) : the custom formatter for the message in the terminal
* `-d`  (boolean, default = true) : clear go-tcp files target directory and go-filter's outDir(used by server) on close
* `-o`  (string, default = "panic") : change the temporary output directory. be sure it exists. won't be created automatically


Server specific flags :
* `-r`      (int, default = 1) : number of image processing routines
* `-scvert` (int, default = 5) : vertical slice count per image
* `-schor`  (int, default = 5) : horizontal slice count per image

### Client Mode
You can use all of the server's common flags, as well as :
* `-i`  (string, default = "") : the path of the image you want to use
* `-fl` (string, default = "copy") [requires `-i`] : the list of filters you want to apply on your image
Using `-i` will put the client in "direct mode". It will immediately send a request to the server using the values provided and close after reception of the result.
Omitting `-i` will put the client in "interactive mode", sequentially asking for the input file path (e.g.: `./path/to/file.png`), the filter list you want to apply (e.g.: `filter1 filter2 ...`)

## Dependencies
* [go-tcp](https://github.com/mxyns/go-tcp)
    * which requires [logrus](https://github.com/sirupsen/logrus)