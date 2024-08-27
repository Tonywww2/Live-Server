package api

import (
	"github.com/gin-gonic/gin"
	"live_server/config"
	"log"
	"strconv"
)

type IApi interface {
	RegisterRouter(server *GinServer)
}

// 参数是Engine和路由分组
type GinServer struct {
	//内联的方式扩展Engine
	*gin.Engine
	// 路由分组，多个不同功能api会用到
	Rg *gin.RouterGroup
}

func Initialize() *GinServer {
	g := &GinServer{Engine: gin.Default()}
	// 还有gin相关都放到这里来
	return g
}

// 服务端口，一般都需要读取配置文件
func (g *GinServer) Listen() {
	err := g.Engine.Run(":" + strconv.Itoa(config.Config.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (g *GinServer) RegisterRouters(apis ...IApi) *GinServer {
	// 遍历所有的控制层，这里使用接口，就是为了将Router实例化
	for _, c := range apis {
		c.RegisterRouter(g)
	}
	return g
}
