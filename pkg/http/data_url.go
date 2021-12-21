package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vincent-petithory/dataurl"
)

func GetDataURL(ctx context.Context, url string) (string, string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf(
			"%d: %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, resp.Body); err != nil {
		return "", "", err
	}

	ct := resp.Header.Get("Content-Type")
	if len(ct) == 0 {
		ct = http.DetectContentType(b.Bytes())
	}
	return dataurl.New(b.Bytes(), ct).String(), ct, nil
}
