package controllers

import (
  "net/url"
  "fmt"

  "github.com/gin-gonic/gin"
  "github.com/ECHibiki/Community-Banners-2.0/templater"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)
// NAME CHECK
func returnAdPage(name string, size string , ip string , board string) string{
  ad_data := getRandomEntry(name , size , ip , board)
  if len(ad_data) == 0{
    return "asdf no ads";
  }
  ad_data[0]["url"] = fmt.Sprintf("/req?s=%v&f=%v" ,
     url.QueryEscape(ad_data[0]["url"]) , ad_data[0]["uri"] )
  ad_data[0]["name"] = ad_data[0]["name"]
  return templater.ReturnFilledTemplate("./templates/banner.html" , ad_data[0])
}
// NAME CHECK
func getRandomEntry(name string,  size string , ip string , board string) []map[string]string{
  if size != "wide" && size != "small"{
    size = "%"
  }
  filter := ""
  var filter_args []interface{}
  if board == ""{
    for _, r_board := range controller_settings.RestrictedBoards{
      if filter == "" {
        filter += " ( board != ? "
      } else{
        filter += " AND board != ? "
      }
      filter_args = append(filter_args , r_board)
    }
    filter += " ) AND "
  } else{
    filter += " ( board = ? OR board = '' ) AND "
    filter_args = append(filter_args , board)
  }

  query := fmt.Sprintf(`
    SELECT ads.fk_name AS name, uri , url, size, clicks, board FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    WHERE %s
      ( hardban IS NULL OR bans.ip = ? OR bans.fk_name = ? )
      AND ads.size LIKE ?
     ORDER BY RAND() LIMIT 1;
    ` , filter )
  query_arr := append(filter_args , []interface{}{ ip , name , size ,}...)
  result , err := bannerdb.Query(query , query_arr)
  if err != nil{
    panic(err)
  }
  return result
}

func returnAdJson(name string, size string , ip string , board string) (gin.H){
  ad_data := getRandomEntry(name , size , ip , board)
  if len(ad_data) == 0{
    return gin.H{
      "url": "",
      "uri": "",
      "name":"",
      "size":"",
      "clicks":"",
      "board":"",
    }
  }
  ad_data[0]["url"] = fmt.Sprintf("/req?s=%v&f=%v" ,
     url.QueryEscape(ad_data[0]["url"]) , ad_data[0]["uri"] )
  // ads.fk_name , uri , url, size, clicks , board
  return gin.H{
    "url": ad_data[0]["url"],
    "uri": ad_data[0]["uri"],
    "name":ad_data[0]["name"],
    "size":ad_data[0]["size"],
    "clicks":ad_data[0]["clicks"],
    "board":ad_data[0]["board"],
  }
}

func getLimitedEntries(name string ,  ip string ) []map[string]string{
  filter := ""
  var filter_args []interface{}
  if name == "" || !checkBanned(name , ip) {
    filter += "WHERE (bans.fk_name = ? OR bans.ip != ? OR bans.ip IS NULL) "
    filter_args = append(filter_args , []interface{}{name, ip}...)
  }
  for _, board := range controller_settings.RestrictedBoards{
    if filter == "" {
      filter += "WHERE board != ? "
    } else{
      filter += " AND board != ? "
    }
    filter_args = append(filter_args , board)
  }

  entry_get , err := bannerdb.Query(fmt.Sprintf(`
    SELECT ads.fk_name AS name , uri , url, size, clicks, board FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    %s ORDER BY ads.id DESC` , filter) , filter_args)
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