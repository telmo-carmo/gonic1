package main

/*

go run github.com/99designs/gqlgen init

https://www.agiliq.com/blog/2020/04/graphql-api-server-using-golang-gin-framework/

https://gqlgen.com/getting-started/


--
	//configuration to SPA applications
	router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.File("./frontend/build/index.html")
		}
		//default 404 page not found
	})
*/

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"tac/gonic1/controllers"
	"tac/gonic1/dal"
)

const (
	MY_JWT_SECRET = "Very Secret Key for alg HS256!"
)

var ErrBadAuth = errors.New("bad auth token")

type UserApi struct {
	ID       uint64 `json:"id"`
	Username string `json:"uid"`
	Password string `json:"pwd"`
	Role     string `json:"role"`
}

var user = UserApi{
	ID:       1,
	Username: "demo",
	Password: "123",
	Role:     "admin",
}

func loginApi(c *gin.Context) {
	var u UserApi
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "gonic1",
		ID:        strconv.FormatUint(user.ID, 10),
		Subject:   user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwt_secret string
	var ok bool
	if jwt_secret, ok = os.LookupEnv("JWT_SECRET"); !ok || jwt_secret == "" {
		jwt_secret = MY_JWT_SECRET
	}

	stoken, err := token.SignedString([]byte(jwt_secret))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	log.Printf("[INFO] Login token= Authorization: Bearer %s", stoken)
	c.JSON(http.StatusOK, gin.H{"username": u.Username, "token": "Bearer " + stoken})
}

func logoutApi(c *gin.Context) {
	log.Println("[INFO] LOGOUT!")
	c.JSON(http.StatusOK, gin.H{"data": "logout!"})
}

func TokenValid(c *gin.Context) error {
	bearToken := c.Request.Header.Get("Authorization")

	//normally Authorization: Bearer the_token_xxx
	strArr := strings.Split(bearToken, " ")
	sjwt := ""
	if len(strArr) == 2 {
		sjwt = strArr[1]
	}

	var jwt_secret string
	var ok bool
	if jwt_secret, ok = os.LookupEnv("JWT_SECRET"); !ok || jwt_secret == "" {
		jwt_secret = MY_JWT_SECRET
	}

	token, err := jwt.Parse(sjwt, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwt_secret), nil
	})

	uIdName := "0;NA"
	//if rclaims, ok := token.Claims.(jwt.RegisteredClaims); !ok && !token.Valid {
	if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		log.Println("[ERROR] bad auth header token", err)
		return ErrBadAuth
	} else {
		log.Printf("[WARN] JWT Claims: %+v", claims)
		uIdName = fmt.Sprintf("%s;%s", claims["jti"], claims["sub"])
		c.Set("uid", uIdName)
		//log.Printf("[WARN] JWT Claims: %v", rclaims)
		//uIdName = rclaims.ID + ";" + rclaims.Subject + ";exp:" + rclaims.ExpiresAt.String()
	}

	log.Printf("[INFO] Auth header token OK for: %s", uIdName)
	return nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := TokenValid(c); err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func spa_handler(c *gin.Context) {
	c.File("./frontend/build/index.html")
}

func main() {
	addrP := flag.String("addr", "0.0.0.0:5000", "server address")
	coloredP := flag.Bool("color", false, "show colored logging")
	flag.Parse()

	//gin.SetMode(gin.ReleaseMode)
	if !*coloredP {
		gin.DisableConsoleColor()
	}

	// starting server config
	r := gin.Default()

	dal.InitDatabase()

	// Serve SPA frontend static files
	//router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))

	//r.Static("/app", "./frontend/build")
	r.Static("/static", "./frontend/build/static")
	r.GET("/app", spa_handler)
	r.GET("/app/:route", spa_handler)

	r.Static("/sta", "./public")
	r.StaticFS("/docs", http.Dir("./updocs"))
	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	r.LoadHTMLGlob("templates/*")

	r.GET("/health", func(c *gin.Context) {
		gv := fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		log.Printf("[INFO] /health %s ", gv)
		c.JSON(http.StatusOK, gin.H{"status": "UP", "now": time.Now().Format("2006-01-02T15:04:05.999Z"), "GO": gv})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"now":     time.Now().UTC(),
		})
	})
	r.GET("/echo/:name", func(c *gin.Context) {
		name := c.Param("name")
		log.Printf("[INFO] %s ; echo %s", c.Request.URL, name)
		c.String(http.StatusOK, "Hello %s", name)
	})
	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"title": "Upload Some Doc",
			"dt1":   time.Now().UTC(),
		})
	})
	r.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		log.Printf("[INFO] Uploading %s from %s", name, file.Filename)

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filepath.Join("updocs", filename)); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		c.JSON(http.StatusOK, gin.H{name: "http://" + c.Request.Host + "/docs/" + filename})
	})

	r.POST("/auth/login", loginApi)
	r.GET("/auth/logout", logoutApi)

	rv1 := r.Group("/api")
	{
		// Routes /api
		rv1.GET("/books", controllers.FindBooks)
		rv1.GET("/books/:id", controllers.FindBook)
		rv1.POST("/books", controllers.CreateBook)
		rv1.PUT("/books/:id", controllers.UpdateBook)
		rv1.DELETE("/books/:id", controllers.DeleteBook)

		rv1.GET("/bonus", controllers.GetAllBonus)
	}
	rv1p := rv1.Group("/")
	rv1p.Use(TokenAuthMiddleware())
	{

		rv1p.GET("/bonus/:id", controllers.GetOneBonus)
		rv1p.POST("/bonus", controllers.CreateBonus)
		rv1p.PUT("/bonus/:id", controllers.UpdateBonus)
		rv1p.DELETE("/bonus/:id", controllers.DeleteBonus)
	}

	log.Printf("Running on  %s", *addrP)

	r.Run(*addrP)
}
