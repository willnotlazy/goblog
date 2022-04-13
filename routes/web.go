package routes

import (
	"github.com/gorilla/mux"
	"goblog/controllers"
	"goblog/http/middlerwares"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router) {
	// use middleware
	r.Use(middlerwares.StartSession)

	pc := controllers.NewPageController()

	// 前段资源
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	//r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	// 404
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	ac := controllers.NewArticlesController()
	r.HandleFunc("/", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9+]}", ac.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9+]}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	//r.Use(middlerwares.ForceHTML)
	// 用户认证
	auc := controllers.NewAuthController()
	r.HandleFunc("/auth/register", auc.Register).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", auc.DoRegister).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login", auc.Login).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/dologin", auc.DoLogin).Methods("POST").Name("auth.dologin")
	r.HandleFunc("/auth/logout", auc.Logout).Methods("POST").Name("auth.logout")

	// forgot password
	fpc := controllers.NewPasswordForgotController()
	r.HandleFunc("/password/forgot", fpc.ForgotPassword).Methods("GET").Name("password.forgotpassword")
	r.HandleFunc("/password/mail", fpc.ForgotPasswordMail).Methods("POST").Name("password.forgot_password_mail")
	r.HandleFunc("/password/reset", fpc.ResetPassword).Methods("GET").Name("password.reset")
}
