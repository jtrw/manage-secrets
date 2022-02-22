package main
import (
  "fmt"
 // "log"
 // "time"
  //"github.com/nilBora/bolt"
  bolt "store/bolt"
)
const b = "MyBucket"

func main() {
  db := bolt.Open("my.db")
  defer db.Close()

  bolt.Set(db, b, "love", "golang")
  v := bolt.Get(db, b, "love")
  fmt.Printf("I love %s\n", v)
}