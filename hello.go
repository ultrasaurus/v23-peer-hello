package main

import (
  "fmt"
  "os"
  vsecurity "v.io/x/ref/lib/security"
  "github.com/manveru/faker"
)


func main() {
  fake, err := faker.New("en")
  if err != nil {
    fmt.Errorf("Error setting up name: %v", err)
    return
  }
  name := fake.FirstName()
  fmt.Printf("hello %v\n", name)

  passphrase := ""   // get this from user eventually

  fmt.Printf("will put stuff here %v\n", "cred/"+name)
  dir := "cred/"+name
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    p, err := vsecurity.CreatePersistentPrincipal(dir, []byte(passphrase))
    if err != nil {
      fmt.Errorf("Error creating principal: %v", err); return
    }
    blessings, err := p.BlessSelf(name)
    if err != nil {
      fmt.Errorf("BlessSelf(%q) failed: %v", name, err); return
    }
    if err := vsecurity.SetDefaultBlessings(p, blessings); err != nil {
      fmt.Errorf("could not set blessings %v as default: %v", blessings, err)
      return
    }
    fmt.Printf("Bless you.\n\n")
  }

}
