package youtube_api

import (
	"fmt"
	"testing"
)

func TestVideoIDs(t *testing.T) {
	search := SearchListResponse{
		Items: []searchItem{
			{ID: id{"abcd"}},
			{ID: id{"efgh"}},
			{ID: id{"abcd"}},
		},
	}
	ids := search.VideoIDs()
	if !contains(ids, "abcd") || len(ids) != 2 {
		t.Fail()
	}
}

func TestSearchURL(t *testing.T) {
	channel := NewClient("UC-ocfaUlccs9DcrBRzwlI-g")
	channel.APIKey = "test-key"
	actual := channel.SearchURL("")
	expected := "https://www.googleapis.com/youtube/v3/search?channelId=UC-ocfaUlccs9DcrBRzwlI-g&key=test-key&maxResults=50&order=date&part=snippet&type=video"
	if expected != actual {
		fmt.Println(actual, "is not expected. failed")
		t.Fail()
	}
}

func contains(s []string, st string) bool {
	for _, v := range s {
		if v == st {
			return true
		}
	}
	return false
}
