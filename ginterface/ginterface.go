package ginterface

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ECHibiki/Community-Banners-2.0/bannerdb"
	"github.com/ECHibiki/Community-Banners-2.0/bannerjwt"
	"github.com/ECHibiki/Community-Banners-2.0/controllers"
	"github.com/gin-gonic/gin"
)

type GinSettings struct {
	Domain         string
	RejectCreation bool
	ReleaseMode    string
}

var gin_engine *gin.Engine
var gin_settings GinSettings

func Init(port string) {
	fmt.Println("\nGin Interface initialization...")

	setting_json_bytes, err := ioutil.ReadFile("./settings/gin-settings.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(setting_json_bytes, &gin_settings)

	// NGINX handles statics such as .html, .js, .css and image media

	gin.SetMode(gin_settings.ReleaseMode)
	gin_engine = gin.New()
	gin_engine.Use(gin.Recovery())
	gin_engine.SetTrustedProxies([]string{"127.0.0.1"})
	gin_engine.Use(JWTDecodeMiddleware())
	{
		gin_engine.GET("/banner", controllers.GenerateAdPage)
		gin_engine.GET("/req", controllers.RedirectSiteRequest)

		public_group := gin_engine.Group("/api/")
		{
			public_group.GET("banner", controllers.GenerateAdJSON)
			public_group.GET("all", controllers.GetLimitedInfo)
			if gin_settings.RejectCreation {
				public_group.POST("create", controllers.RejectUserCreation)
			} else {
				public_group.POST("create", controllers.CreateNewUser)
			}
			public_group.POST("login", controllers.LoginUser(gin_settings.Domain))
		}

		logged_group := public_group.Group("user/")
		logged_group.Use(AuthenticationMiddleware())
		logged_group.Use(BannedMiddleware())
		{
			logged_group.GET("details", controllers.AccessInfo)
			logged_group.POST("details", controllers.CreateBanner)
			logged_group.POST("removal", controllers.RemoveBanner)
			logged_group.POST("edit", controllers.EditBanner)
			logged_group.POST("token", controllers.TestToken(gin_settings.Domain))

			mod_group := logged_group.Group("mod/")
			mod_group.Use(ModeratorMiddleware())
			{
				mod_group.GET("all", controllers.GetAllBanners)
				mod_group.POST("ban", controllers.BanUser)
				mod_group.POST("purge", controllers.DeleteAll)
				mod_group.POST("individual", controllers.DeleteIndividual)
				mod_group.POST("edit", controllers.EditIndividualBanner)
			}
		}

	}
	gin_engine.Run(port)
	fmt.Println("...Gin Interface initialized")
}

/* middleware */
// return function instead of handling directly to potentially pass in command line arguments on initialization

func JWTDecodeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// VALIDATE JWT
		token_string, _ := c.Cookie("freeadstoken")
		name, is_donor, is_mod, err := bannerjwt.IsAuth(token_string)
		c.Set("name", name)
		c.Set("is_donor", is_donor)
		c.Set("is_mod", is_mod)
		c.Set("valid_jwt", err == nil)
		c.Next()
	}
}
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// VALIDATE JWT
		valid := c.MustGet("valid_jwt").(bool)
		if !valid {
			// ABORT IF INVALID
			c.JSON(401, gin.H{"error": "You are not logged in"})
			c.Abort()
		}
		c.Next()
	}
}
func BannedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.MustGet("name")

		var query_params []interface{}
		query_params = append(query_params, name)
		banned_rows, err := bannerdb.Query("SELECT * FROM bans WHERE fk_name=? AND hardban=1", query_params)
		if err != nil {
			// ABORT IF INVALID
			c.JSON(500, gin.H{"error": "Ban Search Error"})
			c.Abort()
		}
		if len(banned_rows) != 0 {
			c.SetCookie("freeadstoken", "", -1, "/", gin_settings.Domain, true, true)
			c.JSON(401, gin.H{"error": "You've been banned..."})
			c.Abort()
		}
		c.Next()
	}
}
func ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.MustGet("name")

		var query_params []interface{}
		query_params = append(query_params, name)
		mod_rows, err := bannerdb.Query("SELECT * FROM mods WHERE fk_name=?", query_params)
		if err != nil {
			// ABORT IF INVALID
			c.JSON(500, gin.H{"error": "Mod Search Error"})
			c.Abort()
		}
		if len(mod_rows) == 0 {
			c.JSON(401, gin.H{"error": "You are not a moderator"})
			c.Abort()
		}
		c.Next()
	}
}
