package youtube_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ChannelClient struct {
	BaseURL    string
	APIKey     string
	ChannelID  string
	HTTPClient *http.Client
}

func NewClient(channelID string) *ChannelClient {
	return &ChannelClient{
		BaseURL:   "https://www.googleapis.com/youtube/v3/",
		APIKey:    os.Getenv("YT_API_KEY"),
		ChannelID: channelID,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, // Set a reasonable timeout
		},
	}
}

func (c *ChannelClient) SearchURL(nextPageToken string) string {
	u, _ := url.Parse(fmt.Sprintf("%ssearch", c.BaseURL))

	v := url.Values{}
	v.Set("key", c.APIKey)
	v.Add("part", "snippet")
	v.Add("type", "video")
	v.Add("maxResults", "50")
	v.Add("order", "date")
	v.Add("channelId", c.ChannelID)

	if len(nextPageToken) != 0 {
		v.Set("pageToken", nextPageToken)
	}

	u.RawQuery = v.Encode()
	return u.String()
}

func (c *ChannelClient) searchList(nextPageToken string) (*SearchListResponse, error) {
	searchURL := c.SearchURL(nextPageToken)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("youtube error response: %v", string(body))
	}

	var searchList SearchListResponse

	if err := json.NewDecoder(resp.Body).Decode(&searchList); err != nil {
		return nil, err
	}

	return &searchList, nil
}

func (c *ChannelClient) videoList(search *SearchListResponse) (*VideoListResponse, error) {

	videoIDs := search.VideoIDs()
	videosURL := c.videosURL(videoIDs)

	var videoList = VideoListResponse{}

	req, err := http.NewRequest("GET", videosURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("youtube error response: %v", string(body))
	}

	err = json.NewDecoder(resp.Body).Decode(&videoList)
	if err != nil {
		return nil, err
	}

	return &videoList, nil
}

func (c *ChannelClient) videosURL(videoIDs []string) string {

	u, _ := url.Parse(fmt.Sprintf("%svideos", c.BaseURL))

	v := url.Values{}
	v.Set("key", c.APIKey)
	v.Add("part", "snippet,statistics,status,contentDetails")
	v.Add("id", strings.Join(videoIDs, ","))

	u.RawQuery = v.Encode()
	return u.String()
}

func (searchList *SearchListResponse) VideoIDs() []string {
	// items.collect{|item| item.video_id}.compact.uniq.join(",")
	keys := make(map[string]bool)
	ids := []string{}
	for _, item := range searchList.Items {
		id := item.ID.VideoID
		if _, value := keys[id]; !value {
			keys[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}
