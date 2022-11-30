package puzzleinput

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

const (
	urlFormat      = "https://adventofcode.com/2021/day/%d/input"
	filePathFormat = "days/%02d/input.txt"
)

func Get(day int) (io.ReadCloser, error) {

	fp := fmt.Sprintf(filePathFormat, day)

	if existsFile(fp) {
		err := download(day, fp)
		if err != nil {
			return nil, err
		}
	}

	return os.Open(fp)
}

func existsFile(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrExist)
}

func getCookie() (_ string, err error) {
	f, err := os.Open("cookie.txt")
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
