package gimg

import (
	_ "fmt"
	"github.com/gographics/imagick/imagick"
)

type ZImage struct {
	MW *imagick.MagickWand
}

func NewImage() *ZImage {
	return &ZImage{MW: imagick.NewMagickWand()}
}

func (z *ZImage) Destroy() {
	z.MW.Destroy()
}
