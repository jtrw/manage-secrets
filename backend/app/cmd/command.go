package command

import (
  "fmt"
  "reflect"
  "os"
)

type Commands struct {
   run string `name:"run" method:"MakeRun" description:"Start server"`
   kv string `name:"kv" method:"MakeKv" description:"Key value storage"`
}

func Parse() {
    var cmd Commands

    t := reflect.TypeOf(cmd)

    run, _ := t.FieldByName(os.Args[1])

    MethodName, _ := run.Tag.Lookup("method")

    reflect.ValueOf(&cmd).MethodByName(MethodName).Call([]reflect.Value{})
//     if ok {
//          cmd.RunServer()
//     }

   // fmt.Println(f.Tag) // one:"1" two:"2"blank:""
   // val, _ := f.Tag.Lookup("name")
   // fmt.Printf("%s\n", val) // 1, true
}


func (c Commands) MakeRun() {
    fmt.Printf("Server started!")
}

func (c Commands) MakeKv() {
     fmt.Printf("MakeKv started!")
}