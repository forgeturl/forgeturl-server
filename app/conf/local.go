package conf

var LocalConfig = `
[base]
debug = true
logLevel = "info" # when online can set to error
domain = "localhost:8080"

[mysql]
  Host = "127.0.0.1" 
  Port = "3306"
  Name = "forget_url"           #database name
  User = "root"
  Passwd = "pwdTest"
  MaxOpenConns = 20 
  MaxIdleConns = 20
  ConnMaxLifeTime = 6000
  Type = "mysql"
  Debug = true

[redisServer]
  host = "127.0.0.1"
  port = ":6379"
  auth = ""
  encryption = 0

`

var OnlConfig = `
[base]
debug = true
logLevel = "info"
domain = "api.forgeturl.com"

[mysql]
  Host = "127.0.0.1" 
  Port = "3306"
  Name = "forget_url"           #database name
  User = "root"
  Passwd = ""
  MaxOpenConns = 20 
  MaxIdleConns = 20
  ConnMaxLifeTime = 6000
  Type = "mysql"
  Debug = true

[redisServer]
  host = "127.0.0.1"
  port = ":6379"
  auth = ""
  encryption = 0

[redisDB]
  db0 = 0

`
