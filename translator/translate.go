package translator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/robertkrimen/otto"
)

func encodeURI(s string) (string, error) {
	eUri := `eUri = encodeURI(sourceText);`
	vm := otto.New()
	err := vm.Set("sourceText", s)
	if err != nil {
		return "err", errors.New("error setting js variable")
	}
	_, err = vm.Run(eUri)
	if err != nil {
		return "err", errors.New("error executing jscript")
	}
	val, err := vm.Get("eUri")
	if err != nil {
		return "err", errors.New("error getting variable value from js")
	}
	v, err := val.ToString()
	if err != nil {
		return "err", errors.New("error converting js var to string")
	}
	return v, nil
}

func Translate(source, sourceLang, targetLang string) (string, error) {
	var text []string
	var result []interface{}

	encodedSource, err := encodeURI(source)
	if err != nil {
		return "err", err
	}
	url := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=" +
		sourceLang + "&tl=" + targetLang + "&dt=t&q=" + encodedSource

	r, err := http.Get(url)
	if err != nil {
		return "err", errors.New("error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "err", errors.New("error reading response body")
	}

	bReq := strings.Contains(string(body), `<title>Error 400 (Bad Request)`)
	if bReq {
		return "err", errors.New("error 400 (Bad Request)")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "err", errors.New("error unmarshaling data")
	}

	if len(result) > 0 {
		inner := result[0]
		for _, slice := range inner.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		cText := strings.Join(text, "")

		return cText, nil
	} else {
		return "err", errors.New("no translated data in response")
	}
}
