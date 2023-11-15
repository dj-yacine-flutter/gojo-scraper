package utils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/buckket/go-blurhash"
)

func GetBlurHash(base, path string) (string, error) {
	if base == "" {
		return "", fmt.Errorf("image URL is empty")
	}

	response, err := http.Get(base + path)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("URL does not point to an image")
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return "", err
	}

	blurHash, err := blurhash.Encode(4, 3, img)
	if err != nil {
		return "", err
	}

	return blurHash, nil
}
