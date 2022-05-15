package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var sqlDb *sql.DB           //数据库连接db
var sqlResponse SqlResponse //响应client的数据

// SqlResponse 应答体（响应client的请求）
type SqlResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

/**
create table products
(
    ID             int primary key,
    Number         varchar(50),
    Category       varchar(50),
    Name           varchar(50),
    made_in         varchar(50),
    production_time datetime
);
*/

// Product 特别注意：结构体名称为：Product，创建的表的名称为：Products
type Product struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Number         string    `gorm:"unique" json:"number"`                       //商品编号（唯一）
	Category       string    `gorm:"type:varchar(256);not null" json:"category"` //商品类别
	Name           string    `gorm:"type:varchar(20);not null" json:"name"`      //商品名称
	MadeIn         string    `gorm:"type:varchar(128);not null" json:"made_in"`  //生产地
	ProductionTime time.Time `json:"production_time"`                            //生产时间
}

//应答体
type GormResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

var gormDB *gorm.DB
var gormResponse GormResponse

func init() {
	var err error
	sqlStr := "root:nohi1234@tcp(10.0.0.210:3306)/go_test?charset=utf8mb4&parseTime=true&loc=Local"
	gormDB, err = gorm.Open(mysql.Open(sqlStr), &gorm.Config{}) //配置项中预设了连接池 ConnPool
	if err != nil {
		fmt.Println("数据库连接出现了问题：", err)
		return
	}

}

func main() {
	r := gin.Default()
	//数据库的CRUD--->gin的 post、get、put、delete方法
	r.POST("gorm/insert", gormInsertData) //添加数据
	r.GET("gorm/get", gormGetData)        //查询数据（单条记录）
	r.GET("gorm/mulget", gormGetMulData)  //查询数据（多条记录）
	//r.PUT("gorm/update", gormUpdateData)    //更新数据
	//r.DELETE("gorm/delete", gormDeleteData) //删除数据
	r.Run(":8080")
}

func gormGetMulData(c *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	name := c.Query("name")
	var products []Product
	tx := gormDB.Where("name like ?", "%"+name+"%").Find(&products)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	gormResponse.Code = http.StatusOK
	gormResponse.Message = "读取成功"
	gormResponse.Data = products
	c.JSON(http.StatusOK, gormResponse)
}

func gormGetData(c *gin.Context) {
	//=============捕获异常============
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	//============
	number := c.Query("number")
	product := Product{}
	tx := gormDB.Where("number=?", number).First(&product)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	gormResponse.Code = http.StatusOK
	gormResponse.Message = "读取成功"
	gormResponse.Data = product
	c.JSON(http.StatusOK, gormResponse)
}

func gormInsertData(c *gin.Context) {
	fmt.Println("========gormInsertData===========")
	//=============捕获异常============
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	//============
	var p Product
	err := c.Bind(&p)
	p.ProductionTime = time.Now()
	fmt.Println(p)
	fmt.Println("=======================")
	if err != nil {
		fmt.Println("========error===============", err)
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "参数错误"
		gormResponse.Data = err
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	fmt.Println("===========gormDB.Create(&p)============")
	tx := gormDB.Create(&p)
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "写入成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	fmt.Printf("insert failed, err:%v\n", err)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "写入失败"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
	fmt.Println(tx) //打印结果
}
