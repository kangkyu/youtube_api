package youtube_api

import (
	"errors"
)

type PaginatedFetcher struct {
	client        *ChannelClient
	nextPageToken string
	hasMorePages  bool
}

func (c *ChannelClient) NewPaginatedFetcher() *PaginatedFetcher {
	return &PaginatedFetcher{
		client:       c,
		hasMorePages: true,
	}
}

func (f *PaginatedFetcher) FetchNextPage() (*VideoListResponse, error) {
	if !f.hasMorePages {
		return nil, errors.New("no more pages to fetch")
	}

	response, err := f.client.searchList(f.nextPageToken)
	if err != nil {
		return nil, err
	}

	f.nextPageToken = response.NextPageToken

	if len(f.nextPageToken) == 0 || len(response.Items) == 0 {
		f.hasMorePages = false
	}

	return f.client.videoList(response)
}

func (f *PaginatedFetcher) HasNextPage() bool {
	return f.hasMorePages
}
