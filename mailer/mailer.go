package mailer

import (
  "fmt"
  "time"
  "strconv"
  "io/ioutil"
  "net/smtp"
  "encoding/json"

  "github.com/ECHibiki/Community-Banners-2.0/templater"
)

type MailSettings struct{
  Host string
  Pass string
  Recipients []string
  SendInterval int64
}

var mail_settings MailSettings
var last_sent_mail int64

func Init(){
  fmt.Println("\nMailer initialization...")
  setting_json_bytes, err := ioutil.ReadFile("./settings/mail-settings.json")
  if err != nil{
    panic(err)
  }
  json.Unmarshal(setting_json_bytes, &mail_settings)
  fmt.Println("...Mailer Interface initialized")

  // fmt.Println("\n\nRUNNING MAIL TEST")
  // fmt.Println("\tNOTE: USE BASE64 EMBEDDED IMAGES")
  // test()
  // fmt.Println("FINISHED MAIL TEST")
}

//https://medium.com/vacatronics/sending-email-with-go-23ae14050914
func sendEmailToAll(mail_body string , mail_title string) ( string ){
  if(time.Now().Unix() - last_sent_mail < mail_settings.SendInterval) {
    fmt.Println("Mail timeout " + strconv.FormatInt(time.Now().Unix() - last_sent_mail , 10) )
    return "Mail timeout " + strconv.FormatInt(time.Now().Unix() - last_sent_mail , 10)
  }

  // Configuration
  from := mail_settings.Host
  password := mail_settings.Pass
  to := mail_settings.Recipients
  smtpHost := "smtp.gmail.com"
  smtpPort := "587"

  message := []byte("Subject: " + mail_title + "\r\n" +
    "Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		mail_body + "\r\n")

  // Create authentication
  auth := smtp.PlainAuth("", from, password, smtpHost)

  // Send actual message
  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
  if err != nil {
    fmt.Println("Mail Error: ", err)
    return err.Error()
  }
  last_sent_mail = time.Now().Unix()
  return "Sent"
}

func SendBannerEmail(name string, file_base64 string, url string, uri string) string {
  // Get the template
  params := map[string]string{
    "time": time.Now().Format(time.UnixDate) + " - UnixDate" ,
    "name": name ,
    "url": url ,
    "uri": uri ,
    "base64": file_base64 ,

  }
  parsed_template := templater.ReturnFilledTemplate("./templates/banner-mail-notice.html" , params)
  // Send as email
  response := sendEmailToAll(parsed_template, "Test banner notification")
  // confirm response
  return response
}

/* Simple test of mail functions */

func test(){
  // Get the template
  params := map[string]string{
    "time": time.Now().Format(time.UnixDate) + " - UnixDate" ,
    "name": "test-name" ,
    "url": "https://art.kissu.moe/storage/image/nOPdjYNHKdr2BVvIVwhGpt1FI6aH2X5stqMV6D7p.gif" ,
    "fname": "https://art.kissu.moe/storage/image/nOPdjYNHKdr2BVvIVwhGpt1FI6aH2X5stqMV6D7p.gif" ,
    "err": "An error message" ,
  }
  parsed_template := templater.ReturnFilledTemplate("./templates/banner-mail-notice.html" , params)
  // Send as email
  response := sendEmailToAll(parsed_template, "Test banner notification")
  // confirm response
  fmt.Println(response)
}