package youtube_api

// generate_csv/generate_csv.go:5:2: "encoding/json" imported and not used
// generate_csv/generate_csv.go:7:2: "io/ioutil" imported and not used
// generate_csv/generate_csv.go:9:2: "net/url" imported and not used
// generate_csv/generate_csv.go:10:2: "os" imported and not used

import (
    "encoding/json"
    "io/ioutil"
    "net/url"
    "os"
    "strings"
)

func VideoListFromVideosURL(vu *url.URL) (videoListResponse, error) {

    var videoList = videoListResponse{}

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

func SearchListFromSearchURL(su *url.URL) (searchListResponse, error) {

    var searchList = searchListResponse{}
    // fmt.Printf("%s\n", su.String())

    resp, err := http.Get(su.String())
    if err != nil {
        return searchList, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
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
