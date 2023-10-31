package utils

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
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

	dst, err := imaging.Open("public/img/mask.jpg")
	if err != nil {
		return
	}

	dst = imaging.Resize(dst, 1920, 1080, imaging.Lanczos)
	dst = imaging.Blur(dst, 5)

	//rand.Seed(time.Now().UnixNano())
	max := dst.Bounds().Max
	count := rand.Intn(35) + 25
	adds := 0
	for adds < count {
		bcp := getBookCoverPath()
		if bcp == "" {
			continue
		}
		//fmt.Println("Open", bcp)
		src, err := imaging.Open(bcp)
		if err != nil {
			return
		}

		szW := rand.Intn(580-150) + 150
		szH := rand.Intn(912-150) + 150
		ang := rand.Float64() * 360

		x := int(rand.Float64()*float64(max.X) - float64(szW)/2)
		y := int(rand.Float64()*float64(max.Y) - float64(szH)/2)
		opacity := 0.7 // rand.Float64()
		if opacity < 0.2 {
			opacity = 0.2
		} else if opacity > 0.7 {
			opacity = 0.7
		}

		src = imaging.Fit(src, szW, szH, imaging.Lanczos)
		src = imaging.Rotate(src, ang, color.Transparent)
		dst = imaging.Overlay(dst, src, image.Pt(x, y), opacity)
		adds++
		//fmt.Println(adds, "/", count)
	}

	err = imaging.Save(dst, filepath.Join("public/img/back.png"))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func getBookCoverPath() string {
	booksDir := "db/storage/books"
	files, err := os.ReadDir(booksDir)
	if err != nil {
		return ""
	}
	i := rand.Intn(len(files))
	bookDir := filepath.Join(booksDir, files[i].Name())
	files, err = os.ReadDir(bookDir)
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
