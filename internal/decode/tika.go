/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
