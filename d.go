package main
/*
import (
	"image/jpeg"
	"os"
	"gocv.io/x/gocv"
	"net/http"
	"strconv"
)

func test(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return 
	}
	defer resp.Body.Close()


	// Check if the URL is a webcam stream
	webcam, err := gocv.VideoCapture(url)
	if err != nil {
		return
	}
	defer webcam.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	frameCount := 0
	for {
		if ok := webcam.Read(&frame); !ok || frame.Empty() {
			break
		}

		// Convert frame to image.Image
		img, err := frame.ToImage()
		if err != nil {
			continue
		}

		// Save the image to a file
		fileName := "frame_" + strconv.Itoa(frameCount) + ".jpg"
		outFile, err := os.Create(fileName)
		if err != nil {
			continue
		}
		defer outFile.Close()

		err = jpeg.Encode(outFile, img, nil)
		if err != nil {
			continue
		}

		frameCount++
	}
}*/