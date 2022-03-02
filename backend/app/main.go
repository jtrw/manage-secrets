package main
import (
  "fmt"
  "log"
  "os"
  //"os/signal"
   "reflect"
  //
  //"time"
  //"flag"
  secret "manager-secrets/backend/app/store"
  server "manager-secrets/backend/app/server"
  "github.com/jessevdk/go-flags"
  //"context"
  //"time"
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
    sec := secret.Store {
        StorePath: "my.db",
    }

    sec.JBolt = sec.NewStore();

    commandName := os.Args[1]
    if commandName == "run" {
        srv := server.Server {
            DataStore: sec,
            PinSize:   1,
            WebRoot:   "/",
            Version:   "1.0",
        }

        //go func () {
            if err := srv.Run(); err != nil {
                log.Printf("[ERROR] failed, %+v", err)
            }
        //}()

//         c := make(chan os.Signal, 1)
//         signal.Notify(c, os.Interrupt)
//         <-c
//
//         // Attempt a graceful shutdown
//         _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel()
        //srv.Shutdown(ctx)
        //time.Sleep(100)
    }

    if commandName == "kv" {
        command := os.Args[2]
        if command == "get" {
            keyStore := os.Args[3]
            v := sec.Get(b, keyStore)
            fmt.Printf("Secret: %s\n", v)
        }

        if command == "set" {
            keyStore := os.Args[3]
            valueStore := os.Args[4]
            sec.Set(b, keyStore, valueStore)
            fmt.Printf("Done!\n")
        }
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



    //secret.Init();
    //secret.Set(b, "secondSecret", "888")

}

