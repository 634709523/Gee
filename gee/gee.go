package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix string
	middleware []HandlerFunc
	parent *RouterGroup
	engine *Engine
}

func New()*Engine{
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine:engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default()*Engine{
	engine := New()
	engine.Use(Logger(),Recovery())
	return engine
}

func (group *RouterGroup)Group(prefix string)*RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:prefix,
		parent:group,
		engine:engine,
	}
	engine.groups = append(engine.groups,newGroup)
	return newGroup
}

func (group *RouterGroup)Use(middleware ...HandlerFunc){
	group.middleware = append(group.middleware,middleware...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine)addRoute(method,pattern string,handler HandlerFunc){
	engine.router.addRoute(method,pattern,handler)
}

func (engine *Engine)GET(parttern string,handler HandlerFunc){
	engine.addRoute("GET",parttern,handler)
}

func (engine *Engine)POST(parttern string,handler HandlerFunc){
	engine.addRoute("POST",parttern,handler)
}

func (engine *Engine)Run(addr string)(err error){
	return http.ListenAndServe(addr,engine)
}



func (engine *Engine)ServeHTTP(w http.ResponseWriter,r *http.Request){
	var middlewares []HandlerFunc
	for _,group := range engine.groups{
		if strings.HasPrefix(r.URL.Path,group.prefix){
			middlewares = append(middlewares,group.middleware...)
		}
	}
	c := newContext(w,r)
	c.handlers = middlewares
	engine.router.handle(c)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix,relativePath)
	fileServer := http.StripPrefix(absolutePath,http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filePath")
		if _,err := fs.Open(file);err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer,c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}