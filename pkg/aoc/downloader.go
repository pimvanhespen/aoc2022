package aoc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	urlFormat      = "https://adventofcode.com/2022/day/%d/input"
	filePathFormat = "days/%02d/input.txt"
)

func Get(day int) (io.Reader, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fp := fmt.Sprintf(filePathFormat, day)

	fp = filepath.Join(wd, fp)

	if !existsFile(fp) {
		err = download(day, fp)
		if err != nil {
			return nil, err
		}
	}

	reader, writer := io.Pipe()
	go func() {

		file, openErr := os.Open(fp)
		if openErr != nil {
			_ = writer.CloseWithError(openErr)
			return
		}
		defer func() {
			_ = file.Close()
		}()

		_, cErr := io.Copy(writer, file)
		if cErr != nil {
			_ = writer.CloseWithError(cErr)
			return
		}

		_ = writer.Close()
	}()

	return reader, nil
}

func Load[Result any](day int, parseFunc func(io.Reader) (Result, error)) (Result, error) {
	r, err := Get(day)
	if err != nil {
		return *new(Result), err
	}

	return parseFunc(r)
}

func existsFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getCookie() (_ string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fp := filepath.Join(wd, "cookie.txt")

	f, err := os.Open(fp)
	if err != nil {
		return "", err
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func download(day int, destination string) error {
	remote := fmt.Sprintf(urlFormat, day)

	req, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		return err
	}

	cookie, err := getCookie()
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", fmt.Sprintf("session=%s", cookie))

	resp, err := http.DefaultClient.Do(req)

	defer func() {
		cerr := resp.Body.Close()
		if err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	f, err := os.Create(destination)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
