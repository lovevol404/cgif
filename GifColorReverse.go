package main

import (
	"image"
	"image/color"
)

type GifColorReverse struct {

}

func (sf *GifColorReverse) NewImage(ctx GifContext, oldImage image.Paletted) *image.Paletted  {
	return image.NewPaletted(oldImage.Rect, oldImage.Palette)
}

func (sf *GifColorReverse) ConvertColor(ctx GifContext, oldColor color.Color) color.Color  {
	r,b,g,a := oldColor.RGBA()
	r_ := r >> 8
	g_ := g >> 8
	b_ := b >> 8
	return color.RGBA{R: uint8(255 - r_), G: uint8(255 - g_), B: uint8(255 - b_), A: uint8(a)}
}
func (sf *GifColorReverse) FileSubFix() string  {
	return "r"
}

