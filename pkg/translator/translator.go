package translator

import (
	gt "github.com/bas24/googletranslatefree"
)

func Translate(text string) (string, error) {
	return gt.Translate(text, "zh-TW", "en-US")
}
