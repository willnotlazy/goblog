package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter().StrictSlash(true)
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "luojy",
		Passwd:               "secret",
		Addr:                 "192.168.4.22:33306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接
	err = db.Ping()
	checkError(err)
}

func createTables() {
	createArtilcesSQL := `CREATE TABLE IF NOT EXISTS articles(
	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	body longtext COLLATE utf8mb4_unicode_ci
);`
	_, err := db.Exec(createArtilcesSQL)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑,请联系我们.</p>")
}

type Article struct {
	Title string
	Body  string
	ID    int64
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)

		err = tmpl.Execute(w, article)
		checkError(err)
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

type articlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesStoreHandle(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务内部错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()

		data := articlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		tml, err := template.ParseFiles("resources/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}

		err = tml.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func saveArticleToDB(title string, body string) (int64, error) {
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1.获取prepare声明语句
	if stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)"); err != nil {
		return 0, err
	}

	// 延迟关闭预执行连接，避免占用连接
	defer stmt.Close()

	// 绑定预执行的参数
	if rs, err = stmt.Exec(title, body); err != nil {
		return 0, err
	}

	// 插入成功返回自增id
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

func articleCreateHandle(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := articlesFormData{
		URL: storeURL,
	}

	tml, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tml.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func articleEditHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := articlesFormData{
			URL:    updateURL,
			Title:  article.Title,
			Body:   article.Body,
			Errors: nil,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkError(err)

		err = tmpl.Execute(w, data)
		checkError(err)
	}
}

func articleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	_, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"

			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)

				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "你没有做任何修改")
			}
		} else {
			updateURL, err := router.Get("articles.update").URL("id", id)
			checkError(err)

			data := articlesFormData{
				URL:    updateURL,
				Title:  title,
				Body:   body,
				Errors: errors,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkError(err)

			err = tmpl.Execute(w, data)
			checkError(err)
		}
	}
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Myname", "goblog")

		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)

	return vars[parameterName]
}

func getArticleByID(id string) (Article, error) {
	query := "SELECT * FROM articles where id = ?"
	article := Article{}
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err
}

func validateArticleFormData(title, body string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

func main() {

	initDB()
	createTables()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandle).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articleCreateHandle).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9+]}/edit", articleEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9+]}", articleUpdateHandler).Methods("POST").Name("articles.update")

	// 404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// middlerware
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":8088", removeTrailingSlash(router))
}
