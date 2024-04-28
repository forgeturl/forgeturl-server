package conf

var LocalConfig = `
[base]
debug = true
logLevel = "info" #线上环境需要设置为 error

[dbDemo]
Host = ""           #数据库连接地址
Name = "demo"           #数据库名称
User = ""           #数据库用户名
Passwd = ""         #数据库密码
Port = "3306"       #数据库端口号
MaxIdleConns = 20
MaxOpenConns = 20

[demo]
host = "" 
port = ":6379"
auth = ""
prefix = ""

[demo.redisDB]
db0 = 0

			
[aliyunmq]
NameServer = ""
AccessKey = ""
SecretKey = ""
Namespace = ""

			
`
