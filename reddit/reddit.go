//
// Example: Go Tutorial
//

// Package reddit implements a basic client for Reddit API.
package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	kURL = "http://reddit.com/r/%s.json"
)

// Reddit API JSON
// {
//   "data":
//   {
//     "children":
//     [
//       {
//         "data":
//         {
//           "title"       : "The Go homepage",
//           "url"         : "http://golang.org",
//           "num_comments": 10
//           ...
//         }
//       },
//       ...
//     ]
//   }
// }

// Item describes a Reddit items.
type Item struct {
	Title    string
	URL      string
	Comments int `json:"num_comments"`
}

// struct tag: annotates field. Go uses reflect package
// to inspect this at runtime. Annotation tells json
// package to decode "num_comments" field of JSON object
// into Comments field and vice versa

type response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

// Get fetches the most recent items posted to reddit string (subreddit)
func Get(reddit string) ([]Item, error) {
	url := fmt.Sprintf(kURL, reddit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	// Decode the JSON objected into r
	r := new(response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(r.Data.Children))
	for i, child := range r.Data.Children {
		items[i] = child.Data
	}

	return items, nil
}

// Implements: Stringer interface so that it can be used directly in fmt.Print...
func (i Item) String() string {
	com := ""
	switch i.Comments {
	case 0:
	case 1:
		com = "[1 comment]"
	default:
		com = fmt.Sprintf("[%d comments]", i.Comments)
	}
	return fmt.Sprintf("%s %s\n\t%s", i.Title, com, i.URL)
}
