package win

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"win_helper/pkg/util/fileUtils"
)

//go:embed ..\..\..\bin\winsw\WinSW-x64.exe
var server []byte

type ServerXML struct {
	XMLName     xml.Name `xml:"service"`
	Id          string   `xml:"id" json:"id"`
	Name        string   `xml:"name,omitempty" json:"name,omitempty"`
	Description string   `xml:"description,omitempty" json:"description,omitempty"`
	LogPath     string   `xml:"logpath,omitempty" json:"logpath,omitempty"`
	//环境
	Env          *Env          `xml:"env,omitempty" json:"env,omitempty"`
	Log          *Log          `xml:"log" json:"log"`
	Dependencies []*Dependency `xml:"depend,omitempty" json:"dependencies,omitempty"`
	Executable   string        `xml:"executable,omitempty" json:"executable,omitempty"`
	//Arguments
	//Optional The element specifies the arguments to be passed to the executable.<arguments>
	//<arguments>arg1 arg2 arg3</arguments>
	//-or-
	//<arguments>
	//arg1
	//arg2
	//arg3
	//</arguments>
	Arguments        string `xml:"arguments,omitempty" json:"arguments,omitempty"`
	WorkingDirectory string `xml:"workingdirectory,omitempty" json:"workingdirectory,omitempty"`

	//stopargument/stopexecutable
	//Optional When the service is requested to stop, winsw simply calls TerminateProcess function to kill the service instantly.
	//However, if the element is present, winsw will instead launch another process of (or if that's specified) with the
	//specified arguments, and expects that to initiate the graceful shutdown of the service process.<stoparguments><executable><stopexecutable>
	//Winsw will then wait for the two processes to exit on its own, before reporting back to Windows that the service has terminated.
	//When you use the , you must use instead of . See the complete example below:<stoparguments><startarguments><arguments>
	//<executable>catalina.sh</executable>
	//<startarguments>jpda run</startarguments>
	//<stopexecutable>catalina.sh</stopexecutable>
	//<stoparguments>stop</stoparguments>
	StartArguments string `xml:"startarguments,omitempty" json:"startarguments,omitempty"`
	StopExecutable string `xml:"stopexecutable,omitempty" json:"stopexecutable,omitempty"`
	StopArguments  string `xml:"stoparguments,omitempty" json:"stoparguments,omitempty"`

	//关机前
	//在系统关闭时为服务提供更多停止时间。
	PreShutdown string `xml:"preshutdown,omitempty" json:"preshutdown,omitempty"`
	//系统默认的关机前超时时间为三分钟。
	PreShutdownTimeout string `xml:"preshutdownTimeout,omitempty" json:"preshutdownTimeout,omitempty"`

	//StopTimeout
	StopTimeout string `xml:"stoptimeout,omitempty" json:"stoptimeout,omitempty"`

	//哔哔关门
	//可选元素用于在服务关闭时发出简单的提示音。 此功能应仅用于调试，因为某些操作系统和硬件不支持此功能。
	BeepOnShutdown bool `xml:"beeponshutdown,omitempty" json:"beeponshutdown,omitempty"`

	// OnFailures（失败）
	OnFailures []*OnFailure `xml:"onfailure,omitempty" json:"onfailures,omitempty"`

	//Additional commands
	PreStart  *AdditionalCommands `xml:"prestart,omitempty" json:"prestart,omitempty"`
	PostStart *AdditionalCommands `xml:"poststart,omitempty" json:"poststart,omitempty"`
	PreStop   *AdditionalCommands `xml:"prestop,omitempty" json:"prestop,omitempty"`
	PostStop  *AdditionalCommands `xml:"poststop,omitempty" json:"poststop,omitempty"`
}

type AdditionalCommands struct {
	Executable string `xml:"executable,omitempty" json:"executable,omitempty"`
	Arguments  string `xml:"arguments,omitempty" json:"arguments,omitempty"`
	//stdoutPath specifies the path to redirect the standard output to.
	StdoutPath string `xml:"stdoutPath,omitempty" json:"stdoutPath,omitempty"`
	//stderrPath specifies the path to redirect the standard error output to.
	//Specify in or to dispose of the corresponding stream.NULstdoutPathstderrPath
	StderrPath string `xml:"stderrPath,omitempty" json:"stderrPath,omitempty"`
}

type OnFailure struct {
	Action string `xml:"action,omitempty" json:"action,omitempty"`
	Delay  string `xml:"delay,omitempty" json:"delay,omitempty"`
}
type Env struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:"value,attr" json:"value"`
}

type Log struct {
	Mode                string `xml:"mode,attr" json:"mode"`
	Pattern             string `xml:"pattern,omitempty" json:"pattern,omitempty"`
	AutoRollAtTime      string `xml:"autoRollAtTime,omitempty" json:"autoRollAtTime,omitempty"`
	SizeThreshold       string `xml:"sizeThreshold,omitempty" json:"sizeThreshold,omitempty"`
	ZipOlderThanNumDays string `xml:"zipOlderThanNumDays,omitempty" json:"zipOlderThanNumDays,omitempty"`
	ZipDateFormat       string `xml:"zipDateFormat,omitempty" json:"zipDateFormat,omitempty"`
}

type Dependency struct {
	XMLName xml.Name `xml:"depend,omitempty" json:"-"`
	Value   string   `xml:",chardata" json:"value"`
}

type Server struct {
	BasePath string

	Name                   string
	Description            string
	Executable             string
	Arguments              string
	StopExecutable         string
	LogMode                string
	LogPattern             string
	LogAutoRollAtTime      string
	LogSizeThreshold       string
	LogZipOlderThanNumDays string
	LogZipDateFormat       string
	StartMode              string
	Depends                string
}

type Option func(s *Server) error

func WithExecutable(exec string) Option {
	return func(s *Server) error {
		s.Executable = exec
		return nil
	}
}

func WithBasePath(basePath string) Option {
	return func(s *Server) error {
		s.BasePath = basePath
		return nil
	}
}

func WithName(name string) Option {
	return func(s *Server) error {
		s.Name = name
		return nil
	}
}

func WithLogMode(logMode string) Option {
	return func(s *Server) error {
		s.LogMode = logMode
		return nil
	}
}

func NewDefaultServer() *Server {
	return &Server{
		LogMode: "roll",
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			panic(err.(any))
		}
	}
	return s
}

func (s *Server) CreateServer() error {
	fileName := path.Join(s.BasePath, fmt.Sprintf("%s-server.exe", s.Executable))
	if fileUtils.FileExist(fileName) {
		return fmt.Errorf("%s already exists", fileName)
	}
	err := os.WriteFile(fileName, server, 0644)
	if err != nil {
		return fmt.Errorf("写入服务失败。%v", err)
	}
	return nil
}

func (s *Server) CreateServerXML() error {
	serverXML := &ServerXML{
		Id:         s.Name,
		Executable: s.Executable,
		Log: &Log{
			Mode: s.LogMode,
		},
	}
	switch s.LogMode {
	case "append":
		serverXML.Log.Mode = "append"
	case "reset":
		serverXML.Log.Mode = "reset"
	case "none":
		serverXML.Log.Mode = "none"
	case "roll":
		serverXML.Log.Mode = "roll"
	case "roll-by-size":
		serverXML.Log.Mode = "roll-by-size"

	case "roll-by-time":
		serverXML.Log.Mode = "roll-by-time"
		serverXML.Log.Pattern = s.LogPattern
		serverXML.Log.SizeThreshold = s.LogSizeThreshold
		serverXML.Log.AutoRollAtTime = s.LogAutoRollAtTime
	}

	if s.StopExecutable != "" {
		serverXML.StopExecutable = s.StopExecutable
	}

	depends := strings.Split(s.Depends, ",")
	if len(depends) > 1 || (len(depends) == 1 && depends[0] != "") {
		// 存在依赖项
		serverXML.Dependencies = make([]*Dependency, len(depends))
		for _, d := range depends {
			serverXML.Dependencies = append(serverXML.Dependencies, &Dependency{Value: d})
		}
	}

	fmt.Println(serverXML)
	filename := filepath.Join(s.BasePath, fmt.Sprintf("%s-server.xml", s.Executable))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "    ")
	if err := encoder.Encode(serverXML); err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() error {
	var err error
	//err = s.CreateServer()
	//if err != nil {
	//	return err
	//}
	err = s.CreateServerXML()
	if err != nil {
		return err
	}
	return nil
}
