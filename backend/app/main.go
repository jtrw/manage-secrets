package main
import (
  "fmt"
 // "log"
 // "time"
  secret "manager-secrets/backend/app/store"
)
const b = "test/secret"

func main() {
  secret.Init();
  secret.Set(b, "secondSecret", "888")
  v := secret.Get(b, "secondSecret")
  fmt.Printf("Secret: %s\n", v)
}