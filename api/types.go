package api

type ImageData struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Data   []uint8 `json:"data"`
}
