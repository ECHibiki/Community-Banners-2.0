package main

import (
  "fmt"

  "github.com/ECHibiki/Community-Banners-2.0/templater"
  "github.com/ECHibiki/Community-Banners-2.0/mailer"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
  "github.com/ECHibiki/Community-Banners-2.0/ginterface"
)

func main(){
  fmt.Println("Community-Banners-2.0 - ECVerniy 2022")
  port := ":4200"

  // yeah, you could be using init() but I want to control it
  templater.Init()
  mailer.Init()
  bannerdb.Init()

  // final action before thread halts
  ginterface.Init(port)
}