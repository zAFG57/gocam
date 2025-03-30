package main

import (
	"fmt"
	"net/http"
	"image"
	"image/jpeg"
	"os"
)

func main() {
	fmt.Println("début")
	img, err := fetchStreamImage("http://http://192.168.1.49:81/stream")
	if err != nil {
		fmt.Println("Erreur lors de la récupération de l'image :", err)
		return
	}
	err = saveImageToFile(img, "image.jpg")
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde de l'image :", err)
		return
	}
	fmt.Println("fin")
}

func fetchStreamImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func saveImageToFile(img image.Image, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		return err
	}

	return nil
}