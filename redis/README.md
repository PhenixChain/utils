# 使用说明：

## config.ini 配置
```
[redis]
; 地址
host                = 192.168.0.213:16379
; 密码
pwd                 = 123456
; 数据库序号
db_index            = 0
; 最大空闲连接数
max_idle            = 1
; 最大连接数
max_active          = 100
; 空闲连接超时时间(秒)
idle_timeout        = 10
; 连接超时时间(秒)
conn_timeout        = 30
```

## API使用
```
redisPool := InitRedisPool()

//String
_, err := redisPool.Set("key", "123")
if err != nil {
	log.Fatalln(err)
}
val, err := redisPool.Get("key")
if err != nil {
	log.Fatalln(err)
}
fmt.Println("key:", val)

//Hash
hash := map[string]string{
	"field1": "value1",
	"field2": "value2",
	"field3": "value3",
}
_, err = redisPool.Hmset("Hkey", hash)
if err != nil {
	log.Fatalln(err)
}
hVal, err := redisPool.Hmget("Hkey", []string{"field1", "field2"})
if err != nil {
	log.Fatalln(err)
}
for _, item := range hVal {
	fmt.Println(item)
}
```    

命令参考: [redisfans](http://doc.redisfans.com/index.html)