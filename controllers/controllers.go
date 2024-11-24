package controllers

import (
  "os"
  "fmt"
  "time"
  "strings"
  "regexp"
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
  // minutes
  FreeMode bool
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
  RestrictedBoards []string
  ValidBoards []string
}

type Login struct{
  Name string `json:"name"`
  Pass string `json:"pass"`
  Pass_Confirm string `json:"pass_confirmation"`
  Token string `json:"token"`
}

type UserRemoval struct{
  URI string `json:"uri"`
}

type ModAction struct {
  Name string `json:"name"`
  Target string `json:"target"`
  URI string `json:"uri"`
  URL string `json:"url"`
  Hard *int `json:"hard"`
}

var controller_settings ControllerSettings

func init(){
  fmt.Println("init controllers...")
  setting_json_bytes, err := ioutil.ReadFile("./settings/controller-settings.json")
  if err != nil{
    panic(err)
  }
  err = json.Unmarshal(setting_json_bytes, &controller_settings)
  if err != nil{
    panic(err)
  }
  fmt.Println("...controllers init")
}

/* File: PageGenerationController */
func GenerateAdPage(c *gin.Context){
  name := getGet(c, "name")
  size := getParam(c ,"size")
  board := getParam(c ,"board")
  // https://github.com/gin-gonic/gin/issues/2697
  // If it is a trusted IP (which means the request is redirected by a proxy),
  // then it will try to parse the "real user IP" from X-Forwarded-For/X-Real-Ip header.
  ip := c.ClientIP()
  page := returnAdPage(name , size, ip , board)

  c.Header("Content-Type", "text/html")
  c.String(200 , page)
}
/* File: PageGenerationController */
func GenerateAdJSON(c *gin.Context){
  name := getGet(c, "name")
  size := getParam(c ,"size")
  board := getParam(c ,"board")
  ip := c.ClientIP()
  ad_gin_h := returnAdJson(name , size, ip , board)
  c.JSON(200, ad_gin_h)
}
/* File: PageGenerationController */
// return info limited to entries not created by shadow banned users
func GetLimitedInfo(c *gin.Context){
  name := getGet(c, "name")
  ip := c.ClientIP()
  ad_data := getLimitedEntries(name , ip);
  ad_gin_h := []gin.H{}
  for _ , ad := range ad_data{
	clicks , _ := strconv.Atoi(ad["clicks"])
    gh := gin.H{
      "url": ad["url"],
      "uri": ad["uri"],
      "name":ad["name"],
      "size":ad["size"],
      "clicks":clicks,
      "board":ad["board"],
    }
    ad_gin_h = append(ad_gin_h , gh)
  }
  c.JSON(200, ad_gin_h);
}

/* File: RedirectController */
func RedirectSiteRequest(c *gin.Context){

  site := getParam(c , "s")
  file := getParam(c , "f")

  if site == ""{
      c.Header("Content-Type", "text/html")
      c.String(404 , "Non-existing URL")
      return
  }

  if !checkURIExists(file){
    c.Header("Content-Type", "text/html")
    c.String(404 , "Non-existing URI")
    return
  }

  incrementClicksSQL(file);

  c.Redirect(301 , site)
}

/* File: UserCreationController */
func CreateNewUser(c *gin.Context){
  ip := c.ClientIP()
  var login Login
  c.BindJSON(&login)

  if len(login.Name) > 30 {
    c.JSON( 400 , gin.H{"error": "Name should not be longer than 30 characters"});
    return
  }
  if login.Name == ""{
    c.JSON( 400 , gin.H{"error": "Insert a name"});
    return
  }
  if len(login.Pass) < 5 {
    c.JSON( 400 , gin.H{"error": "Password is too short"});
    return
  }
  if login.Pass != login.Pass_Confirm{
    c.JSON( 400 , gin.H{"error": "Passwords do not match"});
    return
  }
  if login.Pass == ""{
    c.JSON( 400 , gin.H{"error": "Insert a password"});
    return
  }

  if !validateIPCreation(ip){
    c.JSON( 400 , gin.H{"error": "Too many accounts"});
    return
  }
  updateIPCreation(ip)

  error_message := addNewUserToDB(login.Name, login.Pass);
  if error_message == ""{
    c.JSON( 201 , gin.H{"log": "Successfully Created"});
  } else{
    c.JSON( 400 , gin.H{"error": error_message});
  }
}
/* File: UserCreationController */
func RejectUserCreation(c *gin.Context){
  c.JSON(200 , gin.H{"error": "Pool Closed - Come back later"})
}

/* File: UserSignInController */
// used on both sign in and creation
func LoginUser(domain string) gin.HandlerFunc{
  return func(c *gin.Context){
    var login Login
    c.BindJSON(&login)
    ip := c.ClientIP()
    if login.Name == ""{
      c.JSON( 400 , gin.H{"error": "Insert a name"});
      return
    }
    if login.Pass == ""{
      c.JSON( 400 , gin.H{"error": "Insert a password"});
      return
    }
    if checkHardBanned(login.Name , ip) {
      c.JSON( 403 , gin.H{"error": "You've been banned..."});
      return
    }
    // N attempts every Xsec for given IP
    if !validateNameBruteForce(ip) {
      c.JSON( 401 , gin.H{"error": "Too many password attempts in timespan(" + strconv.Itoa(int(controller_settings.AttemptInterval)) + "min)"});
      return
    }
    if !checkAuthentication(login.Name, login.Pass) {
      c.JSON( 401 , gin.H{"error": "Invalid username/password"});
      return
    } else{
      updateNameBruteForce()
    }
    login.Token  = strings.TrimSpace(login.Token)
    is_donor := false
    if login.Token != ""{
      is_donor = checkIsDonor(login.Token)
      if !is_donor{
        c.JSON( 400 , gin.H{"error": "Your token was not entered correctly"});
        return
      }
    }
    is_mod := checkIsMod(login.Name)
    token , err := bannerjwt.CreateToken(login.Name , is_donor , is_mod)
    if err != nil{
      panic (err)
    }
    c.SetCookie("freeadstoken", token, int(time.Now().Add(time.Hour * 24).Unix()), "/",
      domain, true, true)
    c.JSON(200 , gin.H{
      // "access_token" : gin.H{
      //   "code" : token,
      //   "token_type" : "bearer",
      // },
      // "refresh_token": gin.H {
      //   "code" : "",
      //   "token_type" : "refresh",
      // },
      // "expires_in" : time.Now().Add(time.Hour * 24).Unix(),
      "log" : "Successfully Logged In",
      "donor" : is_donor,
      "mod" : is_mod,
    })
  }
}
/* File: UserSignInController */
// test a token if authenticated
func TestToken(domain string) gin.HandlerFunc{
  return func(c *gin.Context){
    var login Login
    c.BindJSON(&login)
    name := getGet(c , "name")
    is_mod := getGet(c , "is_mod") == "true"

    login.Token  = strings.TrimSpace(login.Token)
    is_donor := checkIsDonor(login.Token)
    if !is_donor{
      c.JSON( 400 , gin.H{"error": "Your token was not entered correctly"});
      return
    }

    token , err := bannerjwt.CreateToken(name , is_donor , is_mod)
    if err != nil{
      panic (err)
    }
    c.SetCookie("freeadstoken", token, int(time.Now().Add(time.Hour * 24).Unix()), "/",
      domain, true, true)
    c.JSON(200 , gin.H{
      "log" : "Successfully verified token",
    })
  }
}

/* File: ConfidentialInfoController */
func AccessInfo(c *gin.Context){
  name := getGet(c , "name")
  is_mod := getGet(c , "is_mod")
  is_donor := getGet(c , "is_donor")
  ad_arr := getUserData(name);
  c.JSON(200 , gin.H{
    "name" : name,
    "mod" : is_mod == "true",
    "donor" : is_donor == "true",
    "ads" : ad_arr,
  })
}
/* File: ConfidentialInfoController */
func CreateBanner(c *gin.Context){
  ip := c.ClientIP()
  is_donor := getGet(c , "is_donor")
  cooldown_timer := doCreationCooldown(ip)
  if cooldown_timer < 0{
    c.JSON(401 , gin.H{"error": "Posting too fast(" + strconv.Itoa(cooldown_timer) + ") seconds"})
    return
  }
  insertHaltCooldown(ip)

  image , header , err:= c.Request.FormFile("image")
  if err != nil {
    c.JSON(400 , gin.H{"error": "File not set or is invalid"})
    return
  }
  defer image.Close()
  mime := getMimeType(image)
  if mime != "image/jpg" &&
    mime != "image/jpeg" &&
    mime != "image/png" &&
    mime != "image/gif"{
      c.JSON(400 , gin.H{"error": "File should be jpeg, png or gif"})
      return
  }


  fsize := header.Size
  if fsize > controller_settings.MaxFileSize {
    c.JSON(400 , gin.H{"error": "Filesize is larger than " + strconv.Itoa(int(controller_settings.MaxFileSize / (1024 * 1024))) + " MB"})
    return
  }

  size := c.PostForm("size")
  url := c.PostForm("url")
  board := c.PostForm("board")
  name := getGet(c, "name")

  if is_donor == "false" && !controller_settings.FreeMode {
    board = ""
  } else if board != ""{
    found := false
    for _ , v_board := range controller_settings.ValidBoards{
      if v_board == board{
        found = true
        break
      }
    }
    if !found {
      c.JSON(400 , gin.H{"error": "The given board is invalid"})
      return
    }
  }

  if size == "wide" && url == ""{
    c.JSON(400 , gin.H{"error": "URL Required" })
    return
  }
  //https://stackoverflow.com/questions/3809401/what-is-a-good-regular-expression-to-match-a-url
  url_regex := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
  if size== "wide" && !url_regex.Match([]byte(url)){
    c.JSON(400 , gin.H{"error": "URL does not look valid" })
    return
  }

  x , y := getImageDimensions(image)
  if size == "small" &&
    x == controller_settings.SmallDimensionsX &&
    y == controller_settings.SmallDimensionsY {
  } else if size == "wide" &&
    x == controller_settings.WideDimensionsX &&
    y == controller_settings.WideDimensionsY {
  } else{
    var dim_str string
    if size == "small" {
      dim_str = "(" +  strconv.Itoa(controller_settings.SmallDimensionsX) + "x" + strconv.Itoa(controller_settings.SmallDimensionsY) + ")"
    } else{
      dim_str = "(" +  strconv.Itoa(controller_settings.WideDimensionsX) + "x" + strconv.Itoa(controller_settings.WideDimensionsY) + ")"
    }
    c.JSON(400 , gin.H{"error": "Dimensions are incorrect" +  dim_str })
    return
  }

  tmp_location := writeTmpFile(name , image)

  uniqueness_response , hash := checkUniqueBanner(tmp_location)
  if uniqueness_response != ""{
    // os.Remove(tmp_location)
    c.JSON(400 , gin.H{"error": uniqueness_response})
    return
  }

  file_base64 := getBase64(tmp_location)
  ext := strings.Replace(mime , "image/", "", -1)
  filename := generateHashedFileName(tmp_location) + "." + ext
  uri := "storage/image/" + filename
  err = os.Rename(tmp_location , controller_settings.PublicPath + uri)
  if err != nil{
    os.Remove(tmp_location)
    os.Remove(controller_settings.PublicPath + uri)
    panic(err)
  }
  addBanner(name , uri , url, ip ,  size , hash , board )
  updateCreationCooldown(ip)
  response := mailer.SendBannerEmail( name , file_base64 , url , uri , header.Filename , board  )
  if response != "Sent"{
    c.JSON(201, gin.H{"warn" : "Banner Added "})
  } else{
    c.JSON(201, gin.H{"log" : "Banner Added "})
  }
}

/* File: ConfidentialInfoController */
func RemoveBanner(c *gin.Context){
  name := getGet(c , "name")
  var removal UserRemoval
  c.BindJSON(&removal)
  if !affirmImageIsOwned(name , removal.URI) {
    c.JSON(403 , gin.H{
      "error" : "This banner is not owned",
    })
    return
  } else{
    removeAdImage(removal.URI)
    removeAdSQL(name , removal.URI)
    c.JSON(200 , gin.H{
      "log" : "Banner removed",
    })
    return
  }
}

// we assume these are safe because of the middleware
// and because tokens can be rejected from restarting the program

/* File: ModeratorActivityController */
func GetAllBanners(c *gin.Context){
  entires := getAllEntries()
  c.JSON(200 , entires)
}
/* File: ModeratorActivityController */
func BanUser(c *gin.Context){
  var mod_json ModAction
  c.BindJSON(&mod_json)
  name := getGet(c , "name")
  // verify in case mod has been removed but token still authenticates as one
  mod_level := checkModLevel(name)
  if mod_level == "none"{
    c.JSON(401 , gin.H{"error": "You are not a moderator"})
    return
  }

  ip := c.ClientIP()
  if mod_json.Target == "" || mod_json.Hard == nil {
    c.JSON(401 , gin.H{"error": "Fields missing"})
    return
  }
  createNewBan(mod_json.Target, *mod_json.Hard , ip)

  hard_ban_str := ""
  if *mod_json.Hard == 1 {
    hard_ban_str = "hard"
  } else{
    hard_ban_str = "soft"
  }
  c.JSON(200 , gin.H{"log": "User " + mod_json.Target + " was " + hard_ban_str + " banned"})
}
/* File: ModeratorActivityController */
func DeleteAll(c *gin.Context){
  var mod_json ModAction
  c.BindJSON(&mod_json)
  name := getGet(c , "name")
  // verify in case mod has been removed but token still authenticates as one
  mod_level := checkModLevel(name)
  if mod_level == "none"{
    c.JSON(401 , gin.H{"error": "You are not a moderator"})
    return
  }
  if mod_json.Target == "" {
    c.JSON(400 , gin.H{"error": "Fields missing"})
    return
  }
  removeAllUserImages(mod_json.Target)
  removeUserFromDatabase(mod_json.Target)

  c.JSON(200 , gin.H{"log": "User " + mod_json.Target + " was purged"})
}
/* File: ModeratorActivityController */
func DeleteIndividual(c *gin.Context){
  var mod_json ModAction
  c.BindJSON(&mod_json)
  name := getGet(c , "name")
  // verify in case mod has been removed but token still authenticates as one
  mod_level := checkModLevel(name)
  if mod_level == "none"{
    c.JSON(401 , gin.H{"error": "You are not a moderator"})
    return
  }
  if mod_json.Target == "" || mod_json.URI == ""{
    c.JSON(400 , gin.H{"error": "Fields missing"})
    return
  }
  removeIndividualBannerFromImages(mod_json.URI)
  removeIndividualBannerFromDB(mod_json.URI)

  c.JSON(200 , gin.H{"log": mod_json.Target + "'s image was pruned"})
}

/* MISC FUNCTIONS */

func getParam(c *gin.Context, key string) string{
  q := c.Request.URL.Query()
  var value string
  if key_arr, key_exists := q[key] ; key_exists && len(key_arr) != 0{
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
