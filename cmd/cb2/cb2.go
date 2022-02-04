package main

import (
  "fmt"

  "github.com/ECHibiki/Community-Banners-2.0/ginterface"
)

func main(){
  fmt.Println("Community-Banners-2.0 - ECVerniy 2022")
  port := ":4200"

  // final action before thread halts
  ginterface.Init(port)
}