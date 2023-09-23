package api

// #cgo pkg-config: mkl-static-ilp64-tbb
// #include <mkl.h>
import "C"

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type _Feature struct {
	red, green, blue []float32
}

func extractColoredFeature(imageData *ImageData, startX, startY, width int) *_Feature {
	feature := _Feature{
		red:   make([]float32, width*width),
		green: make([]float32, width*width),
		blue:  make([]float32, width*width),
	}

	featureIndex := 0
	for y := startY; y < startY+width; y++ {
		for x := startX; x < startX+width; x++ {
			index := 4 * (x + y*imageData.Width)

			feature.red[featureIndex] = float32(imageData.Data[index])
			feature.green[featureIndex] = float32(imageData.Data[index+1])
			feature.blue[featureIndex] = float32(imageData.Data[index+2])
			featureIndex++
		}
	}

	return &feature
}

func applyGaussianBlur(imageData *ImageData) *ImageData {
	KERNEL := []float32{
		0.0000, 0.0000, 0.0000, 0.0001, 0.0001, 0.0001, 0.0000, 0.0000, 0.0000,
		0.0000, 0.0000, 0.0004, 0.0014, 0.0023, 0.0014, 0.0004, 0.0000, 0.0000,
		0.0000, 0.0004, 0.0037, 0.0146, 0.0232, 0.0146, 0.0037, 0.0004, 0.0000,
		0.0001, 0.0014, 0.0146, 0.0584, 0.0926, 0.0584, 0.0146, 0.0014, 0.0001,
		0.0001, 0.0023, 0.0232, 0.0926, 0.1466, 0.0926, 0.0232, 0.0023, 0.0001,
		0.0001, 0.0014, 0.0146, 0.0584, 0.0926, 0.0584, 0.0146, 0.0014, 0.0001,
		0.0000, 0.0004, 0.0037, 0.0146, 0.0232, 0.0146, 0.0037, 0.0004, 0.0000,
		0.0000, 0.0000, 0.0004, 0.0014, 0.0023, 0.0014, 0.0004, 0.0000, 0.0000,
		0.0000, 0.0000, 0.0000, 0.0001, 0.0001, 0.0001, 0.0000, 0.0000, 0.0000,
	}
	EQ := float32(0.0)
	KERNEL_WIDTH := 11
	INC := int(KERNEL_WIDTH / 2)
	EL_NUM := C.longlong(KERNEL_WIDTH * KERNEL_WIDTH)

	for _, v := range KERNEL {
		EQ += v
	}
	if EQ < 1 {
		EQ = 1
	}

	blurredImageData := ImageData{}
	blurredImageData.Width = imageData.Width
	blurredImageData.Height = imageData.Height
	blurredImageData.Data = make([]uint8, len(imageData.Data))
	copy(blurredImageData.Data, imageData.Data)

	wg := sync.WaitGroup{}
	wg.Add(imageData.Height - KERNEL_WIDTH)

	for y := 0; y < imageData.Height-KERNEL_WIDTH; y++ {
		go func(y int) {
			for x := 0; x < imageData.Width-KERNEL_WIDTH; x++ {
				result := make([]float32, EL_NUM)
				updateIndex := 4 * ((x + INC) + (y+INC)*imageData.Width)
				feature := extractColoredFeature(imageData, x, y, KERNEL_WIDTH)

				// Red channel
				C.vsMul(EL_NUM, (*C.float)(&KERNEL[0]), (*C.float)(&feature.red[0]), (*C.float)(&result[0]))
				blurredImageData.Data[updateIndex] = uint8(C.cblas_sasum(EL_NUM, (*C.float)(&result[0]), 1) / C.float(EQ))

				// Green channel
				C.vsMul(EL_NUM, (*C.float)(&KERNEL[0]), (*C.float)(&feature.green[0]), (*C.float)(&result[0]))
				blurredImageData.Data[updateIndex+1] = uint8(C.cblas_sasum(EL_NUM, (*C.float)(&result[0]), 1) / C.float(EQ))

				// Blue channel
				C.vsMul(EL_NUM, (*C.float)(&KERNEL[0]), (*C.float)(&feature.blue[0]), (*C.float)(&result[0]))
				blurredImageData.Data[updateIndex+2] = uint8(C.cblas_sasum(EL_NUM, (*C.float)(&result[0]), 1) / C.float(EQ))
			}
			wg.Done()
		}(y)
	}

	wg.Wait()

	return &blurredImageData
}

func GaussianBlur(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	imageData := ImageData{}

	requestData, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println("Could not read request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// log.Println(string(requestData))

	err = json.Unmarshal(requestData, &imageData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	start := time.Now()
	blurredImage := applyGaussianBlur(&imageData)
	elapsed := time.Since(start)

	log.Println("Blurring took: ", elapsed.Milliseconds(), "ms")

	responseData, err := json.Marshal(blurredImage)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	bytesWritten, err := w.Write(responseData)
	if err != nil {
		log.Println("Could not respond back")
		log.Println(err)
	}

	if bytesWritten != len(responseData) {
		log.Println("Response byte slice had", len(responseData), "bytes. Only", bytesWritten, "bytes were written")
	}
}
