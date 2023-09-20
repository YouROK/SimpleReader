package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"SimpleReader/server_old/models/users"
)

var captchas = make(map[*users.User]int)

func GenerateCaptcha(u *users.User) string {
	rand.Seed(time.Now().Unix())
	b, err := os.Open("public/img/book_36.png")
	if err != nil {
		log.Println("Error open book img file", err)
		return ""
	}
	defer b.Close()
	captVer := rand.Intn(6) + 4
	book, _, err := image.Decode(b)
	if err != nil {
		log.Println("Error decode book img file", err)
		return ""
	}
	canvas := image.NewRGBA(image.Rect(0, 0, 32*captVer, 150))
	x := 0
	y := rand.Intn(150 - 64)
	for i := 0; i < captVer; i++ {
		draw.Draw(canvas, book.Bounds().Add(image.Point{x, y}), book, image.Point{0, 0}, draw.Over)
		x = x + (canvas.Bounds().Dx() / captVer)
		y = rand.Intn(150 - 64)
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, canvas)
	if err == nil {
		captchas[u] = captVer
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	} else {
		log.Println("Error convert to img to png", err)
		return ""
	}
}

func VerifyCaptcha(u *users.User, captcha string) bool {
	capint, err := strconv.Atoi(captcha)
	if err != nil {
		return false
	}
	truecapt, ok := captchas[u]
	if ok {
		delete(captchas, u)
		return truecapt == capint
	}
	return false
}
