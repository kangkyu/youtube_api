package youtube_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func VideoListFromVideosURL(vu *url.URL) (VideoListResponse, error) {

	var videoList = VideoListResponse{}

	resp, err := http.Get(vu.String())
	if err != nil {
		return videoList, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&videoList)
	if err != nil {
		return videoList, err
	}

	return videoList, nil
}

func SearchListFromSearchURL(su *url.URL) (SearchListResponse, error) {

	var searchList = SearchListResponse{}
	// fmt.Printf("%s\n", su.String())

	resp, err := http.Get(su.String())
	if err != nil {
		return searchList, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return searchList, err
	}

	if resp.StatusCode != 200 {
		return searchList, fmt.Errorf("youtube error response: %v", string(body))
	}

	err = json.Unmarshal(body, &searchList)
	if err != nil {
		return searchList, err
	}

	return searchList, nil
}

func SearchURL(nextPageToken string, uuid string) *url.URL {

	u, _ := url.Parse("https://www.googleapis.com/youtube/v3/search")

	v := url.Values{}
	v.Set("key", os.Getenv("YT_API_KEY"))
	v.Add("part", "snippet")
	v.Add("type", "video")
	v.Add("maxResults", "50")
	v.Add("order", "date")
	v.Add("channelId", uuid)

	if len(nextPageToken) != 0 {
		v.Set("pageToken", nextPageToken)
	}

	u.RawQuery = v.Encode()
	return u
}

func VideosURL(videoIDs []string) *url.URL {

	u, _ := url.Parse("https://www.googleapis.com/youtube/v3/videos")

	v := url.Values{}
	v.Set("key", os.Getenv("YT_API_KEY"))
	v.Add("part", "snippet,statistics,status,contentDetails")
	v.Add("id", strings.Join(videoIDs, ","))

	u.RawQuery = v.Encode()
	return u
}
