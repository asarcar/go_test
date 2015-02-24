package reddit

import "testing"

func TestReddit(t *testing.T) {
	const in = "golang"
	if items, err := Get(in); err != nil {
		t.Errorf("#items %d: error %s", len(items), err)
	}
}
