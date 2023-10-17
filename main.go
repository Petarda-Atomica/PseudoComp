package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type lang struct {
	Name      string   `json:"Name"`
	Extension string   `json:"extension"`
	Compile   []string `json:"compile"`
}

var langs []lang
var thisLang lang
var workingDir string

var file string
var intermediary int
var debug string
var output string

func _init() {
	//* Get WD
	var err error
	workingDir, err = os.Getwd()
	handleError(err)

	//* Parse flags
	var filePTR *string
	var intermediaryPTR *int
	var debugPTR *string
	var outputPTR *string
	filePTR = flag.String("F", filepath.Join(workingDir, "example/main.pseudo"), "Input file")
	intermediaryPTR = flag.Int("L", 0, "Intermediary language index")
	debugPTR = flag.String("D", "true", "Debug")
	outputPTR = flag.String("O", filepath.Join(workingDir, "example/out.exe"), "Output file")
	flag.Parse()
	file = *filePTR
	intermediary = *intermediaryPTR
	debug = *debugPTR
	output = *outputPTR

	//* Get langs
	sourceFile, err := os.Open("langs.json")
	handleError(err)
	defer sourceFile.Close()
	jsonData, err := io.ReadAll(sourceFile)
	handleError(err)
	handleError(json.Unmarshal([]byte(jsonData), &langs))

	//* Build compile args
	thisLang = langs[intermediary]
	for i := 0; i < len(thisLang.Compile); i++ {
		thisLang.Compile[i] = strings.ReplaceAll(thisLang.Compile[i], "{{ INPUT }}", file+"."+thisLang.Extension)
		thisLang.Compile[i] = strings.ReplaceAll(thisLang.Compile[i], "{{ OUTPUT }}", output)
		thisLang.Compile[i] = strings.ReplaceAll(thisLang.Compile[i], "{{ DEBUG }}", debug)
	}
}

func runScript() string {
	// Create an *exec.Cmd object
	command := exec.Command("python", workingDir+"\\gpt.py")

	// Run the command and capture its output
	output, err := command.CombinedOutput()
	handleError(err)

	return string(output)
}

func cp(f string) {
	// Open the source file for reading
	sourceFile, err := os.Open(f)
	handleError(err)
	defer sourceFile.Close()

	// Copy the contents from the source file to the destination file
	data, err := io.ReadAll(sourceFile)
	err = os.WriteFile(filepath.Join(workingDir, "work.pseudo"), []byte(fmt.Sprintf("# LANG: %s\n", thisLang.Name)+string(data)), os.ModeCharDevice)
	handleError(err)
}

func handleError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func main() {
	_init()
	cp(file)

	//fmt.Println(thisLang)

	// Process Data
	data := runScript()
	data = "Generated by PseudoComp: \n" + data
	fields := strings.Split(data, "```")

	// Get code fields
	var codeFields []string
	for i := 0; i < len(fields); i++ {
		if i%2 == 1 {
			codeFields = append(codeFields, fields[i])
		}
	}
	if len(codeFields) == 0 {
		//handleError(fmt.Errorf("Something went wrong, analyzing your document! Output: \n%s\n", data))
	}

	// Write file
	codeFields[0] = strings.Join(strings.Split(codeFields[0], "\n")[1:], "\n")
	fmt.Println(codeFields[0])
	err := os.WriteFile(file+"."+thisLang.Extension, []byte(codeFields[0]), os.ModeCharDevice)
	handleError(err)

	// Compile
	command := exec.Command(thisLang.Compile[0], thisLang.Compile[1:]...)
	output, err := command.CombinedOutput()
	handleError(err)
	fmt.Println(string(output))
}