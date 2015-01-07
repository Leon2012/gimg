package gimg

import (
	_ "fmt"
	"github.com/gographics/imagick/imagick"
)

type ZImage struct {
	MW *imagick.MagickWand
}

func NewImage() *ZImage {
	imagick.Initialize()
	return &ZImage{MW: imagick.NewMagickWand()}
}

func (z *ZImage) Destroy() {
	imagick.Terminate()
	z.MW.Destroy()
}
