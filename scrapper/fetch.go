package scrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultApiBase = "https://www.binance.com/bapi/futures"
)

var defaultHeaders = map[string]string{
	"authority":       "www.binance.com",
	"accept":          "*/*",
	"accept-language": "en-US,en;q=0.8",
	"cache-control":   "no-cache",
	"clienttype":      "web",
	"content-type":    "application/json",
	"lang":            "en",
	"origin":          "https://www.binance.com",
	"pragma":          "no-cache",
	"sec-fetch-dest":  "empty",
	"sec-fetch-mode":  "cors",
	"sec-fetch-site":  "same-origin",
	"sec-gpc":         "1",
	"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
}

type LdbAPIRes[T UserPositionData | UserBaseInfo | []NicknameDetails] struct {
	Success       bool        `json:"success"`
	Code          string      `json:"code"`
	Message       string      `json:"message"`
	Data          T           `json:"data"`
	MessageDetail interface{} `json:"messageDetail"` // ???
}

type UserPositionData struct {
	OtherPositionRetList []rawPosition `json:"otherPositionRetList"`
	UpdateTimeStamp      int64         `json:"updateTimeStamp"`
	UpdateTime           []int         `json:"updateTime"`
}

type rawPosition struct {
	Symbol          string  `json:"symbol"`
	EntryPrice      float64 `json:"entryPrice"`
	MarkPrice       float64 `json:"markPrice"`
	Pnl             float64 `json:"pnl"`
	Roe             float64 `json:"roe"`
	Amount          float64 `json:"amount"`
	UpdateTimeStamp int64   `json:"updateTimeStamp"`
	UpdateTime      []int   `json:"updateTime"`
	Yellow          bool    `json:"yellow"`
	TradeBefore     bool    `json:"tradeBefore"`
	Leverage        int     `json:"leverage"`
}

func GetOtherPosition(ctx context.Context, UUID string, Name string, DelayInterval time.Duration) (LdbAPIRes[UserPositionData], error) {
	return NewUser(UUID, Name, DelayInterval).GetOtherPosition(ctx)
}

func (u *User) GetOtherPosition(ctx context.Context) (LdbAPIRes[UserPositionData], error) {
	var res LdbAPIRes[UserPositionData]
	return res, doPost(ctx, u.client, u.APIBase()+"/v1/public/future/leaderboard", "/getOtherPosition", u.Headers(), strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", u.UID)), &res)
}

type UserBaseInfo struct {
	NickName               string      `json:"nickName"`
	UserPhotoURL           string      `json:"userPhotoUrl"`
	PositionShared         bool        `json:"positionShared"`
	DeliveryPositionShared bool        `json:"deliveryPositionShared"`
	FollowingCount         int         `json:"followingCount"`
	FollowerCount          int         `json:"followerCount"`
	TwitterURL             string      `json:"twitterUrl"`
	Introduction           string      `json:"introduction"`
	TwShared               bool        `json:"twShared"`
	IsTwTrader             bool        `json:"isTwTrader"`
	OpenID                 interface{} `json:"openId"`
}

func GetOtherLeaderboardBaseInfo(ctx context.Context, UUID string, Name string, DelayInterval time.Duration) (LdbAPIRes[UserBaseInfo], error) {
	return NewUser(UUID, Name, DelayInterval).GetOtherLeaderboardBaseInfo(ctx)
}

func (u *User) GetOtherLeaderboardBaseInfo(ctx context.Context) (LdbAPIRes[UserBaseInfo], error) {
	var res LdbAPIRes[UserBaseInfo]
	return res, doPost(ctx, u.client, u.APIBase()+"/v2/public/future/leaderboard", "/getOtherLeaderboardBaseInfo", u.Headers(), strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\"}", u.UID)), &res)
}

type NicknameDetails struct {
	EncryptedUID  string `json:"encryptedUid"`
	Nickname      string `json:"nickname"`
	FollowerCount int    `json:"followerCount"`
	UserPhotoURL  string `json:"userPhotoUrl"`
}

func SearchNickname(ctx context.Context, nickname string) (LdbAPIRes[[]NicknameDetails], error) {
	var res LdbAPIRes[[]NicknameDetails]
	return res, doPost(ctx, http.DefaultClient, defaultApiBase+"/v1/public/future/leaderboard", "/searchNickname", defaultHeaders, strings.NewReader(fmt.Sprintf("{\"nickname\":\"%s\"}", nickname)), &res)
}

func doPost(ctx context.Context, c *http.Client, endpoint, path string, headers map[string]string, data io.Reader, resPtr any) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		endpoint+path,
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return BadStatusError{
			Status:     res.Status,
			StatusCode: res.StatusCode,
			Body:       body,
		}
	}

	return json.Unmarshal(body, resPtr)
}
