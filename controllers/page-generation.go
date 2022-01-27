package controllers

import (
  "net/url"
  "fmt"

  "github.com/ECHibiki/Community-Banners-2.0/templater"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func returnAdPage(name string , size string , ip string) string{
  ad_data := getRandomEntry(name , size , ip)
  if len(ad_data) == 0{
    return "asdf no ads";
  }
  ad_data[0]["url"] = fmt.Sprintf("/req?s=%v&f=%v" ,
     url.QueryEscape(ad_data[0]["url"]) , ad_data[0]["uri"] )

  return templater.ReturnFilledTemplate("./templates/banner.html" , ad_data[0])
}

func getRandomEntry(name string , size string , ip string) []map[string]string{
  if size == ""{
    size = "%"
  }
  query_arr := []interface{}{size , ip , name}
  result , err := bannerdb.Query(`
    SELECT ads.fk_name , uri , url, size, clicks FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
      ( WHERE hardban IS NULL AND ads.size LIKE ? )
      OR bans.ip = ?
      OR bans.fk_name = ?
    ORDER BY RAND() LIMIT 1;
    ` , query_arr)
  if err != nil{
    panic(err)
  }
  return result
}