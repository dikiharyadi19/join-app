package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github/yogabagas/join-app/shared/constant"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/matthewhartstonge/argon2"
	ulid "github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func NewULIDGenerate() string {
	defaultEntropySource := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), defaultEntropySource).String()
}

func Hash(alg, pwd string) ([]byte, error) {

	switch alg {
	case constant.Bcrypt.String():
		h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		return h, nil
	case constant.MD5.String():
		h := md5.Sum([]byte(pwd))
		return h[:], nil
	case constant.Argon.String():
		conf := argon2.DefaultConfig()
		h, err := conf.HashEncoded([]byte(pwd))
		if err != nil {
			return nil, err
		}
		return h, nil
	case constant.SHA.String():
		h := sha256.Sum256([]byte(pwd))
		return h[:], nil
	default:
		return nil, errors.New("[CLIENT] - Unsupported algorithm")
	}

}

func Base64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func ValidateEmail(email string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return emailRegex.MatchString(email)
}

func PageToOffset(limit int, page int) int {

	if page <= 1 {
		return 0
	}

	return (limit * page) - limit
}

func GetTotalPage(totalData int, perPage int) int {
	return (totalData + perPage - 1) / perPage
}

func GetFileContentType(out *os.File) string {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return ""
	}
	contentType := http.DetectContentType(buffer)

	return contentType
}
