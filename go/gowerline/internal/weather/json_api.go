package weather

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/johnstarich/go/gowerline/internal/httpclient"
	"github.com/pkg/errors"
)

func doJSONGet(ctx context.Context, httpClient httpclient.Client, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("failed to fetch from %q: %s", url, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(result)
}
