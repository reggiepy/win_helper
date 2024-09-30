package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"win_helper/pkg/util/fileUtils"
	"win_helper/templates"
)

type ServerXML struct {
	XMLName xml.Name `xml:"service"`
	Id      string   `xml:"id" json:"id"`
	// Executable
	// Required This element specifies the executable to be launched. It can be either absolute path, or you can just specify the executable name and let it be searched from PATH (although note that the services often run in a different user account and therefore it might have different PATH than your shell does.)
	Executable string `xml:"executable" json:"executable"`
	// Name
	// Optional Short display name of the service, which can contain spaces and other characters. This shouldn't be too long, like <id>, and this also needs to be unique among all the services in a given system.
	Name string `xml:"name,omitempty" json:"name"`
	// Description
	// Optional Long human-readable description of the service. This gets displayed in Windows service manager when the service is selected.
	Description string `xml:"description" json:"description"`
	// StartMode
	// Optional This element specifies the start mode of the Windows service. It can be one of the following values: Automatic, or Manual. For more information, see the ChangeStartMode method. The default value is Automatic.
	// Boot Start ("Boot")
	// Device driver started by the operating system loader. This value is valid only for driver services.
	// System ("System")
	// Device driver started by the operating system initialization process. This value is valid only for driver services.
	// Auto Start ("Automatic")
	// Service to be started automatically by the service control manager during system startup.
	// Demand Start ("Manual")
	// Service to be started by the service control manager when a process calls the StartService method.
	// Disabled ("Disabled")
	// Service that can no longer be started.
	StartMode string `xml:"startmode,omitempty" json:"startmode,omitempty"`
	// Optional Specify IDs of other services that this service depends on. When service X depends on service Y, X can only run if Y is running.
	// Multiple elements can be used to specify multiple dependencies.
	Dependencies []*Dependency `xml:"depend,omitempty" json:"dependencies,omitempty"`
	LogPath      string        `xml:"logpath,omitempty" json:"logpath,omitempty"`
	//Arguments
	//Optional The element specifies the arguments to be passed to the executable.<arguments>
	//<arguments>arg1 arg2 arg3</arguments>
	//-or-
	//<arguments>
	//arg1
	//arg2
	//arg3
	//</arguments>
	Arguments string `xml:"arguments,omitempty" json:"arguments,omitempty"`

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
	// Additional commands
	PreStart  *AdditionalCommands `xml:"prestart,omitempty" json:"prestart,omitempty"`
	PostStart *AdditionalCommands `xml:"poststart,omitempty" json:"poststart,omitempty"`
	PreStop   *AdditionalCommands `xml:"prestop,omitempty" json:"prestop,omitempty"`

	// 关机前
	// 在系统关闭时为服务提供更多停止时间。
	PreShutdown string `xml:"preshutdown,omitempty" json:"preshutdown,omitempty"`
	// 系统默认的关机前超时时间为三分钟。
	PreShutdownTimeout string `xml:"preshutdownTimeout,omitempty" json:"preshutdownTimeout,omitempty"`

	// StopTimeout
	StopTimeout string `xml:"stoptimeout,omitempty" json:"stoptimeout,omitempty"`
	// 环境
	Env *Env `xml:"env,omitempty" json:"env,omitempty"`
	// 哔哔关门
	// 可选元素用于在服务关闭时发出简单的提示音。 此功能应仅用于调试，因为某些操作系统和硬件不支持此功能。
	BeepOnShutdown bool `xml:"beeponshutdown,omitempty" json:"beeponshutdown,omitempty"`
	Log            *Log `xml:"log" json:"log"`
	// OnFailures（失败）
	OnFailures []*OnFailure `xml:"onfailure,omitempty" json:"onfailures,omitempty"`

	WorkingDirectory string `xml:"workingdirectory,omitempty" json:"workingdirectory,omitempty"`
}

type AdditionalCommands struct {
	Executable string `xml:"executable,omitempty" json:"executable,omitempty"`
	Arguments  string `xml:"arguments,omitempty" json:"arguments,omitempty"`
	// stdoutPath specifies the path to redirect the standard output to.
	StdoutPath string `xml:"stdoutPath,omitempty" json:"stdoutPath,omitempty"`
	// stderrPath specifies the path to redirect the standard error output to.
	// Specify in or to dispose of the corresponding stream.NULstdoutPathstderrPath
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
	sForce   bool

	SId               string
	SExecutable       string
	SName             string
	SDescription      string
	SStartMode        string
	SDepends          string
	SLogPath          string
	SArguments        string
	SStartArguments   string
	SStopExecutable   string
	SStopArguments    string
	SEnv              string
	SFailure          string
	SWorkingDirectory string

	SLogMode                string
	SLogPattern             string
	SLogAutoRollAtTime      string
	SLogSizeThreshold       string
	SLogZipOlderThanNumDays string
	SLogZipDateFormat       string
}

func NewDefaultServer() *Server {
	return &Server{
		SLogMode: "roll",
	}
}

func NewServer(opts ...Option) (*Server, error) {
	s := NewDefaultServer()
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Server) GenerateServer() error {
	filename := path.Join(s.BasePath, fmt.Sprintf("%s-server.exe", s.SName))
	if fileUtils.FileExist(filename) {
		if s.sForce {
			err := os.Remove(filename)
			if err != nil {
				return fmt.Errorf("删除服务文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("服务文件 %s 已存在", filename)
		}
	}
	err := os.WriteFile(filename, templates.WinSW, 0o644)
	if err != nil {
		return fmt.Errorf("写入服务失败。%v", err)
	}
	return nil
}

func (s *Server) GenerateServerXML() error {
	serverXML := &ServerXML{
		Id:          s.SName,
		Name:        s.SName,
		Description: s.SDescription,
		Executable:  s.SExecutable,
		Log: &Log{
			Mode: s.SLogMode,
		},
	}

	if s.SArguments != "" {
		serverXML.Arguments = s.SArguments
	}
	if s.SStartArguments != "" {
		serverXML.StartArguments = s.SStartArguments
	}
	if s.SStopExecutable != "" {
		serverXML.StopExecutable = s.SStopExecutable
	}
	if s.SStopArguments != "" {
		serverXML.StopArguments = s.SStopArguments
	}

	if s.SLogPath != "" {
		serverXML.LogPath = s.SLogPath
	}
	switch s.SLogMode {
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
		serverXML.Log.Pattern = s.SLogPattern
		serverXML.Log.SizeThreshold = s.SLogSizeThreshold
		serverXML.Log.AutoRollAtTime = s.SLogAutoRollAtTime
	}

	if s.SStopExecutable != "" {
		serverXML.StopExecutable = s.SStopExecutable
	}

	depends := strings.Split(s.SDepends, ",")
	if len(depends) > 1 || (len(depends) == 1 && depends[0] != "") {
		// 存在依赖项
		serverXML.Dependencies = make([]*Dependency, len(depends))
		for _, d := range depends {
			serverXML.Dependencies = append(serverXML.Dependencies, &Dependency{Value: d})
		}
	}

	if jsonData, err := json.Marshal(serverXML); err == nil {
		fmt.Println(string(jsonData))
	}
	filename := filepath.Join(s.BasePath, fmt.Sprintf("%s-server.xml", s.SName))
	if fileUtils.FileExist(filename) {
		if s.sForce {
			err := os.Remove(filename)
			if err != nil {
				return fmt.Errorf("删除服务文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("服务文件 %s 已存在", filename)
		}
	}
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

func (s *Server) Generate() error {
	var err error
	err = s.GenerateServer()
	if err != nil {
		return err
	}
	err = s.GenerateServerXML()
	if err != nil {
		return err
	}
	return nil
}
