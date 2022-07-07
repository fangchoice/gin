package controller

import (
	"gin-learn/common"
	"gin-learn/dto"
	"gin-learn/model"
	"gin-learn/response"
	"gin-learn/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// 注册用户
func Register(ctx *gin.Context) {
	DB := common.GetDB()
	// get param
	name := ctx.PostForm("name")
	phone := ctx.PostForm("phone")
	pwd := ctx.PostForm("pwd")

	// 数据验证
	if len(phone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		// 	"code": 200,
		// 	"msg":  "手机号必须是11位",
		// })

		return
	}

	if len(pwd) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")

		return
	}

	// 如果名称没有传，给一个10为随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, phone, pwd)
	// 判断手机号是否存在( 需要数据库 )
	if isPhoneExist(DB, phone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")

		return
	}

	// 创建用户 加密密码
	hashdPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "加密错误")

		return
	}
	newUser := model.User{
		Name:  name,
		Phone: phone,
		Pwd:   string(hashdPwd),
	}
	DB.Create(&newUser)

	// 返回结果
	// ctx.JSON(200, gin.H{
	// 	"code": 200,
	// 	"msg":  "注册成功",
	// })
	response.Success(ctx, nil, "注册成功")
}

// 登录
func Login(ctx *gin.Context) {
	DB := common.GetDB()

	// 获取参数
	phone := ctx.PostForm("phone")
	pwd := ctx.PostForm("pwd")

	// 数据验证
	if len(phone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		// 	"code": 200,
		// 	"msg":  "手机号必须是11位",
		// })

		return
	}

	if len(pwd) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		// 	"code": 422,
		// 	"msg":  "密码不能小于6位",
		// })

		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Where("phone = ?", phone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		// 	"code": 422,
		// 	"msg":  "用户不存在",
		// })

		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(pwd)); err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		// 	"code": 400,
		// 	"msg":  "密码错误",
		// })

		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		// 	"code": 500,
		// 	"msg":  "系统异常",
		// })
		log.Printf("token generate error : %v", err)

		return
	}

	// 返回结果
	// ctx.JSON(200, gin.H{
	// 	"code": 200,
	// 	"msg":  "登录成功",
	// 	"data": gin.H{"token": token},
	// })
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// Info
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

// -------------------------------------------------------------------------------------
// 判断手机号是否存在
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}
