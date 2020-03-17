package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"movie-system/app/kit"
	"movie-system/app/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UserSave(context *gin.Context) {
	userName := context.Param("name")
	context.String(http.StatusOK, "用户"+userName+"已保存")
}
func UserNameAgeSave(context *gin.Context) {
	name := context.Query("name")
	age := context.Query("age")
	context.String(http.StatusOK, "user"+name+age+"save")
}

func UserRegisterJump(ctx *gin.Context) {
	ctx.Request.URL.Path = "login.tmpl"
	ctx.HTML(http.StatusOK, "login.tmpl", nil)

}

func UserRegister(ctx *gin.Context) {
	var user model.User

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.String(http.StatusBadRequest, "输入的数据不合法")
		log.Panicln("err ->", err.Error())
	}
	passwordAgain := ctx.PostForm("password-again")
	if passwordAgain != user.Password {
		ctx.String(http.StatusBadRequest, "密码校验无效，两次密码不一致")
		log.Panicln("密码校验无效，两次密码不一致")
	}

	id := model.Save(&user)
	println("id:", id)

	ctx.Redirect(http.StatusMovedPermanently, "/")

	//email := ctx.PostForm("email")
	//password := ctx.DefaultPostForm("password", "123")
	//passwordAgain := ctx.DefaultPostForm("password-again", "123")

	println("email", user.Email, "password", user.Password, "password again")

}

func UserLogin(context *gin.Context) {
	var user model.User
	if e := context.Bind(&user); e != nil {
		log.Panicln("login 绑定错误", e.Error())
	}

	u := user.QueryByEmail()
	if u.Password == user.Password {
		log.Println("登录成功", u.Email)
		context.SetCookie(kit.COOKIE_NAME, string(u.ID), 60000, "/", "localhost", false, true)
		context.HTML(http.StatusOK, "index.tmpl", gin.H{
			"email": u.Email,
			"id":    u.ID,
			"user":  u,
		})
	}
}

func UserProfile(context *gin.Context) {
	id := context.Query("id")
	var user model.User
	i, err := strconv.Atoi(id)
	u, e := user.QueryById(i)
	if e != nil || err != nil {
		context.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": e,
		})
	}
	context.HTML(http.StatusOK, "user_profile.tmpl", gin.H{
		"user": u,
	})
}

func UpdateUserProfile(context *gin.Context) {
	var user model.User
	if err := context.ShouldBind(&user); err != nil {
		context.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": err.Error(),
		})
		log.Panicln("绑定发生错误 ", err.Error())
	}
	file, e := context.FormFile("avatar-file")
	if e != nil {
		context.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": e,
		})
		log.Panicln("文件上传错误", e.Error())
	}

	path := kit.RootPath()
	path = filepath.Join(path, "avatar")
	e = os.MkdirAll(path, os.ModePerm)
	if e != nil {
		context.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": e,
		})
		log.Panicln("无法创建文件夹", e.Error())
	}
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
	e = context.SaveUploadedFile(file, path+fileName)
	if e != nil {
		context.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": e,
		})
		log.Panicln("无法保存文件", e.Error())
	}
	avatarUrl := "http://localhost:8080/avatar/" + fileName
	user.Avatar = avatarUrl

	//===============================

	/*	fil := make([][]byte, 0)
		var b int64 = 0
		openFile, _ := file.Open()
		// 通过for循环写入
		for {
			buffer := make([]byte, 1024)
			n, err := openFile.ReadAt(buffer, b)
			b = b + int64(n)
			fil = append(fil, buffer)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
		// 生成最后的文件字节流
		fileStream := bytes.Join(fil, []byte(""))

		s := *(*string)(unsafe.Pointer(&fileStream))
		user.Avatar=s*/
	//=================================

	_ = user.Update(user.ID)

	if e != nil {
		context.HTML(http.StatusOK, "error.tmpl", gin.H{
			"error": e,
		})
		log.Panicln("数据无法更新", e.Error())
	}
	context.Redirect(http.StatusMovedPermanently, "/user/profile?id="+strconv.Itoa(user.ID))

}
