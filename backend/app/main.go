package main
import (
  "fmt"
  "log"
  "os"
  //"os/signal"
 //  "reflect"
  //
  //"time"
  //"flag"
  secret "manager-secrets/backend/app/store"
  server "manager-secrets/backend/app/server"
 // command "manager-secrets/backend/app/cmd"
  "github.com/jessevdk/go-flags"
   "net/http"
   "io/ioutil"
   "bytes"
  //"context"
  //"time"
  //"time"
)

const host = "http://127.0.0.1"
const port = "8080"

type Options struct {
   Name string `long:"name" description:"Your name, for a greeting" default:"Unknown"`
   Verbose string `short:"v" long:"verbose" description:"Show verbose debug information"`

   Type string `short:"t" long:"type" description:"Type content save content"`
   Host string `short:"h" long:"host" default:"http://127.0.0.1" description:"Host web server"`
   Port string `short:"p" long:"port" default:"8080" description:"Port web server"`
}

func main() {
    //command.Parse()
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    commandName := os.Args[1]
    switch commandName {
        case "run":
            makeRunCommand()
        case "kv":
            makeKvCommand(opts)
        default:
            log.Fatal("Command name not found")
    }
}

func makeRunCommand() {
    sec := secret.Store {
        StorePath: "my.db",
    }

    sec.JBolt = sec.NewStore();

    srv := server.Server {
        DataStore: sec,
        PinSize:   1,
        WebRoot:   "/",
        Version:   "1.0",
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}

func makeKvCommand(opts Options) {
    command := os.Args[2]

    switch command {
        case "get":
            makeGetSubCommand()
        case "set":
            makeSetSubCommand(opts)
        default:
            log.Fatal("Sub Command name not found")
    }
}

func makeGetSubCommand() {
    keyStore := os.Args[3]

    resp, err := http.Get(getApiAddr()+keyStore)
    if err != nil {
       log.Fatalln(err)
    }

     body, err := ioutil.ReadAll(resp.Body)
     if err != nil {
       log.Fatalln(err)
    }

    fmt.Printf("Secret: %s\n", body)
}

func makeSetSubCommand(opts Options) {
    keyStore := os.Args[3]
    valueStore := os.Args[4]
    dataByte := []byte(valueStore)
    responseBody := bytes.NewBuffer(dataByte)
    contentType := "text/plain"
    if opts.Type == "json" {
        contentType = "application/json"
    }

    resp, err := http.Post(getApiAddr()+keyStore, contentType, responseBody)
    if err != nil {
       log.Fatalln(err)
    }

     response, errRead := ioutil.ReadAll(resp.Body)
     if errRead != nil {
       log.Fatalln(errRead)
    }

    fmt.Printf("%s\n", response)
}

func getApiAddr() (string) {
    return host+":"+port+"/api/v1/kv/"
}

