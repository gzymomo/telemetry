package main
 
import (
    "net"
 
 
 
    "fmt"
    "log"
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/gin-gonic/gin"
    consulapi "github.com/hashicorp/consul/api"
)
 
const (
    consulAddress = "127.0.0.1:8500"
    serviceId     = "111"
)
 
func main() {
    r := gin.Default()
 
    // consul健康检查回调函数
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "ok",
        })
    })
    r.GET("/metrics",  func(c *gin.Context){
        promhttp.Handler().ServeHTTP(c.Writer, c.Request)
    })
    go http.ListenAndServe(":8081", r)
    // 注册服务到consul
    ConsulRegister()
 
    // 从consul中发现服务
    ConsulFindServer()
 
    ConsulCheckHeath()
    ConsulKVTest()
    // 取消consul注册的服务
    //ConsulDeRegister()
    var str string
    fmt.Scan(&str)
 
}


// 注册服务到consul
func ConsulRegister() {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = consulAddress
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 创建注册到consul的服务到
    registration := new(consulapi.AgentServiceRegistration)
    registration.ID = serviceId                    // 服务节点的名称
    registration.Name = "go-consul-test"           // 服务名称
    registration.Port = 8081                       // 服务端口
    registration.Tags = []string{"web-gin"} // tag，可以为空
    registration.Address = "192.168.124.16"          // 服务 IP 要确保consul可以访问这个ip
    registration.Meta = map[string]string{
        "team":"telemetry",
        "suborgin":"grafana",
        "job":"web",
    }
 
    // 增加consul健康检查回调函数
    check := new(consulapi.AgentServiceCheck)
    check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
    check.Timeout = "5s"
    check.Interval = "5s"                        // 健康检查间隔
    check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
    registration.Check = check
 
    // 注册服务到consul
    err = client.Agent().ServiceRegister(registration)
    if err == nil {
        fmt.Println("ConsulRegister done")
    }
}
 
// 取消consul注册的服务
func ConsulDeRegister() {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = consulAddress
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    client.Agent().ServiceDeregister(serviceId)
}
 
// 从consul中发现服务
func ConsulFindServer() {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = consulAddress
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 获取所有service
    services, _ := client.Agent().Services()
    for _, value := range services {
        fmt.Println("address:", value.Address)
        fmt.Println("port:", value.Port)
    }
 
    fmt.Println("=================================")
    // 获取指定service
    service, _, err := client.Agent().Service(serviceId, nil)
    if err == nil {
        fmt.Println("address:", service.Address)
        fmt.Println("port:", service.Port)
    }
    if err == nil {
        fmt.Println("ConsulFindServer done")
    }
}
 
func ConsulCheckHeath() {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = consulAddress
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // 健康检查
    a, b, _ := client.Agent().AgentHealthServiceByID(serviceId)
    fmt.Println("val1:", a)
    fmt.Println("val2:", b)
    fmt.Println("ConsulCheckHeath done")
}
 
func ConsulKVTest() {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = consulAddress
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }
 
    // KV, put值
    values := "test"
    key := "go-consul-test"
    client.KV().Put(&consulapi.KVPair{Key: key, Flags: 0, Value: []byte(values)}, nil)
 
    // KV get值
    data, _, _ := client.KV().Get(key, nil)
    fmt.Println("data:", string(data.Value))
 
    // KV list
    datas, _, _ := client.KV().List("go", nil)
    for _, value := range datas {
        fmt.Println("val:", value)
    }
    keys, _, _ := client.KV().Keys("go", "", nil)
    fmt.Println("key:", keys)
    fmt.Println("ConsulKVTest done")
}
 
func localIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}