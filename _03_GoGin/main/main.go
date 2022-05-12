package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{{ID: 123, Name: "张三"}, {ID: 456, Name: "李四"}}

func Handle(r *gin.Engine, httpMethods []string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	var routes gin.IRoutes
	for _, httpMethod := range httpMethods {
		routes = r.Handle(httpMethod, relativePath, handlers...)
	}
	return routes
}

// 自定义中件间
func costTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		nowTime := time.Now()
		// 继续请求
		c.Next()
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("this request URL %s cost %v\n", url, costTime)
	}
}

func main() {
	fmt.Println("1111")

	r := gin.Default()
	r.Use(costTime()) // 添加自定义中件间

	// http://127.0.0.1:8080/?wechat=thisisnohi&a=1&a=2&a=3&map[a]=m1&map[1]=1111
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Blog":   "www.flysnow.org",
			"wechat": "flysnow_org",
			"abc":    c.Query("wechat"),
			"a":      c.QueryArray("a"),
			"b":      "默认值：" + c.DefaultQuery("b", "0"),
			"map":    c.QueryMap("map"),
		})
	})
	r.GET("/users", func(c *gin.Context) {
		c.JSON(200, users)
	})
	// 路由参数
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println("id:" + id)
		c.String(200, "This user id is %s", id)
	})
	// 路由参数
	r.GET("/users/start/*id", func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println("id:" + id)
		c.String(200, "This user id is %s", id)
	})
	// curl -d wechat=1111 http://127.0.0.1:8080/users
	r.POST("/users", func(c *gin.Context) {
		//创建一个用户
		c.JSON(200, gin.H{
			"Blog":   "www.flysnow.org",
			"wechat": "flysnow_org",
			"abc":    c.PostForm("wechat"),
		})
	})
	r.DELETE("/usrs/123", func(context *gin.Context) {
		//删除ID为123的用户
	})
	r.PUT("/usrs/123", func(context *gin.Context) {
		//更新ID为123的用户
	})

	r.PATCH("/usrs/123", func(context *gin.Context) {
		//更新ID为123用户的部分信息
	})

	Handle(r, []string{"GET", "POST"}, "/handler", func(c *gin.Context) {
		//同时注册GET、POST请求方法
		c.JSON(200, gin.H{
			"Blog":    "www.flysnow.org",
			"wechat":  "flysnow_org",
			"handler": "this is handler GET/POST",
		})
	})

	// 分组路由
	group(r)
	// JSONP
	jsonp(r)
	xml(r)
	html(r)
	auth(r)

	r.Run(":8080")
}

// 授权
func auth(r *gin.Engine) {
	group := r.Group("/auth")
	group.Use(gin.BasicAuth(gin.Accounts{"admin": "123456", "nohi": "nohi"}))
	group.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Blog": "nohi.online",
			"abc":  c.Query("wechat"),
			"a":    c.QueryArray("a"),
			"b":    "默认值：" + c.DefaultQuery("b", "0"),
			"map":  c.QueryMap("map"),
		})
	})

}

func group(r *gin.Engine) {

	g := r.Group("/v1")
	{
		g.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"Blog": "www.flysnow.org====g1",
				"abc":  c.Query("wechat"),
				"a":    c.QueryArray("a"),
				"b":    "默认值：" + c.DefaultQuery("b", "0"),
				"map":  c.QueryMap("map"),
			})
		})
		g.GET("/users", func(c *gin.Context) {
			c.JSON(200, users)
		})
		g.GET("/users/allUsers", func(c *gin.Context) {
			c.IndentedJSON(200, users)
		})
		g.GET("/users/123", func(c *gin.Context) {
			c.JSON(200, User{123, "USER姓名", 18})
		})

		// PureJSON
		g.GET("/json", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "<b>Hello, world!</b><hr>",
			})
		})

		g.GET("/pureJson", func(c *gin.Context) {
			c.PureJSON(200, gin.H{
				"message": "<b>Hello, world!</b>",
			})
		})

		g.GET("/asciiJSON", func(c *gin.Context) {
			c.AsciiJSON(200, gin.H{"message": "hello 飞雪无情"})
		})
	}

	g2 := r.Group("/v2", func(c *gin.Context) {
		fmt.Println("=======v2 1111======")
		fmt.Println("a:" + c.Query("a"))
		fmt.Println(c.QueryArray("a"))
	}, func(c *gin.Context) {
		fmt.Println("=======v2 22222======")
		fmt.Println("a:" + c.Query("a"))
	})
	{
		g2.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"Blog": "www.flysnow.org===g2",
				"abc":  c.Query("wechat"),
				"a":    c.QueryArray("a"),
				"b":    "默认值：" + c.DefaultQuery("b", "0"),
				"map":  c.QueryMap("map"),
			})
		})
		g2.GET("/users", func(c *gin.Context) {
			c.JSON(200, users)
		})
	}
}
func jsonp(r *gin.Engine) {
	groupJsonp := r.Group("/jsonp")
	{
		groupJsonp.GET("/", func(c *gin.Context) {
			c.JSONP(200, gin.H{"wechat": "flysnow_org"})
		})
		a := []string{"1", "2", "3"}
		r.SecureJsonPrefix("for(;;);")
		groupJsonp.GET("/secureJson", func(c *gin.Context) {
			c.SecureJSON(200, a)
		})
	}
}
func xml(r *gin.Engine) {
	// xml
	groupXml := r.Group("/xml")
	{
		groupXml.GET("/", func(c *gin.Context) {
			c.XML(200, gin.H{"wechat": "flysnow_org", "blog": "nohi.online"})
		})
		groupXml.GET("/users", func(c *gin.Context) {
			c.XML(200, users)
		})

	}
}

// md5函数
func md5Fun(in string) (string, error) {
	hash := md5.Sum([]byte(in))
	return hex.EncodeToString(hash[:]), nil
}

func html(r *gin.Engine) {
	// xml
	group := r.Group("/html")
	{
		group.GET("/", func(c *gin.Context) {
			c.Status(200)
			const templateText = `this is {{ printf "%s" .}}`
			tmpl, err := template.New("htmlTest").Parse(templateText)
			if err != nil {
				log.Fatalf("parsing: %s", err)
			}
			tmpl.Execute(c.Writer, "nohi.online")
		})

		// 必须loadhtmlfiles前加载
		r.SetFuncMap(template.FuncMap{
			"md5": md5Fun,
		})
		//r.LoadHTMLFiles("html/index.html")
		r.LoadHTMLGlob("html/*")
		group.GET("/index", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{"a": "aaa", "name": "NOHI"})
		})
	}
}
