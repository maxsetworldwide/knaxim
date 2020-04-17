package process

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context/ctxhttp"
)

func tikaTextExtract(ctx context.Context, input io.Reader, path string) (io.ReadCloser, error) {
	req, err := http.NewRequest("PUT", path+"/tika", input)
	if err != nil {
		return nil, err
	}
	resp, err := ctxhttp.Do(ctx, http.DefaultClient, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tika response code %v", resp.StatusCode)
	}
	return resp.Body, nil
}
