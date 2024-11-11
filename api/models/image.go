package models

import "os"

type ImageResponse struct {
	ImageMetaData
	File *os.File
	// DownLiadLink string `json:"link"`
}
