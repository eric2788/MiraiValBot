package cosplayer

import "testing"

func TestGetOvooa(t *testing.T) {
	provider := providers["ovooa"]

	data, err := provider.GetImages()

	if err != nil {
		t.Skip(err)
	}

	t.Log(data.Title)
	for _, url := range data.Urls {
		t.Log(url)
	}
}
