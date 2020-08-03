// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package decode

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
