package main

import (
	"fmt"
	//"image"
	//"image/color"
	//"io/ioutil"
	"log"
	//"net/http"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tfacedetect-from-url [image URL] [image file]")
		return
	}

	// parse args
	imageURL := os.Args[1]
	saveFile := os.Args[2]


	cam, err := gocv.VideoCaptureFile(imageURL)
	if err != nil {
		log.Fatalf("Error opening video stream or file: %v", err)
	}
	defer cam.Close()

	img := gocv.NewMat()
	for ;; {
		if ok := cam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", err)
			return
		}
		if img.Empty() {
			continue
		}
		break
	}

	gocv.IMWrite(saveFile, img)
	fmt.Printf("saved to %s\n", saveFile)
}
/*

package main

import (
	"fmt"
	"net/http"
	"image"
	"os"
	"bufio"
	"strings"
	"bytes"
	"image/jpeg"
)

func main() {
	fmt.Println("début")

	test("http://192.168.1.49:81/stream")


/ *
	img, err := fetchImage("http://192.168.1.49:81/stream")
	if err != nil {
		fmt.Println("ici ça plante :", err)
		return
	}
	err = saveImageToFile(img, "image.jpg")
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde de l'image :", err)
		return
	}
	fmt.Println("fin")
	* /
}

func newFetchImage(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return 
	}
	defer resp.Body.Close()
	
	reader := bufio.NewReader(resp.Body)

	var img string = ""
	for {
		line, _ := reader.ReadString('\n')
		//fmt.Println(line)
		if strings.Contains(line, "Content-Type: image/jpeg") {
			//img = line
			line, _ = reader.ReadString('\n')
			fmt.Println("line", line)
			line, _ = reader.ReadString('\n')
			fmt.Println("line", line)
			for {
				line, _ = reader.ReadString('\n')
				fmt.Println("line", line)
				if strings.Contains(line, "--123456789000000000000987654321") {
					fmt.Println(len(img))
					break
				}
				img +=  line
			}
			break
		}	
	}
	fmt.Println("img", img)
	
	/ * Write the JPEG data directly to a file
	file, err := os.Create("frame.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(img)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Frame saved as frame.jpg")
	* /
	
	// Decode the JPEG image
	imgDecode, err := jpeg.Decode(strings.NewReader(img))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Save the image
	file, err := os.Create("frame.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = jpeg.Encode(file, imgDecode, nil)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Frame saved as frame.jpg")
}

func fetchImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ha")
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	return img, nil
}


func fetchStreamImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		fmt.Println(line)
		if strings.Contains(line, "Content-Type: image/jpeg") {
			img := make([]byte, 0)
			for {
				line, err := reader.ReadBytes('\n')


				fmt.Println("line", line)


				if err != nil {
					return nil, err
				}
				if strings.Contains(string(line), "--123456789000000000000987654321") && len(img) > 0 {
					fmt.Println(img)
					fmt.Println(len(img))
					break
				}
				if strings.Contains(string(line), "--123456789000000000000987654321") {
					fmt.Println("continue")
					continue
				}
				img = append(img, line...)
			}

			imgReader := bytes.NewReader(img)
			decodedImg, err := jpeg.Decode(imgReader)
			if err != nil {
				return nil, err
			}
			return decodedImg, nil
		}
		if err != nil {
			return nil, err
		}
	}
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







/*

package main

import (
	"bufio"
	"bytes"
	"fmt"
	//"image/jpeg"
	"net/http"
	"os"
	"regexp"
)

func main() {
	streamURL := "http://192.168.1.49:81/stream" // Replace with your ESP32-CAM IP

	// Connect to the MJPEG stream
	resp, err := http.Get(streamURL)
	if err != nil {
		fmt.Println("Error connecting to stream:", err)
		return
	}
	defer resp.Body.Close()

	// Create a buffered reader
	reader := bufio.NewReader(resp.Body)

	// Define MJPEG boundary regex
	boundaryRegex := regexp.MustCompile(`--(\S+)`)

	var imgBuffer bytes.Buffer
	isCapturing := false

	// Read the MJPEG stream
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading stream:", err)
			break
		}

		// Check for the MJPEG boundary
		if boundaryRegex.Match(line) {
			if imgBuffer.Len() > 0 {
				break // We have captured one frame
			}
			isCapturing = true
			imgBuffer.Reset()
			continue
		}

		// Capture JPEG data
		if isCapturing {
			imgBuffer.Write(line)
		}
	}

	fmt.Println( imgBuffer.String())
	fmt.Println( imgBuffer.Len())

	// Write the JPEG data directly to a file
	file, err := os.Create("frame.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(imgBuffer.Bytes())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Frame saved as frame.jpg")
	/*
	// Decode the JPEG image
	img, err := jpeg.Decode(bytes.NewReader(imgBuffer.Bytes()))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Save the image
	file, err := os.Create("frame.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Frame saved as frame.jpg")
	* /
}*/
