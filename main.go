package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/michimani/gotwi"
)

type TweetData struct {
	TweetId         string
	TweetCount      int
	ZeroTweetStreak int
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

const (
	testMode bool = false
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("GOTWI_ACCESS_TOKEN")
	apiKeySecret := os.Getenv("GOTWI_ACCESS_TOKEN_SECRET")
	accessToken := os.Getenv("GOTWI_ACCESS_TOKEN")
	//accessTokenSecret := os.Getenv("GOTWI_ACCESS_TOKEN_SECRET")
	userID := os.Getenv("USER_ID")
	if userID == "" {
		panic("userIDがありません!")
	}

	tokens := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           apiKey,
		OAuthTokenSecret:     apiKeySecret,
	}
	client, err := gotwi.NewClient(tokens)
	client.SetAccessToken(accessToken)
	if err != nil {
		panic(err)
	}

	//client.SetAccessToken(accessToken)

	isExistYesterdayData := true
	d, err := GetYesterdaysData("data.json")
	if err != nil {
		isExistYesterdayData = false
		fmt.Printf("昨日のデータを取得できませんでした: %v\n", err)
	}

	var todaysData TweetData

	// get *time.Time of yesterday
	jst, _ := time.LoadLocation("Asia/Tokyo")
	yesterday := GetDateRange(-1, jst)
	twoDaysAgo := GetDateRange(-2, jst)

	//search and count yesterday tweets
	todaysData.TweetCount, err = CountTweets(client, userID, &yesterday.Start, &yesterday.End)
	if err != nil {
		panic(err)
	}

	pastTweetCount, err := CountTweets(client, userID, &twoDaysAgo.Start, &twoDaysAgo.End)
	if err != nil {
		panic(err)
	}

	letterCount, err := CountTweetsLetter(client, userID, &yesterday.Start, &yesterday.End)
	if err != nil {
		panic(err)
	}

	dayBeforeMsg := ""
	if isExistYesterdayData {
		todaysData.TweetCount -= 1
		if todaysData.TweetCount <= 0 {
			if !testMode {
				if d.TweetId != "" {
					chk, err := DeleteTweet(client, d.TweetId)
					if err != nil {
						panic(err)
					}
					if chk {
						fmt.Println("昨日のBotツイートを削除しました")
					}
				}
			}

			todaysData.ZeroTweetStreak = d.ZeroTweetStreak + 1
			dayBeforeMsg = fmt.Sprintf("(%d日連続)", todaysData.ZeroTweetStreak)
		} else {
			dayBeforeMsg = fmt.Sprintf("(おととい:%d)", pastTweetCount-1)
		}
	}

	tweetCountMsg := fmt.Sprintf("昨日のツイート数:%d%v\n", todaysData.TweetCount, dayBeforeMsg)
	letterCountMsg := ""
	if todaysData.TweetCount != 0 {
		letterCountMsg = fmt.Sprintf("ツイートの文字数の合計:%d(1ツイートあたり%d文字)\n", letterCount, letterCount/todaysData.TweetCount)
	}

	// calculate the days before entrance exams
	secondTest := time.Date(2023, time.March, 12, 0, 0, 0, 0, jst)
	subS := int((secondTest.Sub(yesterday.Start) - 1).Hours() / 24)

	message := "【Botによる定期ツイート】\n"
	message += "こんばんは。"
	message += fmt.Sprintf("%vになりました\n", time.Now().Format("2006年1月2日"))
	message += tweetCountMsg
	message += letterCountMsg
	message += "\n"
	message += fmt.Sprintf("国立後期試験まで%d日です", subS)

	fmt.Println(message)
	if !testMode {
		id, err := PostTweet(client, message)
		if err != nil {
			panic(err)
		}
		if id != "" {
			todaysData.TweetId = id
			fmt.Printf("ツイートしました: %v\n", id)
		}
	}

	SaveTodaysData("data.json", todaysData)

}

func GetDateRange(daysAfter int, timezone *time.Location) DateRange {
	now := time.Now()
	y := now.AddDate(0, 0, daysAfter)

	start := time.Date(y.Year(), y.Month(), y.Day(), 0, 0, 0, 0, timezone)
	end := time.Date(y.Year(), y.Month(), y.Day(), 23, 59, 59, 0, timezone)

	return DateRange{Start: start, End: end}

}

func GetYesterdaysData(path string) (TweetData, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return TweetData{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return TweetData{}, err
	}

	var jsonData TweetData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return TweetData{}, err
	}

	fmt.Println(jsonData)
	return jsonData, nil
}

func SaveTodaysData(path string, src TweetData) error {
	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, jsonData, 0666); err != nil {
		return err
	}

	return nil
}
