package main

import (
  "bufio"
  "fmt"
  "os"
  vsecurity "v.io/x/ref/lib/security"
  "github.com/manveru/faker"

  "log"
  "./ifc"
  "./service"
  "v.io/v23"
  _ "v.io/x/ref/runtime/factories/generic"
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
  if _, err := os.Stat(dir); os.IsExist(err) {
    fmt.Printf("should load pricipal and do something with it.\n\n")

    // p, err := vsecurity.LoadPersistentPrincipal(dir, []byte(passphrase))
    // if err != nil {
    //   fmt.Errorf("Error loading principal: %v", err); return
    // }

  } else {
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

  ctx, shutdown := v23.Init()
  defer shutdown()
  _, _, err = v23.WithNewServer(ctx, "", ifc.HelloServer(service.Make()), nil)
  if err != nil {
    log.Panic("Error listening: ", err)
  }

  // read some text
  text := ""

  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    text = scanner.Text()
    if text == "bye" { break }
    fmt.Printf("echo %v\n", scanner.Text())
  }

  fmt.Printf("bye bye.\n\n")

}
