package imageio

import (
	"fmt"
	im "image"
	"image/png"
	"os"
	"strings"
)

var in_dir = "./in/"
var out_dir = "./out/"

func init() {

	if _, err := os.Stat(in_dir); os.IsNotExist(err) {
		_ = os.Mkdir(in_dir, os.ModeDir)
	}

	if _, err := os.Stat(out_dir); os.IsNotExist(err) {
		_ = os.Mkdir(out_dir, os.ModeDir)
	}
}

func GetImageNames(args []string) []string {

	start := -1
	for i := range args {
		if args[i] == "-i" || args[i] == "--input" && i+1 < len(args) {
			start = i + 1
		}
	}
	if start > -1 {
		return args[start:]
	} else {
		panic("Pas d'images données en entrée")
	}
}
func LoadImage(filename string) (*im.Image, string) {

	file, err := os.Open(in_dir + filename)
	if err != nil {
		fmt.Printf("Erreur chargement fichier : %v", err)
		os.Exit(1)
	}
	defer file.Close()

	image, format, err := im.Decode(file)
	if err != nil {
		fmt.Printf("Erreur decodage image : %v", err)
		os.Exit(1)
	}

	return &image, format
}
func SaveImage(image *im.RGBA, name string, filter_name string) {

	shards := strings.Split(name, ".")
	new_name := ""
	for i := range shards[:len(shards)-2] {
		new_name += shards[i] + "."
	}
	new_name += shards[len(shards)-2] + "-" + filter_name + "." + shards[len(shards)-1]

	out, err := os.Create(out_dir + new_name)
	if err != nil {
		fmt.Printf("Erreur creation fichier : %v", err)
		os.Exit(1)
	}
	defer out.Close()

	err = png.Encode(out, image)
	if err != nil {
		fmt.Printf("Erreur ecriture image : %v", err)
		os.Exit(1)
	}
}
