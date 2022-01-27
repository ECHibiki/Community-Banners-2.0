package controllers

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

//PageGenerationController
func GenerateAdPage(c *gin.Context){
  q := c.Request.URL.Query()
  var size string
  if size_arr, exists := q["size"] ; exists && len(size_arr) != 0{
    size = size_arr[0]
  } else{
    size = ""
  }
  name := fmt.Sprintf("%v", c.MustGet("name"))
  // https://github.com/gin-gonic/gin/issues/2697
  // If it is a trusted IP (which means the request is redirected by a proxy),
  // then it will try to parse the "real user IP" from X-Forwarded-For/X-Real-Ip header.
  ip := c.ClientIP()
  page := returnAdPage(name , size, ip)

  c.HTML(200 , "banner.html" , page)
}
//PageGenerationController
func GetLimitedInfo(c *gin.Context){

}
//PageGenerationController
func GenerateAdJSON(c *gin.Context){

}

//RedirectController
func RedirectSiteRequest(c *gin.Context){

}

//UserCreationController
func CreateNewUser(c *gin.Context){

}

//UserSignInController
func LoginUser(c *gin.Context){

}
//UserSignInController
func RejectUserCreation(c *gin.Context){

}
//UserSignInController
func LoginMod(c *gin.Context){

}

//ConfidentialInfoController
func AccessInfo(c *gin.Context){

}
//ConfidentialInfoController
func CreateInfo(c *gin.Context){

}
//ConfidentialInfoController
func RemoveInfo(c *gin.Context){

}

//ModeratorActivityController
func GetAllInfo(c *gin.Context){

}
//ModeratorActivityController
func BanUser(c *gin.Context){

}
//ModeratorActivityController
func DeleteAll(c *gin.Context){

}
//ModeratorActivityController
func DeleteIndividual(c *gin.Context){

}