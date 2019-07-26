package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/yansal/youtube-ar/log"
)

// Wrap wraps an http.Client and logs all roundtrips.
func Wrap(client *http.Client, log log.Logger) *http.Client {
	new := *client
	new.Transport = &transport{
		log: log,
		rt:  client.Transport,
	}
	return &new
}

type transport struct {
	rt  http.RoundTripper
	log log.Logger
}

func (tr *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp  *http.Response
		rterr error
		start = time.Now()
	)

	defer func() {
		fields := []log.Field{
			log.Stringer("duration", time.Since(start)),
			log.Stringer("url", req.URL), // TODO: remove private data such as api keys
		}
		if rterr != nil {
			tr.log.Log(req.Context(), rterr.Error(), fields...)
			return
		}

		if resp.StatusCode/100 == 4 || resp.StatusCode/100 == 5 {
			defer func() { _ = resp.Body.Close() }()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				tr.log.Log(req.Context(), rterr.Error(), fields...)
				return
			}
			fields = append(fields, log.String("response-body", string(body)))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		fields = append(fields, log.Int("code", resp.StatusCode))
		tr.log.Log(req.Context(), req.Method+" "+req.URL.Path,
			fields...)
	}()

	resp, rterr = http.DefaultTransport.RoundTrip(req)
	return resp, rterr
}
