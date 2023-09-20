package captcha

import (
	"SimpleReader/server/models/user"
	"github.com/dchest/captcha"

	"bytes"
	"encoding/base64"
)

var ids = make(map[*user.User]string)

func GenerateCaptcha(u *user.User) string {
	id := captcha.New()
	buf := bytes.NewBufferString("")
	captcha.WriteImage(buf, id, 128, 64)
	ids[u] = id
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func VerifyCaptcha(u *user.User, nums string) bool {
	if id, ok := ids[u]; ok {
		return captcha.VerifyString(id, nums)
	}
	return false
}
