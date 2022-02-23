package main
import (
  "fmt"
 // "log"
 // "time"
  //"github.com/nilBora/bolt"
  //"store/store"
  //jbolt "manager-secrets/backend/app/store"
  secret "manager-secrets/backend/app/store"
)
const b = "test/secret"

func main() {
//   db := jbolt.Open("my.db")
//   defer db.Close()
//
//   bolt.Set(db, b, "secondSecret", "888")
//   v := jbolt.Get(db, b, "secondSecret")
//   fmt.Printf("Secret: %s\n", v)

  secret.Init();
  secret.Set(b, "secondSecret", "888")
  v := secret.Get(b, "secondSecret")
  fmt.Printf("Secret: %s\n", v)
}