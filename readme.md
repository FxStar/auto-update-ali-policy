# 介绍
阿里云安全组规则自动更新

# 配置文件
将config_example.json改为 config.json 并修改里面的内容
```
{
    "regionId":"",          // 安全组所属地域ID ，比如 `cn-guangzhou`
    "accessKeyId":"",       // AccessKeyId
    "accessKeySecret":"", // AccessKeySecret
    "securityGroupId":"", // 安全组ID
    "ipProtocol":"tcp", // 协议,可选 tcp,udp, icmp, gre, all：支持所有协议
    "portRange":"22/22", // 端口范围，使用斜线（/）隔开起始端口和终止端口
    "priority":"1", // 安全组规则优先级，数字越小，代表优先级越高。取值范围：1~100
    "policy":"accept", // accept:接受访问, drop: 拒绝访问
    "lastIp":"" // 最后配置安全组的ip(如果首次配置为空则会立即添加当前ip安全组)
}

```

    // RegionId：安全组所属地域ID ，比如 `cn-guangzhou`
	// 访问 [DescribeRegions:查询可以使用的阿里云地域](https://next.api.aliyun.com/api/Ecs/2014-05-26/DescribeRegions) 查阅
	// 国内一般是去掉 ECS 所在可用区的后缀，比如去掉 cn-guangzhou-b 的尾号 -b

    // 前4项配置必填，其他为config_example 的默认值，可自行修改

# 部署运行
 可直接自行编译运行
 也可自己打包docker镜像运行

 也可使用已经打包好的docker镜像运行
 docker run -d --name auto_update_ali_policy -v /yourPath/config:/app/config fx0408/auto-update-ali-policy:latest
 