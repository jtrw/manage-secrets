package main
import (
  "fmt"
  "log"
  "os"
   "reflect"
  //
 // "time"
  //"flag"
  secret "manager-secrets/backend/app/store"
  server "manager-secrets/backend/app/server"
  "github.com/jessevdk/go-flags"
)
const b = "test/secret"

type Options struct {
   Name string `long:"name" description:"Your name, for a greeting" default:"Unknown"`
   Verbose string `short:"v" long:"verbose" description:"Show verbose debug information"`
}

type Commands struct {
   Run string `name:"run" description:"Start server"`
   Kv string `name:"kv" description:"Key value storage"`
}

func main() {
    commandName := os.Args[1]
    if commandName == "run" {
        srv := server.Server{
            PinSize:        1,
            WebRoot:        "/",
            Version:        "1.0",
        }
        if err := srv.Run(); err != nil {
            log.Printf("[ERROR] failed, %+v", err)
        }

        //server.Run();
    }

    t := reflect.TypeOf(Commands{})
    f, _ := t.FieldByName("Run")
    fmt.Println(f.Tag) // one:"1" two:"2"blank:""
    val, _ := f.Tag.Lookup("name")
    fmt.Printf("%s\n", val) // 1, true


    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Name", opts.Name)
    fmt.Println("Verbose: ", opts.Verbose)

    secret.Init();
    //secret.Set(b, "secondSecret", "888")
    v := secret.Get(b, "secondSecret")
    fmt.Printf("Secret: %s\n", v)
}

