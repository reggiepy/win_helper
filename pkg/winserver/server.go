package winserver

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	Env []*Env `xml:"env,omitempty" json:"env,omitempty"`
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
	SizeThreshold       int    `xml:"sizeThreshold,omitempty" json:"sizeThreshold,omitempty"`
	KeepFiles           int    `xml:"keepFiles,omitempty" json:"keepFiles,omitempty"`
	ZipOlderThanNumDays string `xml:"zipOlderThanNumDays,omitempty" json:"zipOlderThanNumDays,omitempty"`
	ZipDateFormat       string `xml:"zipDateFormat,omitempty" json:"zipDateFormat,omitempty"`
}

type Dependency struct {
	XMLName xml.Name `xml:"depend,omitempty" json:"-"`
	Value   string   `xml:",chardata" json:"value"`
}

func (s *ServerXML) ToJson() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func (s *ServerXML) LoadJson(data string) (*ServerXML, error) {
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		return nil, fmt.Errorf("parse json error: %v", err)
	}
	return s, nil
}

func (s *ServerXML) ToXML() (string, error) {
	var b bytes.Buffer
	encoder := xml.NewEncoder(&b)
	encoder.Indent("", "    ")

	// 处理编码错误，返回错误信息
	if err := encoder.Encode(s); err != nil {
		return "", fmt.Errorf("XML 编码失败: %v", err)
	}

	return b.String(), nil
}

func (s *ServerXML) LoadXML(data string) (*ServerXML, error) {
	err := xml.Unmarshal([]byte(data), s)
	if err != nil {
		return nil, fmt.Errorf("parse json error: %v", err)
	}
	return s, nil
}

type Server struct {
	BasePath string
	sForce   bool

	SId               string
	SExecutable       string
	SName             string
	SDescription      string
	SStartMode        string
	SDepends          []string
	SLogPath          string
	SArguments        string
	SStartArguments   string
	SStopExecutable   string
	SStopArguments    string
	SEnv              []string
	SFailure          string
	SWorkingDirectory string

	SLogMode           string
	SLogPattern        string
	SLogAutoRollAtTime string
	SLogSizeThreshold  int
	SLogKeepFiles      int
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
	if fsutil.FileExist(filename) {
		if s.sForce {
			err := os.Remove(filename)
			if err != nil {
				return fmt.Errorf("删除服务文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("服务文件 %s 已存在", filename)
		}
	}
	err := os.WriteFile(filename, WinSW, 0o644)
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

	// 简化参数检查逻辑
	serverXML.Arguments = s.SArguments
	serverXML.StartArguments = s.SStartArguments
	serverXML.StopExecutable = s.SStopExecutable
	serverXML.StopArguments = s.SStopArguments
	serverXML.LogPath = s.SLogPath

	// 使用字典处理模式匹配
	logModeMap := map[string]func(){
		"append": func() { serverXML.Log.Mode = "append" },
		"reset":  func() { serverXML.Log.Mode = "reset" },
		"none":   func() { serverXML.Log.Mode = "none" },
		"roll-by-size": func() {
			serverXML.Log.Mode = "roll-by-size"
			serverXML.Log.SizeThreshold = s.SLogSizeThreshold
			serverXML.Log.KeepFiles = s.SLogKeepFiles
		},
		"roll-by-time": func() {
			serverXML.Log.Mode = "roll-by-time"
			serverXML.Log.Pattern = s.SLogPattern
		},
		"roll-by-size-time": func() {
			serverXML.Log.Mode = "roll-by-time"
			serverXML.Log.Pattern = s.SLogPattern
			serverXML.Log.SizeThreshold = s.SLogSizeThreshold
			serverXML.Log.KeepFiles = s.SLogKeepFiles
			serverXML.Log.AutoRollAtTime = s.SLogAutoRollAtTime
		},
	}

	// 根据模式选择对应的处理逻辑
	if setLogMode, exists := logModeMap[s.SLogMode]; exists {
		setLogMode()
	}

	// 处理依赖项
	for _, d := range s.SDepends {
		if d != "" {
			serverXML.Dependencies = append(serverXML.Dependencies, &Dependency{Value: d})
		}
	}

	for _, e := range s.SEnv {
		if e != "" {
			eSplit := strings.SplitN(e, "=", 2)
			if len(eSplit) != 2 {
				continue
			}
			k, v := eSplit[0], eSplit[1]
			serverXML.Env = append(serverXML.Env, &Env{k, v})
		}
	}

	// 输出调试信息（可选）
	fmt.Println(serverXML.ToJson())

	// 生成文件路径
	filename := filepath.Join(s.BasePath, fmt.Sprintf("%s-server.xml", s.SName))

	// 检查文件是否存在并根据强制标志删除
	if fsutil.FileExist(filename) {
		if s.sForce {
			if err := os.Remove(filename); err != nil {
				return fmt.Errorf("删除服务文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("服务文件 %s 已存在", filename)
		}
	}

	// 创建并写入 XML 文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 编码 XML 并写入文件
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
