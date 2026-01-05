package main

import (
	"fmt"
	"os/exec"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/fatih/color"
)

type Process struct{
	Id int
	ProcessName string
	ProcessID int
	Port int
}

var Processes []Process

func logError(err error){
	if err != nil{
		panic(err)
	}
}

func isInSlice(process Process) bool{
	for _, p := range Processes{
		if(process.ProcessID == p.ProcessID && process.Port == p.Port){
			return true
		}
	}
	return false
}

func returnCommandOutput(command string) []string{
	formattedCommand := strings.Split(command, " ")
	output, err := exec.Command(formattedCommand[0], formattedCommand[1:]...).Output()
	logError(err)

	outputSlice := strings.Split(string(output), "\n")
	fmt.Println(string(output))
	return outputSlice
}

func filterCommandOutput(outputSlice []string) []string{
	var filteredOutputSlice []string
	re := regexp.MustCompile(`users:\(\("(.+?)",pid=(\d+)`)
	for _, outputLine := range outputSlice{
		if re.MatchString(outputLine){
			filteredOutputSlice = append(filteredOutputSlice, outputLine)
		}
	}
	return filteredOutputSlice
}

func extractProcess(filteredOutputSlice []string){
	var id int 
	for _, outputLine := range filteredOutputSlice{
		re := regexp.MustCompile(`users:\(\("(.+?)",pid=(\d+)`)
		rePort := regexp.MustCompile(`(?:(?:\[[^\]]+\])|(?:\S+)):(\*|[A-Za-z][A-Za-z0-9_-]*|\d+)`)
		matchSlice := re.FindStringSubmatch(outputLine)
		portSlice := rePort.FindStringSubmatch(outputLine)

		var isPortAvailable bool = true
		process := matchSlice[1]
		pid, err := strconv.Atoi(matchSlice[2])
		logError(err)
		port, err := strconv.Atoi(portSlice[1])
		if err != nil{
			isPortAvailable = false
		}

		if !isPortAvailable{
			port = 0
		}
		
		p := Process{
			Id: id,
			ProcessName: process,
			ProcessID: pid,
			Port : port,
		}
		if(!isInSlice(p)){
			Processes = append(Processes, p)
			id++
		}
	}
	// fmt.Println(Processes)
}

func generateTable(){
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold},
			BG: renderer.Colors{color.BgBlack},
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			BG: renderer.Colors{color.BgBlack},
			Columns: []renderer.Tint{
				{
				FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow for column 0
				BG: renderer.Colors{color.BgBlack},
				},
			},
		},
	}

	table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewColorized(colorCfg)))
	table.Header("ID", "PROCESS", "PROCESS ID", "PORT")
	for _, p := range Processes{
		var port string
		if p.Port != 0{
			port = strconv.Itoa(p.Port)
		}else{
			port = "Not Available"
		}
		table.Append(p.Id, p.ProcessName, p.ProcessID, port)
	}
	table.Render()
}

func main(){
	outputSlice := returnCommandOutput("ss -tlp")
	filteredOutputSlice := filterCommandOutput(outputSlice)
	extractProcess(filteredOutputSlice)
	generateTable()
}	