package main

import (
  "fmt"
  "os"
  vsecurity "v.io/x/ref/lib/security"
)


func main() {
  fmt.Printf("hello v23\n")
  passphrase := ""   // get this from user eventually

  if _, err := os.Stat("cred"); os.IsNotExist(err) {
    p, err := vsecurity.CreatePersistentPrincipal("cred", []byte(passphrase))
    if err != nil {
      fmt.Printf("Error creating principal: %v", err)
      return
    }
    fmt.Printf("Principal created: %v", p)
  }

}
