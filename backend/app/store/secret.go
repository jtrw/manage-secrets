package store

import (
  //"fmt"
  "github.com/nilBora/bolt"
  jbolt "manager-secrets/backend/app/store/jbolt"
)

var db *bolt.DB;

func Init() {
  db = jbolt.Open("my.db")
  //defer db.Close()
}

func Get(bucket, key string) string {
    return jbolt.Get(db, bucket, key)
}

func Set(bucket, key, value string) {
    jbolt.Set(db, bucket, key, value)
}