// Package search provides a function to do Google searches using the Google Web
// Search API. See https://developers.google.com/web-search/docs/
package backend

import (
	"encoding/json"
	"net/http"

	pb "github.com/asarcar/go_test/search/protos"

	"golang.org/x/net/context"
)

const (
	kGetCmdStr             = "GET"
	kGoogleSearchApiUrlStr = "https://ajax.googleapis.com/ajax/services/search/web?v=1.0"
)

// Search sends query to Google search and returns the results.
func Search(ctx context.Context, query string) (*pb.Results, error) {
	// Prepare the Google Search API request.
	req, err := http.NewRequest(kGetCmdStr, kGoogleSearchApiUrlStr, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)

	// If ctx is carrying the user IP address, forward it to the server.
	// Google APIs use the user IP to distinguish server-initiated requests
	// from end-user requests.
	if userIP, ok := FromContext(ctx); ok {
		q.Set("userip", userIP.String())
	}
	req.URL.RawQuery = q.Encode()

	// Issue the HTTP request and handle the response. The httpDo function
	// cancels the request if ctx.Done is closed.
	var results pb.Results
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Parse the JSON search result.
		// https://developers.google.com/web-search/docs/#fonje
		var data struct {
			ResponseData struct {
				Results []struct {
					TitleNoFormatting string
					URL               string
					Content           string
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		for _, res := range data.ResponseData.Results {
			results.Res = append(results.Res,
				&pb.Result{
					Title:   res.TitleNoFormatting,
					Url:     res.URL,
					Content: res.Content,
				})
		}
		return nil
	})
	// httpDo waits for the closure we provided to return, so it's safe to
	// read results here.
	return &results, err
}

// httpDo issues the HTTP request and calls f with the response. If ctx.Done is
// closed while the request or f is running, httpDo cancels the request, waits
// for f to exit, and returns ctx.Err. Otherwise, httpDo returns f's error.
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
