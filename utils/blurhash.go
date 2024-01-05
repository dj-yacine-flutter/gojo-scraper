package utils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/buckket/go-blurhash"
	"github.com/nfnt/resize"
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

	m := resize.Resize(500, 250, img, resize.Lanczos3)

	blurHash, err := blurhash.Encode(4, 3, m)
	if err != nil {
		return "", err
	}

	return blurHash, nil
}
