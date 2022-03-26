package apiserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"oliujunk/server/apiserver/authentication"
	"oliujunk/server/config"
	"oliujunk/server/database"
	"time"
)

func Start() {
	log.Println("接口服务启动")
	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(cors())

	router.POST("/api/login", login)

	authenticated := router.Group("/api")
	authenticated.Use(authentication.JWTAuth())
	{
		authenticated.GET("/current/:deviceID", getCurrent)
		authenticated.GET("/datas", getDatas)
		authenticated.GET("/devices", getDevices)
		authenticated.POST("/devices", postDevice)
		authenticated.PUT("/devices/:deviceID", putDevice)
		authenticated.DELETE("/devices/:deviceID", deleteDevice)
		authenticated.PUT("/users", putUser)
	}

	_ = router.Run(fmt.Sprintf(":%d", config.GlobalConfiguration.ApiServer.Port))
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func login(context *gin.Context) {
	type Result struct {
		Token string `json:"token"`
	}
	type Param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var param Param
	err := context.BindJSON(&param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	var user database.User
	result, err := database.Orm.Table("xph_user").Where("username = ?", param.Username).Get(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  -1,
			"message": "服务端异常: " + err.Error(),
		})
		return
	}

	if !result {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "用户名或密码错误",
		})
		return
	}
	if user.Username == param.Username && user.Password == param.Password {
		token := authentication.GenerateToken(user)
		context.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "登录成功",
			"data":    Result{Token: token},
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "用户名或密码错误",
		})
	}
}

func getCurrent(context *gin.Context) {
	deviceID := context.Param("deviceID")
	current := database.Current{}
	_, _ = database.Orm.
		Table("xph_current").
		Where("device_id = ?", deviceID).
		Desc("data_time").
		Get(&current)
	context.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "请求成功",
		"data":    current,
	})
}

func getDatas(context *gin.Context) {
	type Param struct {
		DeviceID  int    `form:"deviceID" binding:"required"`
		PageNum   int    `form:"pageNum" binding:"required"`
		PageSize  int    `form:"pageSize" binding:"required"`
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
	}
	type Result struct {
		List     []database.Current `json:"list"`
		Total    int64              `json:"total"`
		PageNum  int                `json:"pageNum"`
		PageSize int                `json:"pageSize"`
	}
	var param Param
	err := context.ShouldBindQuery(&param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	var datas []database.Current
	_ = database.Orm.
		Table("xph_current").
		Where("device_id = ?", param.DeviceID).
		And("? IS NULL OR ? = '' OR data_time >= ?", param.StartTime, param.StartTime, param.StartTime).
		And("? IS NULL OR ? = '' OR data_time <= ?", param.EndTime, param.EndTime, param.EndTime).
		Limit(param.PageSize, (param.PageNum-1)*param.PageSize).
		Desc("data_time").
		Find(&datas)
	total, _ := database.Orm.
		Table("xph_current").
		Where("device_id = ?", param.DeviceID).
		And("? IS NULL OR ? = '' OR data_time >= ?", param.StartTime, param.StartTime, param.StartTime).
		And("? IS NULL OR ? = '' OR data_time <= ?", param.EndTime, param.EndTime, param.EndTime).
		Desc("data_time").
		Count()
	context.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "请求成功",
		"data": Result{
			Total:    total,
			PageNum:  param.PageNum,
			PageSize: param.PageSize,
			List:     datas,
		},
	})
}

func getDevices(context *gin.Context) {
	claims := context.MustGet("claims").(*authentication.CustomClaims)
	type Param struct {
		DeviceID int `form:"deviceID"`
	}
	var param Param
	_ = context.ShouldBindQuery(&param)
	var devices []database.Device
	_ = database.Orm.
		Table("xph_device").
		Where("creator_id = ?", claims.UserID).
		And("? IS NULL OR ? = '' OR device_id = ?", param.DeviceID, param.DeviceID, param.DeviceID).
		Asc("id").
		Find(&devices)
	context.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "请求成功",
		"data":    devices,
	})
}

func postDevice(context *gin.Context) {
	claims := context.MustGet("claims").(*authentication.CustomClaims)
	var param database.Device
	err := context.BindJSON(&param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	param.CreatorID = claims.UserID
	param.CreateTime = time.Now()
	param.UpdateTime = time.Now()
	var device database.Device
	result, err := database.Orm.Table("xph_device").Where("device_id = ?", param.DeviceID).Get(&device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  -1,
			"message": "服务端异常: " + err.Error(),
		})
		return
	}
	if result {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "设备已存在",
		})
		return
	}
	_, err = database.Orm.Table("xph_device").Insert(param)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  -1,
			"message": "服务端异常: " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "添加成功",
		"data":    true,
	})
}

func putDevice(context *gin.Context) {
	claims := context.MustGet("claims").(*authentication.CustomClaims)
	deviceID := context.Param("deviceID")
	var param database.Device
	err := context.BindJSON(&param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	var device database.Device
	result, err := database.Orm.Table("xph_device").
		Where("creator_id = ?", claims.UserID).
		And("device_id = ?", deviceID).
		Get(&device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  -1,
			"message": "服务端异常: " + err.Error(),
		})
		return
	}
	if !result {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "设备不存在或无操作权限",
		})
		return
	}
	_, err = database.Orm.Table("xph_device").ID(device.ID).Update(param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "修改成功",
		"data":    true,
	})
}

func deleteDevice(context *gin.Context) {
	claims := context.MustGet("claims").(*authentication.CustomClaims)
	deviceID := context.Param("deviceID")
	rows, err := database.Orm.Table("xph_device").
		Where("creator_id = ?", claims.UserID).
		And("device_id = ?", deviceID).
		Delete(database.Device{})
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	if rows >= 1 {
		context.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "删除成功",
			"data":    true,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "删除失败",
		})
	}
}

func putUser(context *gin.Context) {
	claims := context.MustGet("claims").(*authentication.CustomClaims)
	var param database.User
	err := context.BindJSON(&param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	_, err = database.Orm.Table("xph_user").ID(claims.UserID).Update(param)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "修改成功",
		"data":    true,
	})
}
