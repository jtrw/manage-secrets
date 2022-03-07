package store

import (
  //"fmt"
  "time"
  "encoding/json"
  jbolt "manager-secrets/backend/app/store/jbolt"
   log "github.com/go-pkgz/lgr"
)

//var BoltDB *jbolt.Bolt

var bucketDefault = []byte("secrets")

type Store struct {
	StorePath string
	JBolt jbolt.Bolt
}

type Message struct {
	Key     string
	Bucket  string
	Exp     time.Time
	Data    string
	//Data    []byte
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

func (s Store) Save(msg *Message) {
    _, jerr := json.Marshal(msg)
    if jerr != nil {
        return
        //return jerr
    }
    log.Printf("Secrert Save:");
    log.Printf(msg.Bucket);
    log.Printf(msg.Key);
    jbolt.Set(s.JBolt.DB, msg.Bucket, msg.Key, "TTTT!!!!")
}

// func (s Store) Load(msg *Message) (result *Message, err error) {
//     log.Printf("Secrert Load:");
//     log.Printf(msg.Bucket);
//     log.Printf(msg.Key);
//     val := jbolt.Get(s.JBolt.DB, msg.Bucket, msg.Key)
//
//     result = &Message{}
//
//     return json.Unmarshal([]byte(val), result)
// }