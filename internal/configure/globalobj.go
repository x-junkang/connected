package configure

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/x-junkang/connected/internal/clog"
)

var GlobalObject *GlobalObj

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	//default
	GlobalObject = &GlobalObj{
		Name:             "Connected",
		Version:          "V1.0",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     pwd + "/conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "",
		LogLevel:         "debug",
	}

	GlobalObject.Reload()
}

type GlobalObj struct {
	Host    string //当前服务器主机IP
	TCPPort int    //当前服务器主机监听端口号
	Name    string //当前服务器名称

	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	/*
		config file path
	*/
	ConfFilePath string

	/*
		logger
	*/
	LogDir   string //日志所在文件夹 默认"./log"
	LogFile  string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogLevel string
}

//PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {

	if confFileExists, _ := PathExists(g.ConfFilePath); !confFileExists {
		//fmt.Println("Config File ", g.ConfFilePath , " is not exist!!")
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

	//Logger 设置
	clog.InitLogger(g.LogFile, "")
}
