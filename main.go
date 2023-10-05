package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func fetchHTML(URL string) (io.ReadCloser, error) {
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code error: %d", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("content type not text/html")
	}

	return resp.Body, nil
}

func fetchHyperLink(body io.ReadCloser, startTag, class string) (string, error) {

	tokenizer := html.NewTokenizer(body)

	defer body.Close()

	for {

		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return "", fmt.Errorf("tag \"%v\" with class \"%v\" not found", startTag, class)
			}
			return "", fmt.Errorf("error tokenizing html: %v", tokenizer.Err())
		}

		if tokenType != html.StartTagToken {
			continue
		}

		token := tokenizer.Token()
		if token.Data != startTag {
			continue
		}

		attr := token.Attr
		if attr[1].Val != class {
			continue
		}

		return attr[0].Val, nil
	}
}

func fetchLink(URL, startTag, class string) (string, error) {
	body, err := fetchHTML(URL)
	if err != nil {
		return "", err
	}
	defer body.Close()

	hl, err := fetchHyperLink(body, startTag, class)
	if err != nil {
		return "", err
	}

	return hl, nil
}

func downloadImage(URL, reg string) error {

	pageLink, err := fetchLink(URL, "a", "result__photoLink")
	if err != nil {
		return err
	}

	pageLink = "http://www.jetphotos.com" + pageLink

	imgSrc, err := fetchLink(pageLink, "img", "large-photo__img")
	if err != nil {
		return err
	}

	dir := filepath.Join(".", "img")

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return err
	}

	fdir := filepath.Join(dir, reg+".jpg")

	f, err := os.Create(fdir)

	if err != nil {
		return err
	}

	defer f.Close()

	res, err := http.Get(imgSrc)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:\nsearcher <registration>")
		os.Exit(1)
	}

	URL := "http://www.jetphotos.com/photo/keyword"

	for i := 1; i < len(os.Args); i++ {
		reg := os.Args[i]
		searchURL := URL + "/" + reg
		err := downloadImage(searchURL, reg)

		if err != nil {
			log.Fatal(err)
		}
	}
}
