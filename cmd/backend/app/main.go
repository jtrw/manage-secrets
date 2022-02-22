package main
import (
  "fmt"
 // "log"
 // "time"
  //"github.com/nilBora/bolt"
  bolt "store/bolt"
)
const b = "test/secret"

func main() {
  db := bolt.Open("my.db")
  defer db.Close()

  bolt.Set(db, b, "secondSecret", "888")
  v := bolt.Get(db, b, "secondSecret")
  fmt.Printf("Secret: %s\n", v)
}