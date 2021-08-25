package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多63个中间件

type HandlerFunc func(*SliceRouterContext)

//route 结构体
type SliceRouter struct {
	Groups []*SliceGroup
}

// group 结构体
type SliceGroup struct {
	*SliceRouter
	Path     string
	Handlers []HandlerFunc
}

//route 上下文
type SliceRouterContext struct {
	Rw  http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	Index int8
}

type SliceRouterHandler struct {
	CoreFunc func(*SliceRouterContext) http.Handler
	Router   *SliceRouter
}

func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}
	//最长url前缀匹配
	matchUrlLen := 0
	for _, group := range r.Groups {
		if strings.HasPrefix(req.RequestURI, group.Path) {
			pathLen := len(group.Path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*newSliceGroup = *group //浅拷贝指针
			}
		}
	}
	c := &SliceRouterContext{Rw: rw, Req: req, Ctx: req.Context(), SliceGroup: newSliceGroup}
	c.Reset()
	return c
}

//服务触发的时候  真正调用的是这个方法
func (w *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := newSliceRouterContext(rw, req, w.Router)
	if w.CoreFunc != nil {
		c.Handlers = append(c.Handlers, func(src *SliceRouterContext) {
			w.CoreFunc(c).ServeHTTP(rw, req)
		})
	}
	c.Reset()
	c.Next()
}

//构造sliceRouterFunc
func NewSliceRouterHandler(coreFunc func(*SliceRouterContext) http.Handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		CoreFunc: coreFunc,
		Router:   router,
	}
}

//构造router 返回sliceRouter结构体
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

//构造group  主要是存放路由前缀  并且给sliceRouter放进sliceGroup中
func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		Path:        path,
		SliceRouter: g,
	}
}

//构造回调方法
func (g *SliceGroup) Use(middlewares ...HandlerFunc) *SliceGroup {
	g.Handlers = append(g.Handlers, middlewares...)
	existsFlag := false
	//该处判断有没有加入到 sliceRouter的数组中去 如果没有，就加入其中去
	for _, oldGroup := range g.SliceRouter.Groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.SliceRouter.Groups = append(g.SliceRouter.Groups, g)
	}
	return g
}

//获取中间件上下文变量
func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

//设置中间件上下文变量
func (c *SliceRouterContext) Set(key, value interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, value)
}

// 从最先加入中间件开始回调
func (c *SliceRouterContext) Next() {
	c.Index++
	for c.Index < int8(len(c.Handlers)) {
		c.Handlers[c.Index](c)
		c.Index++
	}
}

// 跳出中间件方法
func (c *SliceRouterContext) Abort() {
	c.Index = abortIndex
}

// 是否跳过了回调
func (c *SliceRouterContext) IsAbort() bool {
	return c.Index >= abortIndex
}

//重置回调
func (c *SliceRouterContext) Reset() {
	c.Index = -1
}
