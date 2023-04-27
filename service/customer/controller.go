package customer

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"weichatapp/public"
)

type GetCustomerInfoRequest struct {
	Code string `json:"code"`
}

type CustomerController struct {
	CustomerCache *CustomerCache
}

func (pc *CustomerController)getCustomerInfo(c *gin.Context){
	log.Println("CustomerController getCustomerInfo start")
	var req GetCustomerInfoRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusOK, nil)
		return 
	}

	/*userInfo1:=&public.UserInfoResponse{
		OpenID:"oYQq-5Z4Q1Q1Q1Q1Q1Q1Q1Q1Q1Q1",
		NickName:"测试用户",
		HeadImgURL:"https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.duitang.com%2Fuploads%2Fblog%2F202106%2F09%2F20210609081952_51ef5.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1685149221&t=741db3d3833945213fe8db2dd66b88b4",
	}

	c.IndentedJSON(http.StatusOK, userInfo1)
	return */

	log.Println("auth code:",req.Code)
	accessToken:=public.GetUserAccessTokenByCode(req.Code)
	if accessToken==nil {
		log.Println("GetUserAccessTokenByCode failed")
		c.IndentedJSON(http.StatusOK, nil)
		return 
	}

	customerInof:=&CustomerInfo{
		OpenID:accessToken.OpenID,
		AccessToken:accessToken.AccessToken,
		RefreshToken:accessToken.RefreshToken,
		ExpiresIn:accessToken.ExpiresIn,
		Scope:accessToken.Scope,
	}

	pc.CustomerCache.SaveCustomerInfo(customerInof)

	userInfo:=public.RefreshCustomerInfo(accessToken.OpenID,accessToken.AccessToken)
	c.IndentedJSON(http.StatusOK, userInfo)
	log.Println("CustomerController getCustomerInfo end")
}

func (opc *CustomerController) Bind(router *gin.Engine) {
	router.POST("/customer/getCustomerInfo", opc.getCustomerInfo)
}
