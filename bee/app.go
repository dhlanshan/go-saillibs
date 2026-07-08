package bee

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// MagicService 单个HTTP服务实例
type MagicService struct {
	Name        string              // 服务名(用于日志展示)
	Addr        string              // 运行地址端口
	IsDefault   bool                // 是否使用默认路由引擎
	IsHeartbeat bool                // 开启心跳检测, 默认关闭
	RegRouteFun func(r *gin.Engine) // 路由注册
	Router      *gin.Engine         // Gin路由引擎(可外部注入)

	// 可选：HTTP Server 超时参数（不填则使用 net/http 默认值）
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// MagicApp 魔术
type MagicApp struct {
	ConfPath       string          // 配置文件路径
	ConfName       string          // 配置文件名
	ConfHotLoading bool            // 配置启用热加载
	RunMode        string          // 运行模式
	ExitAfter      func()          // 程序结束后的操作
	Services       []*MagicService // 多服务(如果设置则启动多个HTTP服务)
	isInit         bool            // 是否初始化过
}

// AddService 添加一个HTTP服务
func (m *MagicApp) AddService(s *MagicService) {
	if s == nil {
		return
	}
	m.Services = append(m.Services, s)
}

// Init 初始化
func (m *MagicApp) Init() {
	if len(m.Services) == 0 {
		fmt.Printf("[%s] 初始化失败: 未配置服务列表\n", time.Now().Format(time.DateTime))
		return
	}

	// 多服务：RunMode 是进程级别全局设置，只需设置一次
	if m.RunMode == "" {
		m.RunMode = "release"
	}
	gin.SetMode(m.RunMode)

	// 初始化配置文件
	m.initConfig()
	// 初始化日志
	m.initLog()

	// 初始化每个服务的 router，并注册路由
	for idx, s := range m.Services {
		if s == nil {
			continue
		}
		if s.Name == "" {
			s.Name = fmt.Sprintf("svc-%d", idx)
		}
		m.initServiceRouter(s)
		if s.IsHeartbeat {
			addHeartbeatRoute(s.Router)
		}
		if s.RegRouteFun != nil {
			s.RegRouteFun(s.Router)
		}
	}

	m.isInit = true

}

func (m *MagicApp) Run() error {
	if !m.isInit {
		err := errors.New("未初始化数据")
		fmt.Printf("[%s] server start err: %s\n", time.Now().Format(time.DateTime), err)
		return err
	}
	if len(m.Services) == 0 {
		err := errors.New("未配置服务列表")
		fmt.Printf("[%s] server start err: %s\n", time.Now().Format(time.DateTime), err)
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 构建 server 列表并启动
	var srvs []*http.Server
	errCh := make(chan error, 1)
	for idx, s := range m.Services {
		if s == nil {
			continue
		}
		if s.Name == "" {
			s.Name = fmt.Sprintf("svc-%d", idx)
		}
		if s.Addr == "" || s.Router == nil {
			fmt.Printf("[%s] server(%s) start err: addr/router 为空\n", time.Now().Format(time.DateTime), s.Name)
			continue
		}
		srv := &http.Server{
			Addr:         fmt.Sprintf("%s", s.Addr),
			Handler:      s.Router,
			ReadTimeout:  s.ReadTimeout,
			WriteTimeout: s.WriteTimeout,
			IdleTimeout:  s.IdleTimeout,
		}
		srvs = append(srvs, srv)
		go func(name, addr string, srv *http.Server) {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("[%s] server(%s) listen err(%s): %s\n", time.Now().Format(time.DateTime), name, addr, err)
				select {
				case errCh <- err:
				default:
				}
				// 任一服务启动失败则触发整体退出
				stop()
			}
		}(s.Name, s.Addr, srv)
	}

	if len(srvs) == 0 {
		err := errors.New("未配置可运行的服务")
		fmt.Printf("[%s] server start err: %s\n", time.Now().Format(time.DateTime), err)
		return err
	}

	<-ctx.Done()
	stop()

	var runErr error
	select {
	case runErr = <-errCh:
	default:
	}

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	var shutdownErrs []error
	var shutdownErrMu sync.Mutex
	for _, srv := range srvs {
		wg.Add(1)
		go func(srv *http.Server) {
			defer wg.Done()
			if err := srv.Shutdown(ctx); err != nil {
				fmt.Println("Server forced to shutdown: ", err)
				shutdownErrMu.Lock()
				shutdownErrs = append(shutdownErrs, err)
				shutdownErrMu.Unlock()
			}
		}(srv)
	}
	wg.Wait()

	if m.ExitAfter != nil {
		m.ExitAfter()
	}

	fmt.Println("Server exiting")

	if len(shutdownErrs) > 0 {
		return errors.Join(append([]error{runErr}, shutdownErrs...)...)
	}
	return runErr
}

func (m *MagicApp) initServiceRouter(s *MagicService) {
	if s.Router != nil {
		return
	}
	if s.IsDefault {
		s.Router = gin.Default()
	} else {
		s.Router = gin.New()
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
	if m.ConfHotLoading {
		//viper.OnConfigChange(func(e fsnotify.Event) {})  // 配置变化时的回调
		viper.WatchConfig()
	}
}

// InitLog 初始化日志
func (m *MagicApp) initLog() {
	magicLog := &MagicLog{LogKey: "app.log"}
	magicLog.initLogger()
	fmt.Printf("[%s] 初始化日志...ok\n", time.Now().Format(time.DateTime))
}

// addHeartbeatRoute 心跳检测
func addHeartbeatRoute(r *gin.Engine) {
	if r == nil {
		return
	}
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "success! This service is normal.")
	})
}
