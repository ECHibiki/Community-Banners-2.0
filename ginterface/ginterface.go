package ginterface

import (
  "fmt"
  "time"
  "encoding/json"
  "io/ioutil"

  "github.com/gin-gonic/gin"
  "github.com/ECHibiki/Community-Banners-2.0/controllers"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
  "github.com/ECHibiki/Community-Banners-2.0/bannerjwt"
)

type GinSettings struct{
  Domain string
}

var gin_engine *gin.Engine
var gin_settings GinSettings

func Init(port string){
  fmt.Println("\nGin Interface initialization...")

  setting_json_bytes, err := ioutil.ReadFile("./settings/gin-settings.json")
  if err != nil{
    panic(err)
  }
  json.Unmarshal(setting_json_bytes, &gin_settings)

  // NGINX handles statics such as .html, .js, .css and image media
  gin_engine = gin.Default()
  {
    gin_engine.GET("/banner", controllers.GenerateAdPage)
    gin_engine.GET("/req", controllers.RedirectSiteRequest)

    public_group := gin_engine.Group("/api/")
    {
      public_group.GET("banner", controllers.GenerateAdJSON)
      public_group.GET("all", controllers.GetLimitedInfo)
      public_group.POST("create", controllers.CreateNewUser)
      public_group.POST("login", controllers.LoginUser)
      public_group.POST("mod/login", controllers.LoginMod)
    }

    logged_group := gin_engine.Group("/api/")
    logged_group.Use(AuthenticationMiddleware())
    logged_group.Use(BannedMiddleware())
    {
      logged_group.GET("details", controllers.AccessInfo)
      logged_group.POST("details", controllers.CreateInfo)
      logged_group.POST("removal", controllers.RemoveInfo)

      mod_group := logged_group.Group("mod/")
      mod_group.Use(ModeratorMiddleware())
      {
        mod_group.GET("all", controllers.GetAllInfo)
        mod_group.POST("ban", controllers.BanUser)
        mod_group.POST("purge", controllers.DeleteAll)
        mod_group.POST("individual", controllers.DeleteIndividual)
      }
    }

  }
  gin_engine.Run(port)
  fmt.Println("...Gin Interface initialized")
}

/* middleware */
// return function instead of handling directly to potentially pass in command line arguments on initialization

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
    // VALIDATE JWT
    token_string, _ := c.Cookie("jwt")
    name, err := bannerjwt.IsAuth(token_string)
    if err != nil{
      // ABORT IF INVALID
      c.JSON(401 , gin.H{"warn": "You are not logged in"})
      c.Abort()
    }
    token, err := bannerjwt.CreateToken(name)
    if err != nil{
      // ABORT IF INVALID
      c.SetCookie("jwt", "", -1, "/", gin_settings.Domain, true, true)
      c.JSON(500 , gin.H{"error": "Key Create Error"})
      c.Abort()
    }
    c.SetCookie("jwt", token, int(time.Now().Add(time.Hour * 24).Unix()), "/",
      gin_settings.Domain, true, true)
    c.Set("name", name)
    c.Next()
	}
}
func BannedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
    name := c.MustGet("name")

    var query_params []interface{}
    query_params = append( query_params, name)
    banned_rows , err := bannerdb.Query("SELECT * FROM bans WHERE fk_name=? AND hardban=1" , query_params)
    if err != nil{
      // ABORT IF INVALID
      c.JSON(500 , gin.H{"error": "Ban Search Error"})
      c.Abort()
    }
    if len(banned_rows) != 0{
      c.SetCookie("jwt", "", -1, "/", gin_settings.Domain, true, true)
      c.JSON(401, gin.H{"warn": "You've been banned..."})
      c.Abort()
    }
    c.Next()
	}
}
func ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
    name := c.MustGet("name")

    var query_params []interface{}
    query_params = append( query_params, name)
    mod_rows, err := bannerdb.Query("SELECT * FROM mods WHERE fk_name=?" , query_params)
    if err != nil{
      // ABORT IF INVALID
      c.JSON(500 , gin.H{"error": "Mod Search Error"})
      c.Abort()
    }
    if len(mod_rows) == 0{
      c.JSON(401, gin.H{"warn": "You are not a moderator"})
      c.Abort()
    }
    c.Next()
	}
}
