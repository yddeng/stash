package webservice

import (
	"errors"
	"github.com/gin-gonic/gin"
	"initialthree/cluster"
	"initialthree/cluster/priority"
	"initialthree/node/common/config"
	"initialthree/pkg/event"
	"initialthree/pkg/timer"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	app         *gin.Engine
	taskQueue   *event.EventQueue
	accessToken string
)

func StartWebService(conf *config.WebService) error {
	taskQueue = cluster.GetEventQueue()
	accessToken = conf.AccessToken

	var err error

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	// 跨域
	app.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Expose-Headers", "*")
		if ctx.Request.Method == "OPTIONS" {
			// 处理浏览器的options请求时，返回200状态即可
			ctx.JSON(http.StatusOK, "")
			ctx.Abort()
			return
		}

		ctx.Next()
	})

	for name, m := range modules {
		if err = m.Init(app); err != nil {
			return errors.New("module " + name + " init err " + err.Error())
		}
	}

	cluster.RegisterTimer(time.Second*5, func(timer *timer.Timer, i interface{}) {
		for _, m := range modules {
			m.Tick()
		}
	}, nil)

	go func() {
		t := strings.Split(conf.WebAddress, ":")
		if err := app.Run(":" + t[1]); err != nil {
			panic(err)
		}
	}()

	return nil
}

type Module interface {
	Init(app *gin.Engine) error
	Tick()
}

var (
	modules = map[string]Module{}
)

func checkToken(ctx *gin.Context, route string) bool {
	if accessToken != "" {
		if tkn := ctx.GetHeader("Access-Token"); tkn != accessToken {
			return false
		}
	}
	return true
}

func RegisterModule(name string, module Module) {
	if _, ok := modules[name]; ok {
		panic("register module " + name + " repeated")
	}
	modules[name] = module
}

// 应答结构
type Result struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WaitConn struct {
	route    string
	result   Result
	done     chan struct{}
	doneOnce sync.Once
}

func newWaitConn(route string) *WaitConn {
	return &WaitConn{
		route: route,
		done:  make(chan struct{}),
	}
}

func (this *WaitConn) Done() {
	this.doneOnce.Do(func() {
		if this.result.Message == "" {
			this.result.Success = true
		}
		close(this.done)
	})
}

func (this *WaitConn) GetRoute() string {
	return this.route
}

func (this *WaitConn) SetResult(message string, data interface{}) {
	this.result.Message = message
	this.result.Data = data
}

func (this *WaitConn) Wait() {
	<-this.done
}

func transBegin(ctx *gin.Context, fn interface{}, args ...reflect.Value) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != len(args)+1 {
		panic("func argument error")
	}

	route := getCurrentRoute(ctx)
	wait := newWaitConn(route)
	if err := taskQueue.PostFullReturn(priority.LOW, func() {
		ok := checkToken(ctx, route)
		if !ok {
			wait.SetResult("Token验证失败", nil)
			wait.Done()
			return
		}
		val.Call(append([]reflect.Value{reflect.ValueOf(wait)}, args...))
	}); err != nil {
		wait.SetResult("访问人数过多", nil)
		wait.Done()
	}
	wait.Wait()

	ctx.JSON(http.StatusOK, wait.result)
}

func getCurrentRoute(ctx *gin.Context) string {
	return ctx.FullPath()
}

func getJsonBody(ctx *gin.Context, inType reflect.Type) (inValue reflect.Value, err error) {
	if inType.Kind() == reflect.Ptr {
		inValue = reflect.New(inType.Elem())
	} else {
		inValue = reflect.New(inType)
	}
	if err = ctx.ShouldBindJSON(inValue.Interface()); err != nil {
		return
	}
	if inType.Kind() != reflect.Ptr {
		inValue = inValue.Elem()
	}
	return
}

func WarpHandle(fn interface{}) gin.HandlerFunc {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	switch typ.NumIn() {
	case 1: // func(done *WaitConn)
		return func(ctx *gin.Context) {
			transBegin(ctx, fn)
		}
	case 2: // func(done *WaitConn, req struct)
		return func(ctx *gin.Context) {
			inValue, err := getJsonBody(ctx, typ.In(1))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Json unmarshal failed!",
					"error":   err.Error(),
				})
				return
			}

			transBegin(ctx, fn, inValue)
		}
	default:
		panic("func symbol error")
	}
}
