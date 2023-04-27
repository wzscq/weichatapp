package customer

import (
	"weichatapp/public"
)

type CustomerInfo struct {
	OpenID string `json:"openid"`
	AccessToken string  `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
}

type CustomerRepo struct {
	CustomerCache *CustomerCache
}

func (c *CustomerRepo)GetCustomerInfo(openid string) *public.UserInfoResponse {
	//获取客户信息
	customerInfo:=c.CustomerCache.GetCustomerInfo(openid)
	if customerInfo==nil {
		return nil
	}
  //如果customerInfo不是nil，则调用微信接口刷新客户的access_token
	//刷新access_token
	if c.RefreshAccessToken(customerInfo) == false {
		return nil
	}

	c.CustomerCache.SaveCustomerInfo(customerInfo)
	//刷新客户信息
	userInfoResponse:=public.RefreshCustomerInfo(customerInfo.OpenID,customerInfo.AccessToken)
	return userInfoResponse
}

func (c *CustomerRepo)RefreshAccessToken(customerInfo *CustomerInfo) bool {
	rsp:=public.RefreshAccessToken(customerInfo.RefreshToken)
	if rsp==nil {
		return false
	}

	//更新客户信息
	customerInfo.AccessToken=rsp.AccessToken
	customerInfo.ExpiresIn=rsp.ExpiresIn
	customerInfo.RefreshToken=rsp.RefreshToken
	customerInfo.Scope=rsp.Scope
	return true
}


