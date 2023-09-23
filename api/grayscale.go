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
)

func applyGrayscale(imageData *ImageData) {
	wg := sync.WaitGroup{}
	wg.Add(imageData.Height)

	for y := 0; y < imageData.Height; y++ {
		go func(y int) {
			splicedData := make([]float32, 3)
			for x := 0; x < imageData.Width; x++ {
				index := 4 * (x + y*imageData.Width)
				for i, v := range imageData.Data[index : index+3] {
					splicedData[i] = float32(v)
				}

				avg := uint8(C.cblas_sasum(3, (*C.float)(&splicedData[0]), 1) / 3.0)

				imageData.Data[index] = avg
				imageData.Data[index+1] = avg
				imageData.Data[index+2] = avg
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
}

func Grayscale(w http.ResponseWriter, r *http.Request) {
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

	applyGrayscale(&imageData)

	responseData, err := json.Marshal(imageData)
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
