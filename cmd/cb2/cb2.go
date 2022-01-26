package main

import (
  "fmt"
  
  "github.com/ECHibiki/Community-Banners-2.0/ginterface"
  "github.com/ECHibiki/Community-Banners-2.0/mailer"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func main(){
  fmt.Println("Community-Banners-2.0 - ECVerniy 2022")
  port := ":4200"

  ginterface.Init(port)
  mailer.Init()
  bannerdb.Init()

}