package feishu

// @Title  飞书发送消息模块
// @Description  飞书发送消息模块
// @Author  张振华  2020-10-31
// @Update  张振华  2020-10-31

import (
	"fmt"
	"io/ioutil"
	"github.com/garudaen/beluga/feishu/core/body"
	"github.com/garudaen/beluga/feishu/core/json"
	"github.com/garudaen/beluga/feishu/core/urls"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/widuu/gojson"
)

//Auth 授权
type Auth struct {
	AppID             string `json:"app_id"`
	UserID            string `json:"user_id"`
	AppSecret         string `json:"app_secret"`
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
	Authorization     string `json:"authorization"`
}

//GetAccessToken 获取 access token
func (auth *Auth) GetAccessToken() *body.AppAccessTokenInternalRespBody {
	postBody := json.ToJSONIgnoreError(auth)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urls.APIAppAccessTokenInternal, strings.NewReader(postBody))
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	res := &body.AppAccessTokenInternalRespBody{}
	json.FromJSONIgnoreError(string(respBody), res)
	//_ = json.Unmarshal(respBody, &res)
	auth.AppAccessToken = res.AppAccessToken
	auth.TenantAccessToken = res.TenantAccessToken
	return res
}

//GetToken 获取 access token
func (auth *Auth) getToken() {
	auth.AppID = viper.GetString("feishu.AppID")
	auth.AppSecret = viper.GetString("feishu.AppSecret")
	postBody := json.ToJSONIgnoreError(auth)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urls.APIAppAccessTokenInternal, strings.NewReader(postBody))
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	res := &body.AppAccessTokenInternalRespBody{}
	json.FromJSONIgnoreError(string(respBody), res)
	//_ = json.Unmarshal(respBody, &res)
	auth.AppAccessToken = res.AppAccessToken
	auth.TenantAccessToken = res.TenantAccessToken
	auth.Authorization = fmt.Sprintf("Bearer %s", auth.TenantAccessToken)
}

//GetUserID 通过手机号获取用户 ID
func (auth *Auth) GetUserID(phone string) *body.UserInfo {
	auth.getToken()
	reqURL := fmt.Sprintf("%s?mobiles=%s", urls.APIGetUserFromPhone, phone)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth.Authorization)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	res := gojson.Json(string(respBody))
	userID := res.Get("data").Get("mobile_users").Get(phone).Getindex(1).Get("user_id").Tostring()
	defer resp.Body.Close()
	cb := &body.UserInfo{}
	if userID == "" {
		cb.Code = 1
		cb.Msg = "手机号码有误"
		return cb
	}
	auth.UserID = userID
	cb.UserID = userID
	cb.Code = 0
	cb.Msg = "ok"
	return cb
}

//SendTextMessageToUser 发送消息给用户
func (auth *Auth) SendTextMessageToUser(phone string, content string) *body.SendMessageRespBody {
	info := auth.GetUserID(phone)
	sendMessageToUserRespBody := &body.SendMessageRespBody{}
	if info.Code != 0 {
		sendMessageToUserRespBody.Code = info.Code
		sendMessageToUserRespBody.Msg = fmt.Sprintf("消息发送失败，原因：%s", info.Msg)
		return sendMessageToUserRespBody
	}
	//client := &http.Client{}
	postBody := &body.SendTextMessageToUserPostBody{}
	postBody.UserID = info.UserID
	postBody.MsgType = "text"
	postBody.Content.Text = content
	postbody := json.ToJSONIgnoreError(postBody)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urls.APIRobotSendMessage, strings.NewReader(postbody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth.Authorization)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.FromJSONIgnoreError(string(respBody), sendMessageToUserRespBody)
	defer resp.Body.Close()
	return sendMessageToUserRespBody
}

// SendTextMessageToChatGroup 发送群消息
// @param chatID string 会话（包括单聊、群聊）的唯一标识。通过创建群接口返回
// @param content string 消息内容文本
func (auth *Auth) SendTextMessageToChatGroup(chatID string, content string) *body.SendMessageRespBody {
	auth.getToken()
	sendMessageToUserRespBody := &body.SendMessageRespBody{}
	postBody := &body.SendTextMessageToChatGroupPostBody{}
	postBody.ChatID = chatID
	postBody.MsgType = "text"
	postBody.Content.Text = content
	postbody := json.ToJSONIgnoreError(postBody)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urls.APIRobotSendMessage, strings.NewReader(postbody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth.Authorization)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.FromJSONIgnoreError(string(respBody), sendMessageToUserRespBody)
	defer resp.Body.Close()
	return sendMessageToUserRespBody
}

// CreateChatGroup 创建飞书聊天群
// params users 用户列表
// params groupName 群名称
func (auth *Auth) CreateChatGroup(users []string, groupName string) *body.CreateChatGroupRespBody {
	auth.getToken()
	createChatGroupRespBody := &body.CreateChatGroupRespBody{}
	postBody := body.CreateChatGroupPostBodyDefault()
	postBody.UserIDs = users
	postBody.GroupName = groupName
	postbody := json.ToJSONIgnoreError(postBody)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urls.APICreateChatGroup, strings.NewReader(postbody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth.Authorization)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	json.FromJSONIgnoreError(string(respBody), createChatGroupRespBody)
	defer resp.Body.Close()
	return createChatGroupRespBody
}
