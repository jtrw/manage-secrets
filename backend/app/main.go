package main
import (
  "fmt"
  "log"
  "os"
 // "time"
  //"flag"
  secret "manager-secrets/backend/app/store"
  "github.com/jessevdk/go-flags"
)
const b = "test/secret"

type Options struct {
   Name string `long:"name" description:"Your name, for a greeting" default:"Unknown"`
   Verbose string `short:"v" long:"verbose" description:"Show verbose debug information"`
}

func main() {

    myProgramName := os.Args[1]

    fmt.Println(myProgramName)

    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Name", opts.Name)
    fmt.Println("Verbose: ", opts.Verbose)

    secret.Init();
    secret.Set(b, "secondSecret", "888")
    v := secret.Get(b, "secondSecret")
    fmt.Printf("Secret: %s\n", v)
}

