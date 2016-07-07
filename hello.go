package main

import (
  "bufio"
  "flag"
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
  var name string
  flag.StringVar(&name, "id", "", "identifies unique instance, usually a name")
  flag.Parse()

  if name == "" {
    fake, err := faker.New("en")
    if err != nil {
      fmt.Errorf("Error setting up name: %v", err)
      return
    }
    name = fake.FirstName()
  }

  fmt.Printf("hello %v\n", name)

  passphrase := ""   // get this from user eventually

  fmt.Printf("will put credentials here %v\n", "cred/"+name)
  dir := "cred/"+name
  _, tmp_err := os.Stat(dir)
  fmt.Printf("os.Stat(dir) err=%v\n", tmp_err);
  if _, err := os.Stat(dir); err == nil {
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
  _, server, err := v23.WithNewServer(ctx, "", ifc.HelloServer(service.Make()), security.AllowEveryone())

  // start server with random blessing
  // ep: (@host:port@...@blessing)
  // send an advertisement: with (localmomentID, endpoint)
  AdvertiseServer(ctx, nil, server, "")

  // receive an advertisement
  Scan()  // print who is nearby
  

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
    // f := ifc.HelloClient(*server)
    // ctx, cancel := context.WithTimeout(ctx, time.Minute)
    // defer cancel()
    // hello, _ := f.Get(ctx)
    // fmt.Println(hello)

  }

  fmt.Printf("bye bye.\n\n")

}
