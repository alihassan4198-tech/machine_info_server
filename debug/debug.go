package debug

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

type ENV int

const (
	B4dmain ENV = iota
	Hivemain
	B4dsandbox
	Local
)

// Debug
func Debug() bool {
	debug := false
	return debug
}

func GetENV() ENV {
	return Hivemain
}

// Trace Functions
func Trace_enter() {
	myStack := string(debug.Stack())
	var goRoutineName string = myStack[:strings.IndexByte(myStack, ':')]
	lines := linesStringCount(myStack)
	var routinueLine int = 1
	var perFuncLines int = 2
	var extraFunc int = 2
	// var minusMainFunc int = 1
	var tabs int = ((lines - routinueLine) / perFuncLines) - extraFunc //- minusMainFunc
	var tabsStr string = ""
	for i := 1; i < tabs; i++ {
		tabsStr += "\t"
	}
	pc, file, lineNum, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	fnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")
	var logStr string = ""
	logStr += " --> " + time.Now().Format("15:04:05")
	logStr += "\t" + goRoutineName
	logStr += "\t" + tabsStr
	logStr += " --> " + fnName
	logStr += "\tfile name: " + file + ","
	logStr += " line num: " + fmt.Sprint(lineNum)
	logStr += "\n"
	fmt.Print(logStr)
}

func GoRoutineName() string {
	myStack := string(debug.Stack())
	var goRoutineName string = myStack[:strings.IndexByte(myStack, ':')]

	return goRoutineName
}

func Trace_exit() {
	myStack := string(debug.Stack())
	lines := linesStringCount(myStack)
	var goRoutineName string = myStack[:strings.IndexByte(myStack, ':')]
	var routinueLine int = 1
	var perFuncLines int = 2
	var extraFunc int = 2
	// var minusMainFunc int = 1
	var tabs int = ((lines - routinueLine) / perFuncLines) - extraFunc //- minusMainFunc
	var tabsStr string = ""
	for i := 1; i < tabs; i++ {
		tabsStr += "\t"
	}
	pc, file, lineNum, _ := runtime.Caller(1)
	_, file = filepath.Split(file)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	fnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")
	var logStr string = ""
	logStr += " <-- " + time.Now().Format("15:04:05")
	logStr += "\t" + goRoutineName
	logStr += "\t" + tabsStr
	logStr += " <-- " + fnName
	logStr += "\tfile name: " + file + ","
	logStr += " line num: " + fmt.Sprint(lineNum)
	logStr += "\n"
	fmt.Print(logStr)
}

func GetFunctionName(vars ...interface{}) string {
	pc, _, _, _ := runtime.Caller(1)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	fnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")
	return fnName
}

// formatRequest generates ascii representation of a request
func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// formatResponse generates ascii representation of a request
func FormatResponse(r *http.Response) string {
	if r == nil {
		return "nil response pointer"
	}
	// Create return string
	var request []string
	// Add the request string
	// Add the host
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	b, err := io.ReadAll(r.Body)
	// b, err := ioutil.ReadAll(resp.Body) //Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	request = append(request, string(b))

	request = append(request, r.Status)

	// Return the request as a string
	return strings.Join(request, "\n")
}

func PrintVariables(vars ...interface{}) {
	return
	myStack := string(debug.Stack())
	pc, _, lineNum, _ := runtime.Caller(1)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	fnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")
	lines := linesStringCount(myStack)
	var goRoutineName string = myStack[:strings.IndexByte(myStack, ':')]
	var routinueLine int = 1
	var perFuncLines int = 2
	var extraFunc int = 2
	var minusMainFunc int = 1
	var tabs int = ((lines - routinueLine) / perFuncLines) - extraFunc - minusMainFunc
	var tabsStr string = ""
	for i := 1; i < tabs+2; i++ {
		tabsStr += "\t"
	}
	var funcAndLineStr = ""
	funcAndLineStr += " func name: " + fnName + ","
	funcAndLineStr += " line number: " + fmt.Sprint(lineNum)
	fmt.Printf(".\t\t"+goRoutineName+" %s>>>>>>>>>>>>>>>\n", tabsStr)
	fmt.Printf(".\t\t"+goRoutineName+" %s%s\n", tabsStr, funcAndLineStr)
	for i, variable := range vars {
		var logStr = "."
		logStr += "\t\t" + goRoutineName
		logStr += " " + tabsStr
		logStr += " " + fmt.Sprint(i) + " ) "
		logStr += fmt.Sprintf("%#+v", variable)
		logStr += "\n"
		fmt.Print(logStr)
	}
	fmt.Printf(".\t\t"+goRoutineName+" %s<<<<<<<<<<<<<<<\n", tabsStr)
}

// DEMO Line
// defer app.TimeTrack(time.Now(), app.FileFunctionLine())
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("time tracking: %s took %s", name, elapsed)
}

func FileFunctionLine() string {
	myStack := string(debug.Stack())
	pc, file, lineNum, _ := runtime.Caller(1)
	// _, file = filepath.Split(file)
	file = strings.Replace(file, "/go/src/github.com/Infinite-Compute/", "", -1)
	functionObject := runtime.FuncForPC(pc)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	fnName := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	pc2, _, _, _ := runtime.Caller(2)
	functionObject2 := runtime.FuncForPC(pc2)
	extractFnName2 := regexp.MustCompile(`^.*\.(.*)$`)
	fnName2 := extractFnName2.ReplaceAllString(functionObject2.Name(), "$1")

	pc3, _, _, _ := runtime.Caller(3)
	functionObject3 := runtime.FuncForPC(pc3)
	extractFnName3 := regexp.MustCompile(`^.*\.(.*)$`)
	fnName3 := extractFnName3.ReplaceAllString(functionObject3.Name(), "$1")

	var goRoutineName string = myStack[:strings.IndexByte(myStack, ':')]
	return fmt.Sprintf("theard:%s , file:%s , func:%s -> %s -> %s , line:%d", goRoutineName, file, fnName3, fnName2, fnName, lineNum)
}

func linesStringCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
