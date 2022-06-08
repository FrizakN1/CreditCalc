package router

import (
	"creditCalc/database"
	"creditCalc/utils"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Initialized() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("template/*.html")
	router.Static("assets", "assets")
	word := sessions.NewCookieStore([]byte("SecretX"))
	router.Use(sessions.Sessions("session", word))

	router.GET("/", index)
	router.GET("/registration", registration)
	router.GET("/login", login)
	router.GET("/credits", getCredits)

	routerUser := router.Group("/user")
	routerUser.POST("/login", loginCheck)
	routerUser.POST("/checkEmail", checkEmail)
	routerUser.PUT("/addUser", addUser)
	routerUser.DELETE("/exit", exit)

	routerAPI := router.Group("/api")
	routerAPI.PUT("/applyCredit", applyCredit)

	return router
}

func getCredits(c *gin.Context) {
	session := getSession(c)
	if session.User.ID > 0 {
		var credits []database.Credit
		credits = session.GetCredits()
		c.HTML(200, "credits", gin.H{
			"Title":   "Мои кредиты",
			"ID":      session.User.ID,
			"Credits": credits,
		})
	} else {
		c.Redirect(301, "/login")
	}
}

func exit(c *gin.Context) {
	session := sessions.Default(c)
	_session := getSession(c)

	_, ok := session.Get("SessionSecretKey").(string)
	if ok {
		session.Clear()
		_ = session.Save()
		c.SetCookie("hello", "", -1, "/", c.Request.URL.Hostname(), false, true)
		session.Delete("SessionSecretKey")
	}

	_session.DeleteSession()

	c.JSON(301, true)
}

func loginCheck(c *gin.Context) {
	session := sessions.Default(c)

	var user database.User
	e := c.BindJSON(&user)
	if e != nil {
		utils.Logger.Println(e)
		c.Status(400)
		return
	}

	user.Password, e = utils.Encrypt(user.Password)
	if e != nil {
		utils.Logger.Println(e)
		c.Status(400)
		return
	}

	if user.LoginCheck() {
		hash, ok := database.CreateSession(&user)
		if ok {
			session.Set("SessionSecretKey", hash)
			e = session.Save()
			if e != nil {
				utils.Logger.Println(e)
				return
			}

			c.JSON(200, true)

			return
		}
	}

	c.JSON(400, nil)
}

func login(c *gin.Context) {
	session := getSession(c)

	if session.User.ID > 0 {
		c.Redirect(301, "/")
	} else {
		c.HTML(200, "login", gin.H{
			"Title": "Авторизация",
			"JS":    "login.js",
			"ID":    session.User.ID,
		})
	}
}

func addUser(c *gin.Context) {
	var userData database.User
	e := c.BindJSON(&userData)
	if e != nil {
		utils.Logger.Println(e)
		return
	}

	c.JSON(200, userData.AddUser())
}

func checkEmail(c *gin.Context) {
	var email database.User
	e := c.BindJSON(&email)
	if e != nil {
		utils.Logger.Println(e)
		c.Status(400)
		return
	}

	c.JSON(200, database.CheckEmail(email.Email))
}

func registration(c *gin.Context) {
	session := getSession(c)
	fmt.Println(session.User.ID)

	if session.User.ID > 0 {
		c.Redirect(301, "/")
	} else {
		c.HTML(200, "registration", gin.H{
			"JS":    "registration.js",
			"Title": "Регистрация",
			"ID":    session.User.ID,
		})
	}
}

func applyCredit(c *gin.Context) {
	session := getSession(c)
	if session.User.ID > 0 {
		var creditData database.Credit

		e := c.BindJSON(&creditData)
		if e != nil {
			utils.Logger.Println(e)
			return
		}

		c.JSON(200, session.ApplyCredit(creditData))
	} else {
		c.JSON(200, false)
	}

}

func index(c *gin.Context) {
	session := getSession(c)

	c.HTML(200, "index", gin.H{
		"JS":    "main.js",
		"Title": "Кредитный калькулятор",
		"ID":    session.User.ID,
	})
}

func getSession(c *gin.Context) *database.Session {
	_session := sessions.Default(c)

	sessionHash, ok := _session.Get("SessionSecretKey").(string)
	if ok {
		session := database.GetSession(sessionHash)
		if session != nil {
			session.Exists = true
			return session
		}
	}

	return &database.Session{
		Exists: false,
	}
}
