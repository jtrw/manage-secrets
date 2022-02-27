package store

import (
  //"fmt"
  jbolt "manager-secrets/backend/app/store/jbolt"
)

//var BoltDB *jbolt.Bolt

type Store struct {
	StorePath string
	DB jbolt.Bolt
}


func (s Store) NewStore() jbolt.Bolt {
    bolt := jbolt.Open(s.StorePath)

    return *bolt
}

func (s Store) Get(bucket, key string) string {
    return jbolt.Get(s.DB.DB, bucket, key)
}

// func Init() {
//   //var boltDB jbolt.Bolt
//   BoltDB = jbolt.Open("my.db")
//   //defer db.Close()
// }
//

//
// func Set(bucket, key, value string) {
//     jbolt.Set(BoltDB.DB, bucket, key, value)
// }