package store

import (
 // "fmt"
  "log"
  "time"
  "github.com/nilBora/bolt"
)

func Open(file string) *bolt.DB {
  db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
  if err != nil {
    //handle error
    log.Fatal(err)
  }
  return db
}

func Set(db *bolt.DB, bucket, key, value string) {
  db.Update(func(tx *bolt.Tx) error {
    b, _ := tx.CreateBucketIfNotExists([]byte(bucket))
    err := b.Put([]byte(key), []byte(value))
    return err
  })
}

func Get(db *bolt.DB, bucket, key string) string {
  s := ""
  db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket))
    s = string(b.Get([]byte(key)))
    return nil
  })
  return s
}

func Del(db *bolt.DB, bucket, key string) {
  db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket))
    b.Delete([]byte(key))
    return nil
  })
}