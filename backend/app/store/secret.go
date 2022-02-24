package store

import (
  //"fmt"
  jbolt "manager-secrets/backend/app/store/jbolt"
)

var BoltDB *jbolt.Bolt


func Init() {
  //var boltDB jbolt.Bolt
  BoltDB = jbolt.Open("my.db")
  //defer db.Close()
}

func Get(bucket, key string) string {
    return jbolt.Get(BoltDB.DB, bucket, key)
}

func Set(bucket, key, value string) {
    jbolt.Set(BoltDB.DB, bucket, key, value)
}