package controllers

import (
  "net/url"
  "fmt"

  "github.com/gin-gonic/gin"
  "github.com/ECHibiki/Community-Banners-2.0/templater"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)
// NAME CHECK
func returnAdPage(name string, size string , ip string) string{
  ad_data := getRandomEntry(name , size , ip)
  if len(ad_data) == 0{
    return "asdf no ads";
  }
  ad_data[0]["url"] = fmt.Sprintf("/req?s=%v&f=%v" ,
     url.QueryEscape(ad_data[0]["url"]) , ad_data[0]["uri"] )
  ad_data[0]["name"] = ad_data[0]["name"]
  return templater.ReturnFilledTemplate("./templates/banner.html" , ad_data[0])
}
// NAME CHECK
func getRandomEntry(name string,  size string , ip string) []map[string]string{
  if size != "wide" && size != "small"{
    size = "%"
  }
  query_arr := []interface{}{size , ip , name}
  result , err := bannerdb.Query(`
    SELECT ads.fk_name AS name, uri , url, size, clicks FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    WHERE ( hardban IS NULL AND ads.size LIKE ? )
      OR bans.ip = ?
      OR bans.fk_name = ?
    ORDER BY RAND() LIMIT 1;
    ` , query_arr)
  if err != nil{
    panic(err)
  }
  return result
}

func returnAdJson(name string, size string , ip string) (gin.H){
  ad_data := getRandomEntry(name , size, ip)
  if len(ad_data) == 0{
    return gin.H{
      "url": "",
      "uri": "",
      "name":"",
      "size":"",
      "clicks":"",
    }
  }
  ad_data[0]["url"] = fmt.Sprintf("/req?s=%v&f=%v" ,
     url.QueryEscape(ad_data[0]["url"]) , ad_data[0]["uri"] )
  // ads.fk_name , uri , url, size, clicks
  fmt.Println(ad_data[0])
  return gin.H{
    "url": ad_data[0]["url"],
    "uri": ad_data[0]["uri"],
    "name":ad_data[0]["name"],
    "size":ad_data[0]["size"],
    "clicks":ad_data[0]["clicks"],
  }
}

func getLimitedEntries(name string ,  ip string ) []map[string]string{
  ban_filter_query := ""
  ban_filter_args := []interface{}{}
  if name == "" || !checkBanned(name , ip) {
    ban_filter_query = "WHERE bans.fk_name = ? OR bans.ip != ? OR bans.ip IS NULL"
    ban_filter_args = []interface{}{name, ip}
  }
  fmt.Println(fmt.Sprintf(`
    SELECT ads.fk_name AS name , uri , url, size, clicks FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    %s ORDER BY ads.id DESC` , ban_filter_query))
  entry_get , err := bannerdb.Query(fmt.Sprintf(`
    SELECT ads.fk_name AS name , uri , url, size, clicks FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    %s ORDER BY ads.id DESC` , ban_filter_query) , ban_filter_args)
  if err != nil{
    panic(err)
  }
  return entry_get
}

func checkBanned(name string , ip string) bool{
  ban_get , err := bannerdb.Query(`
    SELECT fk_name FROM bans
    WHERE fk_name = ? OR ip = ?
  `, []interface{}{name , ip})
  if err != nil{
    panic(err)
  }
  return len(ban_get) > 0
}