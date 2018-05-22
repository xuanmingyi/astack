package main

import (
  "os/exec"
  "fmt"
)

func volume_create(args []string){

}

func volume_list(args []string) {

}

func volume_delete(args []string) {

}

func volume_info(args []string) {
  cmd := exec.Command("vgdisplay", global_option.volume_vg)
  stdout, err1 := cmd.Output()
  if err1 != nil {
	panic_exit(err1.Error())
  }
  fmt.Println(string(stdout))
}
