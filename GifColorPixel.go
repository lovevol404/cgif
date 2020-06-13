package main

import (
	"image"
	"image/color"
)

type GifColorPixel struct {

}

func (sf *GifColorPixel) NewImage(ctx GifContext, oldImage image.Paletted) *image.Paletted  {
	return image.NewPaletted(oldImage.Rect, ctx.Palette)
}

func (sf *GifColorPixel) ConvertColor(ctx GifContext, oldColor color.Color) color.Color  {
	r,b,g,_ := oldColor.RGBA()
	r_ := r >> 8
	b_ := b >> 8
	g_ := g >> 8

	if int(r_+ b_ + g_) > ctx.BlackLine {
		return color.Black
	}else {
		return color.White
	}
}

func (sf *GifColorPixel) FileSubFix() string  {
	return "c"
}