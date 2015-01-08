package gimg

import (
	_ "fmt"
)

type ZStorage interface {
	SaveImage(data []byte) (string, error)
	//NewImage(saveName string, data []byte) error
	GetImage(request *ZRequest) ([]byte, error)
	InfoImage(md5 string) (*ZImageInfo, error)
}
