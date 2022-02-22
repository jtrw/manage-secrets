package store

import (
  "fmt"
  jbolt "bolt"
  "github.com/nilBora/bolt"
)

var db *bolt.DB;

func Init() {
  fmt.Printf("Secret!!!")
}

func Get(db *bolt.DB, bucket, key string) string {
    bolt.Get(db, b, "secondSecret")
}