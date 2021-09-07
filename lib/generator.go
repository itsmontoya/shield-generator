package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	// ErrEmptyFirstPart is returned when the first part of a badge is empty
	ErrEmptyFirstPart = errors.New("invalid first part, cannot be empty")
	// ErrEmptySecondPart is returned when the second part of a badge is empty
	ErrEmptySecondPart = errors.New("invalid second part, cannot be empty")
	// ErrEmptyColor is returned when the color of a badge is empty
	ErrEmptyColor = errors.New("invalid color, cannot be empty")
)

const urlFmt = "https://img.shields.io/badge/%s-%s-%s"

func GenerateURL(first, second, color string) (generated string, err error) {
	switch {
	case len(first) == 0:
		err = ErrEmptyFirstPart
		return

	case len(second) == 0:
		err = ErrEmptySecondPart
		return
	case len(color) == 0:
		err = ErrEmptyColor
		return
	}

	first = url.QueryEscape(first)
	second = url.QueryEscape(second)
	generated = fmt.Sprintf(urlFmt, first, second, color)
	return
}

func GenerateImage(url string, w io.Writer) (err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		err = handleError(resp.Body)
		err = fmt.Errorf("error downloading image with URL <%s>: %s", url, err)
		return
	}

	_, err = io.Copy(w, resp.Body)
	return
}

func handleError(r io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, r); err != nil {
		return
	}

	err = errors.New(buf.String())
	return
}
