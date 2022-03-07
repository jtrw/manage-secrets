package store

import (
  //"fmt"
  "time"
  jbolt "manager-secrets/backend/app/store/jbolt"
)

//var BoltDB *jbolt.Bolt

type Store struct {
	StorePath string
	JBolt jbolt.Bolt
}

type Message struct {
	Key     string
	Exp     time.Time
	Data    []byte
	PinHash string
	Errors  int
}

func (s Store) NewStore() jbolt.Bolt {
    bolt := jbolt.Open(s.StorePath)

    return *bolt
}

func (s Store) Get(bucket, key string) string {
    return jbolt.Get(s.JBolt.DB, bucket, key)
}

func (s Store) Set(bucket, key, value string) {
    jbolt.Set(s.JBolt.DB, bucket, key, value)
}