package main

import (
	//"encoding/json"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

/*type CountryInfo struct {
	Code           string          `json:"code,omitempty"  db:"Code"`
	Name           string          `json:"name,omitempty"  db:"Name"`
	Continent      string          `json:"continent"  db:"Continent"`
	Region         string          `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64         `json:"surfacearea,omitempty"  db:"SurfaceArea"`
	IndepYear      sql.NullInt64   `json:"indepyear,omitempty"  db:"IndepYear"`
	Population     int             `json:"population,omitempty"  db:"Population"`
	LifeExpectancy sql.NullFloat64 `json:"lifeexpectancy,omitempty"  db:"LifeExpectancy"`
	GNP            sql.NullFloat64 `json:"gnp,omitempty"  db:"GNP"`
	GNPOld         sql.NullFloat64 `json:"gnpold,omitempty"  db:"GNPOld"`
	LocalName      string          `json:"localname,omitempty"  db:"LocalName"`
	GovernmentForm string          `json:"governmentform,omitempty"  db:"GovernmentForm"`
	HeadOfState    sql.NullString  `json:"headofstate,omitempty"  db:"HeadOfState"`
	Capital        sql.NullInt64   `json:"capital,omitempty"  db:"Capital"`
	Code2          string          `json:"code2,omitempty"  db:"Code2"`
}*/
type CountryInfo struct {
	Code           string          `json:"code,omitempty"  db:"Code"`
	Name           string          `json:"name,omitempty"  db:"Name"`
	Continent      string          `json:"continent"  db:"Continent"`
	Region         string          `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64         `json:"surfacearea,omitempty"  db:"SurfaceArea"`
	IndepYear      sql.NullInt64   `json:"indepyear,omitempty"  db:"IndepYear"`
	Population     int             `json:"population,omitempty"  db:"Population"`
	LifeExpectancy sql.NullFloat64 `json:"lifeexpectancy,omitempty"  db:"LifeExpectancy"`
	GNP            sql.NullFloat64 `json:"gnp,omitempty"  db:"GNP"`
	GNPOld         sql.NullFloat64 `json:"gnpold,omitempty"  db:"GNPOld"`
	LocalName      string          `json:"localname,omitempty"  db:"LocalName"`
	GovernmentForm string          `json:"governmentform,omitempty"  db:"GovernmentForm"`
	HeadOfState    sql.NullString  `json:"headofstate,omitempty"  db:"HeadOfState"`
	Capital        sql.NullInt64   `json:"capital,omitempty"  db:"Capital"`
	Code2          string          `json:"code2,omitempty"  db:"Code2"`
}

type Country struct {
	Code string `json:"code,omitempty"  db:"Code"`
	Name string `json:"name,omitempty"  db:"Name"`
}

type CityInfo struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	CityName    string `json:"cityname,omitempty"  db:"CityName"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
	Code        string `json:"code,omitempty"  db:"Code"`
	CountryName string `json:"countryname,omitempty"  db:"Name"`
}

/*
type NullInt64 struct {    // 新たに型を定義
    sql.NullInt64
}

type NullString struct {    // 新たに型を定義
    sql.NullString
}

type someModel struct {
    code NullInt64    // 新しい型を指定
    name NullString    // 新しい型を指定
}

func (ni *NullInt64) UnmarshalJSON(value []byte) error {
    err := json.Unmarshal(value, ni.Int64)
    ni.Valid = err == nil
    return err
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
    if !ni.Valid {
        return json.Marshal(nil)
    }
    return json.Marshal(ni.Int64)    // 値のフィールドのみ返す
}

func (ns *NullString) UnmarshalJSON(value []byte) error {
    err := json.Unmarshal(value, ns.String)
    ns.Valid = err == nil
    return err
}

func (ns NullString) MarshalJSON() ([]byte, error) {
    if !ns.Valid {
        return json.Marshal(nil)
    }
    return json.Marshal(ns.String)    // 値のフィールドのみ返す
}
*/
var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/login", postLoginHandler)
	e.POST("/signup", postSignUpHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)
	withLogin.GET("/cities/:cityName", getCityInfoHandler)
	//withLogin.GET("/countries", getCountryInfoHandler)
	withLogin.GET("/countries", getAllCountryNameHandler)
	withLogin.GET("/countries/:countryName", getCityListHandler)
	withLogin.GET("/whoami", getWhoAmIHandler)

	//e.GET("/countries/:countryName", getCountryInfoHandler)
	e.GET("/login/username", getUserName)
	e.Start(":12500")
}

type LoginRequestBody struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

type Me struct {
	Username string `json:"username,omitempty"  db:"username"`
}

func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	// もう少し真面目にバリデーションするべき
	if req.Password == "" || req.Username == "" {
		// エラーは真面目に返すべき
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	// ユーザーの存在チェック
	var count int

	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username=?", req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	_, err = db.Exec("INSERT INTO users (Username, HashedPass) VALUES (?, ?)", req.Username, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE username=?", req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Username
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))
		name = sess.Values["userName"].(string)

		return next(c)
	}
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city := City{}
	//
	err := db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	//
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}
	//
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong")
	}
	//
	return c.JSON(http.StatusOK, city)
}

func getAllCountryNameHandler(c echo.Context) error {
	//countryName := c.Param("countryName")
	country := []Country{}
	err := db.Select(&country, "SELECT Name,Code FROM country ORDER BY name ")
	/*if country.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}*/
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong")
	}

	return c.JSON(http.StatusOK, country)
}

func getCityListHandler(c echo.Context) error {
	countryName := c.Param("countryName")
	cityInfo := []CityInfo{}
	err := db.Select(&cityInfo, "SELECT city.*,city.name AS CityName,country.Code,country.Name  FROM `city` LEFT OUTER JOIN `country` ON city.CountryCode=country.Code WHERE country.Name=? ORDER BY city.district ASC", countryName)

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong")
	}
	return c.JSON(http.StatusOK, cityInfo)
}

var name = ""

func getUserName(c echo.Context) error {
	return c.String(http.StatusOK, name)
}

func getWhoAmIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Me{
		Username: c.Get("userName").(string),
	})
}
