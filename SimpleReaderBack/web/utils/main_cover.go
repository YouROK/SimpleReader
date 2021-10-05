package utils

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"

	"SimpleReader/web/settings"
)

func MakeCover() {
	makeCover()
	timer := time.NewTimer(time.Minute * 15)
	for {
		<-timer.C
		makeCover()
	}
}

func makeCover() {

	dst, err := imaging.Open(filepath.Join(settings.Path, "img", "mask.jpg"))
	if err != nil {
		return
	}

	dst = imaging.Resize(dst, 1920, 1080, imaging.Lanczos)
	dst = imaging.Blur(dst, 5)

	rand.Seed(time.Now().UnixNano())
	max := dst.Bounds().Max
	count := rand.Intn(50) + 15

	for i := 0; i < count; i++ {
		src, err := imaging.Open(getBookCoverPath())
		if err != nil {
			return
		}

		szW := rand.Intn(580-150) + 150
		szH := rand.Intn(912-150) + 150
		ang := rand.Float64()*60 - 30

		x := int(rand.Float64()*float64(max.X) - float64(szW))
		y := int(rand.Float64()*float64(max.Y) - float64(szH))
		opacity := rand.Float64()
		if opacity < 0.2 {
			opacity = 0.2
		} else if opacity > 0.7 {
			opacity = 0.7
		}

		src = imaging.Fit(src, szW, szH, imaging.Lanczos)
		src = imaging.Rotate(src, ang, color.Transparent)
		dst = imaging.Overlay(dst, src, image.Pt(x, y), opacity)
		fmt.Println(i+1, "/", count)
	}

	err = imaging.Save(dst, filepath.Join(settings.Path, "img", "back.png"))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

var randomizeCover = map[string]struct{}{}

func getBookCoverPath() string {
	booksDir := filepath.Join(settings.Path, "books")
	files, err := ioutil.ReadDir(booksDir)
	if err != nil {
		return ""
	}
	i := rand.Intn(len(files))
	bookDir := filepath.Join(booksDir, files[i].Name())
	files, err = ioutil.ReadDir(bookDir)
	if err != nil {
		return ""
	}
	for _, f := range files {
		if !f.IsDir() && (filepath.Ext(f.Name()) == ".png" || filepath.Ext(f.Name()) == ".jpg") {
			return filepath.Join(bookDir, f.Name())
		}
	}
	return ""
}
