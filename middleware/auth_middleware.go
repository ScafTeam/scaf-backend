package middleware

import (
	"backend/database"
	"backend/model"
	"log"
	"net/http"
	"time"

	"github.com/ScafTeam/firebase-go-client/auth"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var IdentityKey = "id"
var AuthMiddleware *jwt.GinJWTMiddleware

func SetupAuthMiddleware(server *gin.Engine) {
	var err error
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",   //标识
		SigningAlgorithm: "HS256",       //加密算法
		Key:              []byte("111"), //密钥
		Timeout:          36 * time.Hour,
		MaxRefresh:       time.Hour,   //刷新最大延长时间
		IdentityKey:      IdentityKey, //指定cookie的id
		SendCookie:       true,
		PayloadFunc: func(data interface{}) jwt.MapClaims { //负载，这里可以定义返回jwt中的payload数据
			v, ok := data.(*model.ScafUser)
			// log.Println(data.(*model.ScafUser))
			if ok {
				return jwt.MapClaims{
					IdentityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Println(claims)
			user := model.ScafUser{}
			user.Email = claims[IdentityKey].(string)
			return &user
		},
		Authenticator: UserLogin,
		Authorizator: func(data interface{}, c *gin.Context) bool { //当用户通过token请求受限接口时，会经过这段逻辑
			v, ok := data.(*model.ScafUser)
			log.Println(v, ok)
			return v != nil && ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) { //错误时响应
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) { //登录成功时响应
			c.JSON(http.StatusOK, gin.H{
				"status":  "Authorized",
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "Sign in success",
			})
		},
		// 指定从哪里获取token 其格式为："<source>:<name>" 如有多个，用逗号隔开
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Println("JWT Error:" + err.Error())
	}
	log.Printf("%p\n", AuthMiddleware.Authenticator)
}

func MemberCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := AuthMiddleware.ParseToken(c)
		claims := jwt.ExtractClaimsFromToken(token)
		email := claims[IdentityKey].(string)

		project_name := c.Param("project_name")
		project_author := c.Param("user_email")

		res, err := database.GetProjectDetail(project_author, project_name)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			c.Abort()
		}

		log.Println("What?")

		if res == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "Not Found",
				"message": "Project not found",
				"debug":   project_author + "/" + project_name,
			})
			c.Abort()
			return
		}

		for _, member := range res.Data()["members"].([]interface{}) {
			if member == email {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not the member of this project",
		})
		c.Abort()
	}
}

func OwnerCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := AuthMiddleware.ParseToken(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": "Internal Server Error",
			})
			c.Abort()
		}
		claims := jwt.ExtractClaimsFromToken(token)
		email := claims[IdentityKey].(string)
		url_email := c.Param("user_email")
		if email != url_email {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "Unauthorized",
				"message": "You are not the author of this project",
			})
			c.Abort()
		}
		c.Abort()
	}
}

func UserLogin(c *gin.Context) (interface{}, error) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return nil, jwt.ErrMissingLoginValues
	}

	res := auth.SignInWithEmailAndPassword(req.Email, req.Password)
	if res.Status() {
		user := res.Result()
		scaf_user := model.ScafUser{
			Email:    user.Email,
			Projects: []string{},
		}
		return &scaf_user, nil
	} else {
		return nil, jwt.ErrFailedAuthentication
	}
}

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(IdentityKey)
	c.JSON(200, gin.H{
		"uesrEmail_claims": claims[IdentityKey].(string),
		"userEmail_c_Get":  user.(*model.ScafUser).Email,
	})
}
