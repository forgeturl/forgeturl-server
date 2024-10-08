package conf

var LocalConfig = `
[base]
debug = true
logLevel = "info" #线上环境需要设置为 error

[sqliteDB]
Path = ""../../identifier.sqlite""
MaxIdleConns = 20
MaxOpenConns = 20

[demo.redisDB]
db0 = 0
`

var OnlConfig = `
[base]
debug = true
logLevel = "info" #线上环境需要设置为 error

[sqliteDB]
Path = ""/data/forgeturl/forgeturl.sqlite""
MaxIdleConns = 20
MaxOpenConns = 20

[demo.redisDB]
db0 = 0
`
