package main

import (
	"context"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

func PostTweet(c *gotwi.Client, text string) (string, error) {
	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}

func DeleteTweet(c *gotwi.Client, id string) (bool, error) {
	p := &types.DeleteInput{
		ID: id,
	}

	res, err := managetweet.Delete(context.Background(), c, p)
	if err != nil {
		return false, err
	}

	return gotwi.BoolValue(res.Data.Deleted), nil
}
