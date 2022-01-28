package controllers

import (
  "fmt"
  "time"
  "io/ioutil"
  "encoding/json"

  "github.com/ECHibiki/Community-Banners-2.0/bannerjwt"
  "github.com/gin-gonic/gin"
)
// the functions in this file link to the routes in ginterface.go
// functions within are outbound to the files in the comments above the functions
// a few misc shared functions at the bottom of this file

type ControllerSettings struct{
  RedirectDomain string
  // minutes
  AccountInterval int64
  AttemptInterval int64
  MaxAttempts int
}

var controller_settings ControllerSettings

func Init(){
  setting_json_bytes, err := ioutil.ReadFile("./settings/controller-settings.json")
  if err != nil{
    panic(err)
  }
  json.Unmarshal(setting_json_bytes, &controller_settings)
}

//File: PageGenerationController
func GenerateAdPage(c *gin.Context){
  size := getParam(c ,"size")
  name := getGet(c, "name")
  // https://github.com/gin-gonic/gin/issues/2697
  // If it is a trusted IP (which means the request is redirected by a proxy),
  // then it will try to parse the "real user IP" from X-Forwarded-For/X-Real-Ip header.
  ip := c.ClientIP()
  page := returnAdPage(name , size, ip)

  c.HTML(200 , "banner.html" , page)
}
//File: PageGenerationController
func GenerateAdJSON(c *gin.Context){
  size := getParam(c ,"size")
  name := getGet(c, "name")
  ip := c.ClientIP()
  ad_gin_h := returnAdJson(name , size, ip)
  c.JSON(200, ad_gin_h)
}
//File: PageGenerationController
// return info limited to entries not created by shadow banned users
func GetLimitedInfo(c *gin.Context){
  name := getGet(c, "name")
  ip := c.ClientIP()
  ad_data := getLimitedEntries(name , ip);
  var ad_gin_h []gin.H
  for _ , ad := range ad_data{
    gh := gin.H{
      "url": ad["url"],
      "uri": ad["uri"],
      "name":ad["ads.fk_name"],
      "size":ad["size"],
      "clicks":ad["clicks"],
    }
    ad_gin_h = append(ad_gin_h , gh)
  }
  c.JSON(200, ad_gin_h);
}

//File: RedirectController
func RedirectSiteRequest(c *gin.Context){

  site := getParam(c , "s")
  file := getParam(c , "f")

  if site == ""{
      c.HTML(404, "url-error.html", "Non-existing URL")
      return
  }

  if !checkURIExists(file){
    c.HTML(404, "uri-error.html", "Non-existing URI")
    return
  }

  incrementClicksSQL(file);

  c.Redirect(301 , site)
}

//File: UserCreationController
func CreateNewUser(c *gin.Context){
  ip := c.ClientIP()
  name := getParam(c, "name")
  if len(name) > 30 {
    c.JSON( 401 , gin.H{"error": "Name should not be longer than 30 characters"});
    return
  }
  if name == ""{
    c.JSON( 401 , gin.H{"error": "Insert a name"});
    return
  }
  pass := getParam(c, "pass")
  pass_confirm := getParam(c, "pass_conf")
  if pass != pass_confirm{
    c.JSON( 401 , gin.H{"error": "Passwords do not match"});
    return
  }
  if pass == ""{
    c.JSON( 401 , gin.H{"error": "Insert a password"});
    return
  }

  if !validateIPCreation(ip){
    c.JSON( 401 , gin.H{"error": "Too many accounts"});
    return
  }
  updateIPCreation(ip)


  error_message := addNewUserToDB(name, pass);
  if error_message == ""{
    c.JSON( 200 , gin.H{"log": "Successfully Created"});
  } else{
    c.JSON( 401 , gin.H{"error": error_message});
  }
}
//File: UserCreationController
func RejectUserCreation(c *gin.Context){
  c.JSON(200 , gin.H{"warn": "Pool Closed - Come back later"})
}

//File: UserSignInController
func LoginUser(c *gin.Context){
  name := getParam(c, "name")
  if name == ""{
    c.JSON( 401 , gin.H{"error": "Insert a name"});
    return
  }
  pass := getParam(c, "pass")
  if pass == ""{
    c.JSON( 401 , gin.H{"error": "Insert a password"});
    return
  }
  if !checkHardBanned(name) {
    c.JSON( 401 , gin.H{"error": "You've been banned..."});
    return
  }
  // N attempts every X for given IP
  ip := c.ClientIP()
  if !validateNameBruteForce(ip) {
    c.JSON( 401 , gin.H{"error": "Too many password attempts in timespan"});
    return
  }
  if !checkAuthentication(name, pass) {
    c.JSON( 401 , gin.H{"error": "Invalid username/password"});
    return
  } else{
    updateNameBruteForce(ip)
  }

  token , err := bannerjwt.CreateToken(name)
  if err != nil{
    panic (err)
  }
  c.JSON(200 , gin.H{
    "access_token" : token,
    "token_type" : "bearer",
    "expires_in" : time.Now().Add(time.Hour * 24).Unix(),
    "log" : "Successfully Logged In",
  })
}

//File: ConfidentialInfoController
func AccessInfo(c *gin.Context){

}
//File: ConfidentialInfoController
func CreateInfo(c *gin.Context){

}
//File: ConfidentialInfoController
func RemoveInfo(c *gin.Context){

}

//File: ModeratorActivityController
func GetAllInfo(c *gin.Context){

}
//File: ModeratorActivityController
func BanUser(c *gin.Context){

}
//File: ModeratorActivityController
func DeleteAll(c *gin.Context){

}
//File: ModeratorActivityController
func DeleteIndividual(c *gin.Context){

}

/* MISC FUNCTIONS */

func getParam(c *gin.Context, key string) string{
  q := c.Request.URL.Query()
  var value string
  if key_arr, key_exists := q["f"] ; key_exists && len(key_arr) != 0{
    value = key_arr[0]
  } else{
    value = ""
  }
  return value
}

func getGet(c *gin.Context , key string) string{
  context_value , c_exists := c.Get(key)
  var value string
  if !c_exists{
    value = ""
  } else{
    value = fmt.Sprintf("%v", context_value)
  }
  return value
}