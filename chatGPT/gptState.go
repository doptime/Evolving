package chatGPT

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/yangkequn/saavuu/data"
)

var tokenLoad, XyOverload = false, false

func gptReportOverLoad(err error) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "no AccessToken") {
		tokenLoad = true
		go func() {
			time.Sleep(time.Duration(120) * time.Second)
			tokenLoad = false
		}()
	}
	if strings.Contains(err.Error(), "请求量过大") {
		XyOverload = true
		go func() {
			time.Sleep(30 * time.Second)
			XyOverload = false
		}()
	}
}
func gptOverLoadDelay(forceDelaySececonds int) {
	if tokenLoad {
		//random sleep between 3min to 5min
		time.Sleep(time.Duration(rand.Intn(120)+180) * time.Second)
	} else if XyOverload {
		//random sleep between 20s to 80s
		time.Sleep(time.Duration(rand.Intn(60)+20) * time.Second)
	} else if forceDelaySececonds > 0 {
		time.Sleep(time.Duration(forceDelaySececonds) * time.Second)
	}
}

type AuthToken struct {
	ConsecutiveFailNum int
	AccountInfo        string
	AccessToken        string
	AvailableHench4    int64
	AvailableHench35   int64
	RefreshToekn       string
	RefreshAt          int64
}

var keyAuthTokens = data.NewStruct[string, *AuthToken]().WithRedis("Auth")

func NewAuthToken(AccountInfo, RefreshToekn string) (authToken *AuthToken) {
	authToken = &AuthToken{AccountInfo: AccountInfo, RefreshToekn: RefreshToekn}
	keyAuthTokens.HSet(authToken.AccountInfo, authToken)
	return authToken
}

func refreshCookieToAccessToken(authToken *AuthToken) (ok bool, err error) {
	var (
		refreshCookie        string = authToken.RefreshToekn
		accessToken, messege string
	)
	//messege, err := FakeAIPost("https://yangkequn.xyhelper.cn/getsession", "", "xyhelper", map[string]string{"refreshCookie": refreshCookie})
	messege, err = FakeAIPost("https://demo.xyhelper.cn/getsession", "", "xyhelper", map[string]string{"refreshCookie": refreshCookie})
	if err != nil {
		return false, err
	}
	js, ok := gjson.Parse(messege).Value().(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("update account's js error, %s", messege)
	}
	if accessToken, ok = js["accessToken"].(string); !ok || len(accessToken) == 0 {
		return false, fmt.Errorf("update account's accessToekn fail, %s", messege)
	}
	if _refreshCookie, ok := js["refreshCookie"].(string); ok && len(_refreshCookie) > 0 {
		authToken.RefreshToekn = _refreshCookie
		authToken.AccessToken = accessToken
		authToken.RefreshAt = time.Now().Unix()
		//update refreshCookie
		if err = keyAuthTokens.HSet(authToken.AccountInfo, authToken); err != nil {
			return false, err
		}
	}
	return true, nil
}

func UpdateAccountContinuousely() {
	var (
		err                                 error
		NewAuthTokenFromDBCnt               int = 0
		AuthTokenMapLocal, AuthTokenMapInDB map[string]*AuthToken
	)
	//load last Data or save to db
	AuthTokenMapLocal = lo.SliceToMap(AuthTokenSlice, func(v *AuthToken) (key string, value *AuthToken) { return v.AccountInfo, v })
	//not FirstRun, save to db
	if len(AuthTokenMapLocal) > 0 {
		//save to db
		keyAuthTokens.HSet(AuthTokenMapLocal)
		log.Info().Msgf("save %v AuthToken to db", len(AuthTokenMapLocal))
	}
	//download new accessToken from db
	if AuthTokenMapInDB, err = keyAuthTokens.HGetAll(); err == nil && len(AuthTokenMapInDB) > 0 {
		for _, authToken := range AuthTokenMapInDB {
			if existedTokens := lo.Filter(AuthTokenSlice, func(v *AuthToken, index int) bool { return v.AccountInfo == authToken.AccountInfo }); len(existedTokens) == 0 {
				AuthTokenSlice = append(AuthTokenSlice, authToken)
				NewAuthTokenFromDBCnt++
			}
		}
		if NewAuthTokenFromDBCnt > 0 {
			log.Info().Msgf("New authToken from db: +%v ", NewAuthTokenFromDBCnt)
		}
	}
	//remove  AuthToken theat nologer exists in db
	for _, authToken := range AuthTokenSlice {
		if _, ok := AuthTokenMapInDB[authToken.AccountInfo]; !ok {
			AuthTokenSlice = lo.Filter(AuthTokenSlice, func(v *AuthToken, index int) bool { return v.AccountInfo != authToken.AccountInfo })
			log.Info().Msgf("remove invalid AuthToken: %v", authToken.AccountInfo)
		}
	}

	//refresh refresh Cookie every day
	for _, authToken := range AuthTokenSlice {
		if authToken.RefreshAt < time.Now().Unix()-86400 || len(authToken.AccessToken) == 0 {
			refreshCookieToAccessToken(authToken)
		}
	}
	//retry every 10 minutes
	go func() {
		time.Sleep(time.Minute * 10)
		go UpdateAccountContinuousely()
	}()
}

// {0, "19339", "", "", ""},
// {0, "19339", "", "", ""},
// {0, "19339", "", "", ""},
var AuthTokenSlice = []*AuthToken{}

type AuthTokenQuota struct {
	Cursor35 uint
	Cursor4  uint
}

var AuthQuata *AuthTokenQuota = &AuthTokenQuota{
	Cursor35: uint(rand.Float32() * float32(len(AuthTokenSlice))),
	Cursor4:  uint(rand.Float32() * float32(len(AuthTokenSlice)))}

func init() {
	log.Info().Msg("init apiGptCompletion.go, start UpdateAccountContinuousely ")
	if os.Getenv("UpdateRefreshtoken") == "true" {
		UpdateAccountContinuousely()
	}
}

func QuotaCheckOK(Model string) (AccessToken string) {
	var token *AuthToken
	var now = time.Now().Unix()
	for i := 0; i < len(AuthTokenSlice); i++ {
		//case 1: GPT4
		if strings.Contains(Model, "4") {
			AuthQuata.Cursor4 = (AuthQuata.Cursor4 + 1) % uint(len(AuthTokenSlice))
			token = AuthTokenSlice[int(AuthQuata.Cursor4)]
			//prevent "deactivated"
			if now > token.AvailableHench4 && len(token.AccessToken) > 60 {
				return token.AccessToken
			}
		} else { //case 2: GPT3.5
			AuthQuata.Cursor35 = (AuthQuata.Cursor35 + 1) % uint(len(AuthTokenSlice))
			token = AuthTokenSlice[int(AuthQuata.Cursor35)]
			//prevent "deactivated"
			if now > token.AvailableHench35 && len(token.AccessToken) > 60 {
				return token.AccessToken
			}
		}
	}
	return ""
}
func ReportError(AccessToekn string, Model string, err error) {
	if len(AccessToekn) == 0 {
		return
	}
	for _, v := range AuthTokenSlice {
		if v.AccessToken == AccessToekn {
			if err == nil {
				v.ConsecutiveFailNum = 0

			} else if err != nil && strings.Contains(err.Error(), "detail") && (strings.Contains(err.Error(), "limit") || strings.Contains(err.Error(), "model_cap_exceeded") || strings.Contains(err.Error(), "deactivated")) {
				//quota limit reached, which means the account is valid
				v.ConsecutiveFailNum = 0
				if strings.Contains(Model, "4") {
					v.AvailableHench4 = time.Now().Unix() + 3600*4
				} else {
					v.AvailableHench35 = time.Now().Unix() + 3600*4
				}
				if strings.Contains((err.Error()), "deactivated") {
					v.AvailableHench4 = time.Now().Unix() + 100*86400 + (rand.Int63() % 100)
					v.AvailableHench35 = time.Now().Unix() + 100*86400 + (rand.Int63() % 100)
					v.AccessToken = "deactivated"
				}
			} else {
				v.ConsecutiveFailNum++
			}
			return
		}
	}
}
