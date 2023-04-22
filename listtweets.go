package main

import (
	"context"
	"regexp"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweet/timeline"
	"github.com/michimani/gotwi/tweet/timeline/types"
	"github.com/rivo/uniseg"
)

var twitterUrl = regexp.MustCompile("https://t.co/([a-zA-Z0-9]+).")

func CountTweets(c *gotwi.Client, id string, startTime *time.Time, endTime *time.Time) (int, error) {
	res, err := ListTweets(c, id, startTime, endTime)
	if err != nil {
		return 0, err
	}
	return len(res), nil
}

func ListTweets(c *gotwi.Client, id string, startTime *time.Time, endTime *time.Time) ([]resources.Tweet, error) {
	p := &types.ListTweetsInput{
		ID:         id,
		StartTime:  startTime,
		EndTime:    endTime,
		MaxResults: 100,
	}
	res, err := timeline.ListTweets(context.Background(), c, p)
	if err != nil {
		return []resources.Tweet{}, err
	}
	return res.Data, nil
}

func CountTweetsLetter(c *gotwi.Client, id string, startTime *time.Time, endTime *time.Time) (int, error) {
	res, err := ListTweets(c, id, startTime, endTime)
	if err != nil {
		return 0, err
	}

	letterCount := 0
	for _, t := range res {
		tweetStr := gotwi.StringValue(t.Text)
		twitterUrl.ReplaceAllString(tweetStr, "")

		letterCount += uniseg.GraphemeClusterCount(tweetStr)
	}
	return letterCount, nil

}
