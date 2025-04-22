package main

import (
	"fmt"
	"log"
	"os"
	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Pour lancer le programe:\n\t./webCam [stream URL]")
		return
	}

	imageURL := os.Args[1]

	cam, err := gocv.VideoCaptureFile(imageURL)
	window := gocv.NewWindow("camera")
	defer window.Close()
	if err != nil {
		log.Fatalf("Error opening video stream or file: %v", err)
	}
	defer cam.Close()

	img := gocv.NewMat()
	lastFrames := [3]gocv.Mat{
		gocv.NewMat(),
		gocv.NewMat(),
		gocv.NewMat(),
	}
	defer func() {
		for _, frame := range lastFrames {
			frame.Close()
		}
	}()
	
	
	for iteration := uint8(0);; iteration = (iteration+1)%3{
		if ok := cam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", err)
			return
		}
		if img.Empty() {
			continue
		}
		clearFrames(img, &lastFrames, iteration)
		window.IMShow(img)
		window.WaitKey(1)
	}
}

func clearFrames(frame gocv.Mat, lastFrames *[3]gocv.Mat, iteration uint8) {
	if lastFrames[0].Empty() || lastFrames[1].Empty() || lastFrames[2].Empty() {
		shiftFrames(frame, lastFrames, iteration)
		return
	}
	avgDiff := gocv.NewMat()
	diff1 := gocv.NewMat()
	diff2 := gocv.NewMat()
	diff3 := gocv.NewMat()
	defer avgDiff.Close()
	defer diff1.Close()
	defer diff2.Close()
	defer diff3.Close()

	gocv.AbsDiff((*lastFrames)[0], (*lastFrames)[1], &diff1)
	gocv.AbsDiff((*lastFrames)[1], (*lastFrames)[2], &diff2)
	gocv.AbsDiff((*lastFrames)[2], (*lastFrames)[0], &diff3)

	// Add the differences together
	gocv.Add(diff1, diff2, &avgDiff)
	gocv.Add(avgDiff, diff3, &avgDiff)

	// Divide by 3 to get the average
	scalarMat := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(3, 3, 3, 0), avgDiff.Rows(), avgDiff.Cols(), avgDiff.Type())
	defer scalarMat.Close()
	gocv.Divide(avgDiff, scalarMat, &avgDiff)

	sumDiff := avgDiff.Sum()
	rsum := sumDiff.Val1 + sumDiff.Val2 + sumDiff.Val3
	if rsum < 255000*3 {
		removeImpurty(frame,avgDiff, lastFrames)
		shiftFrames(frame, lastFrames, iteration)
		fmt.Printf("c'est ok: %v\n", rsum)
	} else {
		frame.CopyTo(&(*lastFrames)[0])
		frame.CopyTo(&(*lastFrames)[1])
		frame.CopyTo(&(*lastFrames)[2])
		fmt.Printf("ça bouge trop: %v\n", rsum)
	}
}

func removeImpurty(frame gocv.Mat, avgDiff gocv.Mat, lastFrames *[3]gocv.Mat) {
	// Create a mask from the average difference
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.Threshold(avgDiff, &mask, 7, 255, gocv.ThresholdBinary)

	// Set all pixels in the frame to 255 where the mask is non-zero
	whiteMat := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), frame.Type())
	defer whiteMat.Close()
	whiteMat.SetTo(gocv.NewScalar(255, 255, 255, 0))
	
	// replace corrupted pixels with white
	for y := 0; y < frame.Rows(); y++ {
		for x := 0; x < frame.Cols(); x++ {
			if mask.GetUCharAt(y, x) != 0 {
				for i := 0; i < 3; i++ {
					c1 := (*lastFrames)[0].GetUCharAt(y, x*frame.Channels()+i)
					c2 := (*lastFrames)[1].GetUCharAt(y, x*frame.Channels()+i)
					c3 := (*lastFrames)[2].GetUCharAt(y, x*frame.Channels()+i)
					frame.SetUCharAt(y, x*frame.Channels()+i, uint8((int(c1)+int(c2)+int(c3))/3))
				}
			}
		}
	}
}

func shiftFrames(frame gocv.Mat, lastFrames *[3]gocv.Mat, iteration uint8) {
	frame.CopyTo(&(*lastFrames)[iteration])
}