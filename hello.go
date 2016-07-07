package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "github.com/manveru/faker"

  "log"
  "./ifc"
  "./service"
  "v.io/v23"
  "v.io/v23/context"
  "v.io/v23/security"
  "v.io/x/ref/lib/discovery"
  vdiscovery "v.io/v23/discovery"
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

  //V23_CREDENTIALS => dir
  root_ctx, shutdown := v23.Init()  // create random in-memory principal
  // v23.getPrincipal(ctx)  temp, in memory
  defer shutdown()

  ctx, cancel := context.WithCancel(root_ctx)

  d, err := v23.NewDiscovery(ctx)
  if err != nil {
    log.Panic("Error from NewDiscovery: ", err)
  }

  _, server, err := v23.WithNewServer(ctx, "", ifc.HelloServer(service.Make()), security.AllowEveryone())
  // TODO: error check

  ad := vdiscovery.Advertisement{
		InterfaceName: "hello",
	}

  // this advertises all of our endpoints
  discovery.AdvertiseServer(ctx, d, server, "", &ad, nil)

  // receive an advertisement (scan for all endpoints)
  updates, err := d.Scan(ctx, "")
  if err != nil {
    log.Panic("Error listening: ", err)
  } 
  fmt.Printf("updates %v\n", updates)
  // ??? the following code seems to hang... if uncommented, "bye" doens't work
  // for update := range updates {
  //   fmt.Printf("update %v\n", update)
  // }

  fmt.Printf("type some text and press return.\n")
  fmt.Printf("to exit: type 'bye' and press return %v\n")

  // read some text
  text := ""

  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    text = scanner.Text()
    if text == "bye" { cancel(); break }
    fmt.Printf("echo %v\n", scanner.Text())
    // f := ifc.HelloClient(*server)
    // ctx, cancel := context.WithTimeout(ctx, time.Minute)
    // defer cancel()
    // hello, _ := f.Get(ctx)
    // fmt.Println(hello)

  }

  fmt.Printf("bye bye.\n\n")

}
