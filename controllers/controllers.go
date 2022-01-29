package controllers

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "io/ioutil"
  "encoding/json"

  "github.com/ECHibiki/Community-Banners-2.0/bannerjwt"
  "github.com/ECHibiki/Community-Banners-2.0/mailer"
  "github.com/gin-gonic/gin"
)
// the functions in this file link to the routes in ginterface.go
// functions within are outbound to the files in the comments above the functions
// a few misc shared functions at the bottom of this file

type ControllerSettings struct{
  RedirectDomain string
  // minutes
  AccountInterval int64
  BannerInterval int64
  AttemptInterval int64
  MaxAttempts int
  PublicPath string
  MaxFileSize int64
  SmallDimensionsX int
  SmallDimensionsY int
  WideDimensionsX int
  WideDimensionsY int
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
  name := c.PostForm("name")
  if len(name) > 30 {
    c.JSON( 401 , gin.H{"error": "Name should not be longer than 30 characters"});
    return
  }
  if name == ""{
    c.JSON( 401 , gin.H{"error": "Insert a name"});
    return
  }
  pass := c.PostForm("pass")
  pass_confirm := c.PostForm("pass_conf")
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
  c.JSON(200 , gin.H{"error": "Pool Closed - Come back later"})
}

//File: UserSignInController
func LoginUser(c *gin.Context){
  name := c.PostForm("name")
  if name == ""{
    c.JSON( 401 , gin.H{"error": "Insert a name"});
    return
  }
  pass := c.PostForm("pass")
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
  is_mod := checkIsMod(name)
  token , err := bannerjwt.CreateToken(name , is_mod)
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
  name := getParam(c , "name")
  is_mod := getParam(c , "is_mod")
  ad_arr := getUserData(name);

  c.JSON(200 , gin.H{
    "name" : name,
    "mod" : is_mod,
    "ads" : ad_arr,
  })
}
//File: ConfidentialInfoController
func CreateBanner(c *gin.Context){
  ip := c.ClientIP()
  cooldown_timer := doCreationCooldown(ip)
  if cooldown_timer <= 0{
    c.JSON(401 , gin.H{"error": "Posting too fast(" + strconv.Itoa(cooldown_timer) + ") seconds"})
    return
  }
  updateCreationCooldown(ip)

  image , header , _:= c.Request.FormFile("image")
  fsize := header.Size
  if fsize > controller_settings.MaxFileSize {
    c.JSON(401 , gin.H{"error": "Filesize is larger than " + strconv.Itoa(int(controller_settings.MaxFileSize / (1000 * 1000))) + " MB"})
  }

  size := c.PostForm("size")
  url := c.PostForm("url")
  name := getParam(c , "name")

  if size == "wide" && url == ""{
    c.JSON(401 , gin.H{"error": "URL Required" })
  }

  x , y := getImageDimensions(image)
  if size == "small" &&
    x == controller_settings.SmallDimensionsX &&
    y == controller_settings.SmallDimensionsY {
  } else if size == "wide" &&
    x == controller_settings.WideDimensionsX &&
    y == controller_settings.WideDimensionsY {
  } else{
    c.JSON(401 , gin.H{"error": "Dimensions are incorrect" })
    return
  }

  tmp_location := controller_settings.PublicPath + "storage/tmp/" + "tmp-" + ip + "-name"
  handle, err := os.Create(tmp_location)
  if err != nil {
    panic(err)
  }
  handle.Close();

  unique := checkUniqueBanner(tmp_location)
  if !unique{
    os.Remove(tmp_location)
    c.JSON(401 , gin.H{"error": "This banner is not unique"})
    return
  }

  file_base64 := getBase64(tmp_location)
  filename := generateHashedFileName(tmp_location)
  file_path := "storage/image/" + filename
  err = os.Rename(tmp_location , controller_settings.PublicPath + file_path)
  if err != nil{
    panic(err)
  }
  addBanner(name , file_path , url, ip ,  size  )
  response := mailer.SendBannerEmail( name , file_base64 , url , file_path )
  if response != "sent"{
    c.JSON(200, gin.H{"warn" : "Banner Added "})
  } else{
    c.JSON(200, gin.H{"log" : "Banner Added "})
  }
}

//File: ConfidentialInfoController
func RemoveBanner(c *gin.Context){
  name := getParam(c , "name")
  uri := getGet(c , "uri")
  if !affirmImageIsOwned(name , uri) {
    c.JSON(401 , gin.H{
      "error" : "This banner is not owned",
    })
    return
  } else{
    removeAdSQL(name , uri)
    removeAdImage(uri)
    c.JSON(200 , gin.H{
      "log" : "Banner removed",
    })
    return
  }
}

//File: ModeratorActivityController
func GetAllBanners(c *gin.Context){
  entires := getAllEntries()
  c.JSON(200 , entires)
}
//File: ModeratorActivityController
func BanUser(c *gin.Context){
  target := getGet(c , "target")
  hard := getGet(c , "hard")
  if target == "" || hard == "" {
    c.JSON(401 , gin.H{"error": "Fields missing"})
    return
  }
  hard_ban , err := strconv.Atoi(hard)
  if err != nil{
    panic(err)
  }
  createNewBan(target, hard_ban)

  hard_ban_str := ""
  if hard_ban == 1 {
    hard_ban_str = "hard"
  } else{
    hard_ban_str = "soft"
  }
  c.JSON(200 , gin.H{"log": "User " + target + " was " + hard_ban_str + " banned"})
}
//File: ModeratorActivityController
func DeleteAll(c *gin.Context){
  target := getGet(c , "target")
  if target == "" {
    c.JSON(401 , gin.H{"error": "Fields missing"})
    return
  }
  removeAllUserImages(target)
  removeUserFromDatabase(target)

  c.JSON(200 , gin.H{"log": "User " + target + " was purged"})
}
//File: ModeratorActivityController
func DeleteIndividual(c *gin.Context){
  name := getGet(c , "name")
  uri := getGet(c , "uri")
  if name == "" || uri == ""{
    c.JSON(401 , gin.H{"error": "Fields missing"})
    return
  }
  removeIndividualBannerFromImages(uri)
  removeIndividualBannerFromDB(uri)

  c.JSON(200 , gin.H{"log": name + "'s image was pruned"})
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