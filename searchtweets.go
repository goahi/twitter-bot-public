package main

import (
	"context"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/searchtweet"
	"github.com/michimani/gotwi/tweet/searchtweet/types"
)

func SearchTweet(c *gotwi.Client, q string, startTime *time.Time, endTime *time.Time) (int, error) {
	p := &types.ListRecentInput{
		Query:     q,
		StartTime: startTime,
		EndTime:   endTime,
	}
	res, err := searchtweet.ListRecent(context.Background(), c, p)
	if err != nil {
		return 0, err
	}
	return len(res.Data), nil
}
