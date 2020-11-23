package io

import (
	"fmt"
	im "image"
	"image/png"
	"os"
	"strings"
)

const (
	InDir  = "./in/"
	OutDir = "./out/"
)

func init() {

	if _, err := os.Stat(InDir); os.IsNotExist(err) {
		_ = os.Mkdir(InDir, os.ModeDir)
	}

	if _, err := os.Stat(OutDir); os.IsNotExist(err) {
		_ = os.Mkdir(OutDir, os.ModeDir)
	}
}

func LoadImage(filename *string) (*im.Image, string) {

	file, err := os.Open(*filename)
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
func SaveImage(image *im.RGBA, name *string, filter_name *string) *string {

	shards := strings.Split(*name, "/")
	filename := shards[len(shards)-1]

	shards = strings.Split(filename, ".")
	new_name := ""
	if len(shards) > 1 {
		for i := range shards[:len(shards)-2] {
			new_name += shards[i] + "."
		}
		new_name += shards[len(shards)-2] + "-" + *filter_name + "." + shards[len(shards)-1]
	} else {
		new_name = shards[0] + "-" + *filter_name + ".png"
	}

	outPath := OutDir + new_name
	out, err := os.Create(outPath)
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

	return &outPath
}
func RemoveFile(filename string) error {

	return os.Remove(filename)
}
