package main

import (
	"flag"
	"fmt"
	"os/exec"
	"portmanager/helpers"
	"portmanager/internal"
	"strconv"
)

func killFlag(port int){
	var isUtilized bool = false
	for _, process := range internal.Processes{
		if(process.Port == port){
			_, err := exec.Command("kill", strconv.Itoa(process.ProcessID)).Output()
			helpers.LogError(err)
			fmt.Printf("Process %s utilising port %d killed.\n", process.ProcessName, process.Port)
			isUtilized = true
		}
	}
	if !isUtilized {
		fmt.Printf("Port %d is not being utilised by any process.\n", port)
	}
}

func main(){
	internal.GenerateProcesses()

	kill := flag.Int("kill",0,"Kills the process occupying the specified port")
	flag.Parse()

	if(*kill != 0){
		killFlag(*kill)
	}
	if(*kill == 0){
		internal.Runner()
	}
}