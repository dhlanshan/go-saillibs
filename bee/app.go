package bee

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// MagicApp 魔术师
type MagicApp struct {
	Addr        string // 运行地址端口
	ConfPath    string // 配置文件路径
	ConfName    string // 配置文件名
	IsDefault   bool   // 是否使用默认路由引擎
	RunMode     string // 运行模式
	RegRouteFun func(r *gin.Engine)
	router      *gin.Engine
	isInit      bool // 是否初始化过
}

// Init 初始化
func (m *MagicApp) Init() {
	// 初始化路由引擎
	m.initRouter()
	// 初始化配置文件
	m.initConfig()
	// 初始化日志
	m.initLog()
	// 心跳检测
	m.testRoute()
	// 路由注册
	if m.RegRouteFun != nil {
		m.RegRouteFun(m.router)
	}
	// 初始化完毕
	m.isInit = true

}

func (m *MagicApp) Run() {
	if !m.isInit {
		fmt.Printf("[%s] server start err: 未初始化数据\n", time.Now().Format(time.DateTime))
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{Addr: fmt.Sprintf("%s", m.Addr), Handler: m.router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("[%s] server listen err: %s\n", time.Now().Format(time.DateTime), err)
		}
	}()
	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown: ", err)
	}
	fmt.Println("Server exiting")
}

// initRouter 初始化路由引擎
func (m *MagicApp) initRouter() {
	if m.RunMode == "" {
		m.RunMode = "release"
	}
	gin.SetMode(m.RunMode)
	if m.IsDefault {
		m.router = gin.Default()
	} else {
		m.router = gin.New()
	}
}

// initConfig 初始化配置
func (m *MagicApp) initConfig() {
	if m.ConfPath == "" || m.ConfName == "" {
		fmt.Printf("[%s] 初始化配置...跳过\n", time.Now().Format(time.DateTime))
		return
	}
	viper.SetConfigName(m.ConfName)
	viper.AddConfigPath(m.ConfPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("初始化配置失败")
	}
	fmt.Printf("[%s] 初始化配置...ok\n", time.Now().Format(time.DateTime))
}

// InitLog 初始化日志
func (m *MagicApp) initLog() {
	magicLog := &MagicLog{LogKey: "app.log"}
	magicLog.initLogger()
	fmt.Printf("[%s] 初始化日志...ok\n", time.Now().Format(time.DateTime))
}

// testRoute 心跳检测
func (m *MagicApp) testRoute() {
	m.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "success! This service is normal.")
	})
}
