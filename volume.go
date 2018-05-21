package main

import (
  "os/exec"
  "fmt"
  "io/ioutil"
)

func volume_create(args []string){

}

func volume_list(args []string) {

}

func volume_delete(args []string) {

}

func volume_info(args []string) {
  cmd := exec.Command(fmt.Sprintf("sudo vgdisplay %s", global_option.volume_vg))
  stdout, err1 := cmd.StdoutPipe()
  if err1 != nil {
    panic_exit(err1.Error())
  }
  defer stdout.Close()
  cmd.Start()
  content, err2 := ioutil.ReadAll(stdout)
  if err2 != nil {
    panic_exit(err2.Error())
  }
  fmt.Println(string(content))
}
