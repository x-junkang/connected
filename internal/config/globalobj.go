package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
		Host:             "0.0.0.0",
		TCPPort:          8090,
		HttpPort:         8080,
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     pwd + "/conf/connected.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "tmp.log",
		LogLevel:         "debug",
	}

	GlobalObject.Reload()
}

type GlobalObj struct {
	Host     string `json:"host"`      //当前服务器主机IP
	TCPPort  int    `json:"tcp_port"`  //当前服务器主机监听端口号
	HttpPort int    `json:"http_port"` //http监听端口
	Name     string //当前服务器名称

	Version          string `json:"name"`                //当前Zinx版本号
	MaxPacketSize    uint32 `json:"max_packer_size"`     //都需数据包的最大值
	MaxConn          int    `json:"max_conn"`            //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 `json:"worker_pool_size"`    //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 `json:"max_worker_task_len"` //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 `json:"max_msg_chan_len"`    //SendBuffMsg发送消息的缓冲最大长度

	ConfFilePath string `json:"conf_file_path"`
	LogDir       string `json:"log_dir"`  //日志所在文件夹 默认"./log"
	LogFile      string `json:"log_file"` //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogLevel     string `json:"log_level"`
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
		fmt.Println("Config File ", g.ConfFilePath, " is not exist!!")
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
	fmt.Println(g)
}
