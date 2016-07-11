package main

import (
  "bufio"
  "fmt"
  "os"
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
  if err != nil {
    log.Panic("Error from WithNewServer: ", err)
  }

  ad := vdiscovery.Advertisement{
		InterfaceName: "hello",
	}

  // this advertises all of our endpoints
  _, err = discovery.AdvertiseServer(ctx, d, server, "", &ad, nil)
  if err != nil {
    log.Panic("Error from AdvertiseServer: ", err)
  }

  // receive an advertisement (scan for all endpoints)
  updates, err := d.Scan(ctx, "")
  if err != nil {
    log.Panic("Error listening: ", err)
  }

  fmt.Printf("\n\nto exit: type 'bye' and press return\n\n")


  go func() {
    // the flow is to keep reading from the channel till it is closed
    // (which won't happen until the scan is aborted by cancelling the context).
    // in this example, the cancel happens via user input below,
    // so this is has to happen in parallel in a go routine
    for update := range updates {
      if (update.IsLost()) {
        fmt.Printf("lost: %v\n", update.Id())
      } else {
       fmt.Printf("new: %v %v\n", update.Id(), update.InterfaceName())
       addressList := update.Addresses()
       for addrIndex := range addressList {
         fmt.Printf("  %v %v\n", addrIndex, addressList[addrIndex])
       }
     }
    }
  }()

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
