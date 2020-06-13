package main

import (
	"fmt"
	"github.com/urfave/cli"
	"image"
	"image/color"
	"image/gif"
	"os"
	"strconv"
	"strings"
)

var (
	blackAndWhitePalette = []color.Color{color.White, color.Black}
	defaultBlackLine = 300
	dot = "."
	defaultThreadNum = 4
)


type GifContext struct {
	BlackLine   int
	Palette     []color.Color
	GifFilePath string
	ThreadNum int

}

type GifImageHandler interface {
	NewImage(ctx GifContext, oldImage image.Paletted) *image.Paletted
	ConvertColor(ctx GifContext, oldColor color.Color) color.Color
	FileSubFix() string
}


func main() {

	var gifFilePath string
	var blackLine, threadNum int
	app := cli.NewApp()
	app.Name = "Gif转换"
	app.Usage = "把Gif进行转换"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "进行黑白转换",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "file, f",
					Usage:       "Gif文件名称",
					Destination: &gifFilePath,
				},
				cli.IntFlag{
					Name:        "black, b",
					Usage:       "转换为黑色的阈值（源像素点的rgb和超过此值时为黑色，否则为白色）",
					Destination: &blackLine,
					Value: defaultBlackLine,
				},
				cli.IntFlag{
					Name:        "thread, t",
					Usage:       "线程数",
					Destination: &threadNum,
					Value: defaultThreadNum,
				},
			},
			Action: func(c *cli.Context) error {
				if gifFilePath == ""{
					panic("源文件不能为空！")
				}

				return convertGif(GifContext{BlackLine:blackLine, GifFilePath:gifFilePath, ThreadNum:threadNum, Palette:blackAndWhitePalette}, &GifColorPixel{})
			},
		},
		{
			Name:    "reverse",
			Aliases: []string{"r"},
			Usage:   "进行颜色反转",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "file, f",
					Usage:       "Gif文件名称",
					Destination: &gifFilePath,
				},
				cli.IntFlag{
					Name:        "thread, t",
					Usage:       "线程数",
					Destination: &threadNum,
					Value: defaultThreadNum,
				},
			},
			Action: func(c *cli.Context) error {
				if gifFilePath == ""{
					panic("源文件不能为空！")
				}

				return convertGif(GifContext{BlackLine:blackLine, GifFilePath:gifFilePath, ThreadNum:threadNum, Palette:blackAndWhitePalette}, &GifColorReverse{})
			},
		},
	}
	// 接受os.Args启动程序
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func convertGif(ctx GifContext, handler GifImageHandler)  error{
	gifFilePath := ctx.GifFilePath

	dotIndex := strings.LastIndex(gifFilePath, dot)
	subFix := gifFilePath[dotIndex+1:]

	if subFix != "gif" {
		panic("文件后缀需要为gif")
	}
	resultFilePath := gifFilePath[0:dotIndex] + "_result_" + handler.FileSubFix() + ".gif"

	file, err := os.Open(gifFilePath)
	defer file.Close()
	if err != nil{
		panic(err)
	}

	sourceGif, _ := gif.DecodeAll(file)
	newGif := gif.GIF{LoopCount: sourceGif.LoopCount, Delay:sourceGif.Delay}

	imageChans := *new([]chan []*image.Paletted)

	sourceImages := sourceGif.Image
	sourceImageLen := len(sourceImages)

	threadNum := ctx.ThreadNum
	if threadNum > sourceImageLen {
		threadNum = sourceImageLen
	}
	perThreadImageNum := sourceImageLen / threadNum

	progressChan := make(chan int, sourceImageLen)
	go showProgress(sourceImageLen, progressChan)

	for i := 0; i < threadNum; i++ {

		newChan := make(chan []*image.Paletted, 1)
		imageChans = append(imageChans, newChan)
		if i == threadNum - 1 {
			go getImageAsync(ctx,handler,sourceImages[perThreadImageNum * i:], newChan, progressChan)
		}else {
			go getImageAsync(ctx,handler,sourceImages[perThreadImageNum * i: perThreadImageNum * (i + 1)], newChan, progressChan)
		}
	}

	for _,imageChan := range imageChans{
		newImages := <- imageChan
		newGif.Image = append(newGif.Image, newImages...)
	}

	close(progressChan)

	out,err := os.Create(resultFilePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = gif.EncodeAll(out, &newGif)
	if err != nil {
		panic(err)
	}
	return nil
}

func getImageAsync(ctx GifContext, handler GifImageHandler, oldImages []*image.Paletted, imageChan chan []*image.Paletted, progressChan chan int) {
	images := *new([]*image.Paletted)
	for _, img := range oldImages {
		newImg := handler.NewImage(ctx, *img)
		rect := newImg.Rect
		x := rect.Max.X
		y := rect.Max.Y
		for i := 0; i < x ; i += 1{
			for j:= 0; j < y ; j += 1{
				c := img.At(i,j)
				newImg.Set(i, j, handler.ConvertColor(ctx,c))
			}
		}
		images = append(images, newImg)
		progressChan <- 1
	}
	imageChan <- images
}

func showProgress(totalImageCount int, progressChan chan int)  {
	successCount := 0
	progrssStr := ""
	progressNum := 30
	percentage := 0.0
	for {
		progrssStr = ""
		percentage = float64(successCount * 1.0) / float64(totalImageCount)
		for i := 0; i < progressNum; i++ {
			if percentage >= float64(i * 1.0) / float64(progressNum){
				progrssStr += "#"
			}else {
				progrssStr += "*"
			}
		}
		progrssStr +=  strconv.FormatFloat(float64(percentage * 100), 'f', 6, 64) + "%" +"(" + strconv.Itoa(successCount) + ")"
		fmt.Printf("%s\r", progrssStr)
		n, isClose := <- progressChan
		if !isClose {
			break
		}
		successCount += n
	}
	fmt.Println()
}
