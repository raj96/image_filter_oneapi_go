package main

import (
	"image_filter_oneapi/api"
	"log"
	"net/http"
)

func main() {

	rootMux := http.NewServeMux()

	uiHandler := http.FileServer(http.Dir("./ui"))
	rootMux.Handle("/", uiHandler)

	rootMux.HandleFunc("/oneapi/grayscale", api.Grayscale)
	rootMux.HandleFunc("/oneapi/gaussian_blur", api.GaussianBlur)

	log.Println("Listening at :8080")
	err := http.ListenAndServe(":8080", rootMux)
	if err != nil {
		log.Fatalln(err)
	}

	// input := []float32{1, 2, 3, 4, 5}
	// res := C.cblas_sasum(C.longlong(len(input)), (*C.float)(&input[0]), 1)

	// fmt.Println("Result: ", res)
}
