package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func main() {
	c, err := initConfig()
	if err != nil {
		fmt.Printf("init config failed")
		panic(err)
	}
	// 每3分钟循环一次判断ip
	for {
		if ok, newIp := CheckIP(c.LastIp); !ok {
			UpdateIpPolice(c, newIp)

		}
		fmt.Printf("sleep 3 min\n")
		time.Sleep(180 * time.Second)

	}
}

// {"result":true,"code":"querySuccess","message":"Query Success","IP":"1.193.37.211","IPVersion":"IPv4"}
type Response struct {
	Result    bool   `json:"result"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	IP        string `json:"IP"`
	IPVersion string `json:"IPVersion"`
}

// config
//
//	{
//	    "regionId":"",
//	    "accessKeyId":"",
//	    "accessKeySecret":"",
//	    "securityGroupId":"",
//	    "ipProtocol":"",
//	    "portRange":"",
//	    "priority":"",
//	    "policy":"",
//	    "lastIp":""
//	}
type Config struct {
	RegionId        string `json:"regionId"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityGroupId string `json:"securityGroupId"`
	IpProtocol      string `json:"ipProtocol"`
	PortRange       string `json:"portRange"`
	Priority        string `json:"priority"`
	Policy          string `json:"policy"`
	LastIp          string `json:"lastIp"`
}

func initConfig() (*Config, error) {
	c := Config{}
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonFile), &c)
	if err != nil {
		return nil, err
	}
	if c.RegionId == "" || c.AccessKeyId == "" || c.AccessKeySecret == "" || c.SecurityGroupId == "" {
		return nil, errors.New("config.json 配置错误")
	}

	if c.IpProtocol == "" {
		c.IpProtocol = "tcp"
	}
	if c.PortRange == "" {
		c.PortRange = "22/22"
	}
	if c.Priority == "" {
		c.Priority = "1"
	}
	if c.Policy == "" {
		c.Policy = "accept"
	}
	// 判断必填项
	return &c, err
}

func updateConfig(config *Config) {
	// 更新 写入config.json
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal config:", err)
		return
	}

	// 3. 将 JSON 写入 config.json 文件
	err = os.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		fmt.Println("Failed to write config file:", err)
		return
	}
}

func UpdateIpPolice(config *Config, clientIP string) {

	// {"result":true,"code":"querySuccess","message":"Query Success","IP":"1.193.37.211","IPVersion":"IPv4"}
	fmt.Printf("Client IP: %s\n", clientIP)
	// <accessKeyId>, <accessSecret>: 前往 https://ram.console.aliyun.com/manage/ak 添加 accessKey
	// RegionId：安全组所属地域ID ，比如 `cn-guangzhou`
	// 访问 [DescribeRegions:查询可以使用的阿里云地域](https://next.api.aliyun.com/api/Ecs/2014-05-26/DescribeRegions) 查阅
	// 国内一般是去掉 ECS 所在可用区的后缀，比如去掉 cn-guangzhou-b 的尾号 -b
	client, err := ecs.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Println("ali client init failed")
		return
	}
	request := ecs.CreateAuthorizeSecurityGroupRequest()
	request.Scheme = "https"
	request.SecurityGroupId = config.SecurityGroupId // 安全组ID
	request.IpProtocol = config.IpProtocol           // "tcp"                       // 协议,可选 tcp,udp, icmp, gre, all：支持所有协议
	request.PortRange = config.PortRange             // "22/22"                        // 端口范围，使用斜线（/）隔开起始端口和终止端口
	request.Priority = config.Priority               //"1"                           // 安全组规则优先级，数字越小，代表优先级越高。取值范围：1~100
	request.Policy = config.Policy                   //"accept"                           // accept:接受访问, drop: 拒绝访问
	request.NicType = "internet"                     // internet：公网网卡, intranet：内网网卡。
	request.SourceCidrIp = clientIP                  // 源端IPv4 CIDR地址段。支持CIDR格式和IPv4格式的IP地址范围。

	response, err := client.AuthorizeSecurityGroup(request)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Println("ali client AuthorizeSecurityGroup failed")
		return
	}
	fmt.Printf("Response: %#v\nClient IP: %s  was successfully added to the Security Group.\n", response, clientIP)

	config.LastIp = clientIP

	updateConfig(config)

}

func CheckIP(oldIp string) (bool, string) {
	responseClient, errClient := http.Get("https://4.ipw.cn/api/ip/myip") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		panic(errClient)
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()
	body, _ := ioutil.ReadAll(responseClient.Body)
	myresponse := Response{}
	json.Unmarshal(body, &myresponse)
	clientIP := myresponse.IP
	if oldIp == clientIP {
		return true, clientIP
	} else {
		return false, clientIP
	}

}
