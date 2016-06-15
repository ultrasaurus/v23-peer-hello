package main

import (
  "fmt"
  vsecurity "v.io/x/ref/lib/security"
)


func main() {
  fmt.Printf("hello v23\n")
  passphrase := ""   // get this from user eventually
  p, err := vsecurity.CreatePersistentPrincipal("cred", []byte(passphrase))
  if err != nil {
    fmt.Printf("Error creating principal: %v", err)
    return
  }
  fmt.Printf("Principal created: %v", p)

}
