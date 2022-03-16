package cmd
import (
    "fmt"
    "encoding/json"
    "log"
    secret "manager-secrets/backend/app/store"
)

func Print(body []byte) {

     result := &secret.Message{}

     err := json.Unmarshal(body, result)
     if err != nil {
            log.Fatal("Error Json")
     }

    fmt.Printf("======= Metadata =======\n")
    fmt.Printf("Key    Value\n")
    fmt.Printf("---    -----\n")
    fmt.Printf("Exp     %s\n", result.Exp)

    fmt.Printf("\n======= Data =======\n")
    fmt.Printf("Key    Value\n")
    fmt.Printf("---    -----\n")
    fmt.Printf("%s      %s\n", result.Key, result.Data)

}