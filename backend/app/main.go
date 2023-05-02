package main
import (
  "fmt"
  "log"
  "os"
  secret "manager-secrets/backend/app/store"
  server "manager-secrets/backend/app/server"
  output "manager-secrets/backend/app/cmd"
  "github.com/jessevdk/go-flags"
  "net/http"
  "io/ioutil"
  "bytes"
  "github.com/joho/godotenv"
)

const ENV_HOST_KEY  = "JTRW_MANAGER_SECRETS_HOST"
const ENV_PORT_KEY  = "JTRW_MANAGER_SECRETS_PORT"

type Options struct {
   Type string `short:"t" long:"type" description:"Type content save content"`
   Host string `short:"h" long:"host" default:"0.0.0.0" description:"Host web server"`
   Port string `short:"p" long:"port" default:"8080" description:"Port web server"`
   StoragePath string `short:"s" long:"storage_path" default:"/var/tmp/jtrw_manager_s.db" description:"Storage Path"`
}

type KvCommand struct {
	ContentType string
}

type MainCommand struct {
    CommandName string
    Opts Options
}

func main() {
    godotenv.Load()

    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

     mc := MainCommand {
        CommandName: os.Args[1],
        Opts: opts,
    }
    mc.Start()
}

func (mc MainCommand) Start() {
    switch mc.CommandName {
        case "run":
            mc.makeRunCommand()
        case "kv":
            mc.makeKvCommand()
        default:
            log.Fatal("Command name not found")
    }
}

func (mc MainCommand) makeRunCommand() {
    sec := secret.Store {
        StorePath: mc.Opts.StoragePath,
    }

    sec.JBolt = sec.NewStore()

    host := mc.Opts.Host
    port := mc.Opts.Port
    hostEnv := os.Getenv(ENV_HOST_KEY)
    if len(hostEnv) > 0 {
        host = hostEnv
    }
    portEnv := os.Getenv(ENV_PORT_KEY)
    if len(portEnv) > 0 {
        port = portEnv
    }

    srv := server.Server {
        DataStore: sec,
        Host: host,
        Port: port,
        PinSize:   1,
        WebRoot:   "/",
        Version:   "1.0",
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}

func (mc MainCommand) makeKvCommand() {
    kvc := KvCommand {
        ContentType: mc.Opts.Type,
    }
    command := os.Args[2]

    switch command {
        case "get":
            kvc.makeKvGetCommand()
        case "set":
            kvc.makeKvSetCommand()
        default:
            log.Fatal("Sub Command name not found")
    }
}

func (kvc KvCommand) makeKvGetCommand() {
    keyStore := os.Args[3]

    resp, err := http.Get(getApiAddr()+"kv/"+keyStore)
    if err != nil {
       log.Fatalln(err)
    }

     body, err := ioutil.ReadAll(resp.Body)
     if err != nil {
       log.Fatalln(err)
    }

    output.Print(body)
}

func (kvc KvCommand) makeKvSetCommand() {
    keyStore := os.Args[3]
    valueStore := os.Args[4]
    dataByte := []byte(valueStore)
    responseBody := bytes.NewBuffer(dataByte)
    contentType := "text/plain"
    if kvc.ContentType == "json" {
        contentType = "application/json"
    }

    resp, err := http.Post(getApiAddr()+"kv/"+keyStore, contentType, responseBody)
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
    host := os.Getenv(ENV_HOST_KEY)
    port := os.Getenv(ENV_PORT_KEY)
    return "http://"+host+":"+port+"/api/v1/"
}

