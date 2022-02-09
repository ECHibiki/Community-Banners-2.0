package templater

import (
  "io/ioutil"
  "bytes"
  "fmt"
  "log"

  "github.com/tyler-sommer/stick"
)
var env *stick.Env
func init(){

  fmt.Println("Starting Templater...")
  env = stick.New(nil)

  m := map[string]string{"user": "Vermin" , "permision": "Admin"}
  m_typecast := StringMapToStickValue(m)

  buff := new(bytes.Buffer)
  if err := env.Execute("\tHello, {{permision}} {{ user }}!", buff, m_typecast ); err != nil {
    log.Fatal(err)
  }
  fmt.Println(buff);

  fmt.Println("...Templater Initialized")
}

func ReturnFilledTemplate(template_path string, value_map map[string]string) string{
  template_bytes, err := ioutil.ReadFile(template_path)
  if err != nil{
    panic(err)
  }
  value_typecast := StringMapToStickValue(value_map)
  template_buffer := new(bytes.Buffer)
  if err := env.Execute(string(template_bytes), template_buffer, value_typecast ); err != nil {
    log.Fatal(err)
  }
  return template_buffer.String()
}

func StringMapToStickValue(value_map map[string]string) map[string]stick.Value {
  m_typecast := map[string]stick.Value{}
  for key, value := range value_map {
    m_typecast[key] = stick.Value(value)
  }
  return m_typecast
}