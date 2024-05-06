package controller

import (
	"lzhuk/clients/pkg/config"
	"strconv"
)

// Точки для работы с сервисом forum-api
var (
	allPost          = "" // Получение всех постов
	registry         = "" // Регистрация нового пользователя
	login            = "" // Вход пользователя в учетную запись
	createPosts      = "" // Создание нового поста пользователем
	userPost         = "" // Получение постов созданных пользователем
	getUserPost      = "" // Получение конкретного поста по индентификатору
	createComments   = "" // Создание коментариев пользователем под постом
	updatePosts      = "" // Изменение поста пользователем
	deletePosts      = "" // Удаление поста пользователем
	votePosts        = "" // Постановка голоса на пост пользователем
	voteComments     = "" // Постанока голоса на комментарий пользователем
	likePosts        = "" // Получение постов которые понравились пользователю
	updateComments   = "" // Изменение комментария пользователем
	deleteComments   = "" // Удаление комментария пользователем
	logoutUsers      = "" // Выход пользовтаеля из учетной записи
	categoryGet      = "" // Получение постов по конкретной категории
	auth             = "" // Регистрация и вход пользователя через стороннего провайдера
	UploadImagePosts = "" // Передача данных о загрузке картинки к новому посту пользователя
)

func InitPointApi(cfg config.Config) {
	port := strconv.Itoa(cfg.PortServer)

	allPost = "http://backend:" + port + "/d3"                       // Получение всех постов
	registry = "http://backend:" + port + "/register"                // Регистрация нового пользователя
	login = "http://backend:" + port + "/login"                      // Вход пользователя в учетную запись
	createPosts = "http://backend:" + port + "/d3/post-create"       // Создание нового поста пользователем
	userPost = "http://backend:" + port + "/d3/user-posts"           // Получение постов созданных пользователем
	getUserPost = "http://backend:" + port + "/d3/post?id="          // Получение конкретного поста по индентификатору
	createComments = "http://backend:" + port + "/d3/comment-create" // Создание коментариев пользователем под постом
	updatePosts = "http://backend:" + port + "/d3/post-update?id="   // Изменение поста пользователем
	deletePosts = "http://backend:" + port + "/d3/post-delete?id="   // Удаление поста пользователем
	votePosts = "http://backend:" + port + "/d3/post-like"           // Постановка голоса на пост пользователем
	voteComments = "http://backend:" + port + "/d3/comment-like"     // Постанока голоса на комментарий пользователем
	likePosts = "http://backend:" + port + "/d3/user-likes"          // Получение постов которые понравились пользователю
	updateComments = "http://backend:" + port + "/d3/comment-update" // Изменение комментария пользователем
	deleteComments = "http://backend:" + port + "/d3/comment-delete" // Удаление комментария пользователем
	logoutUsers = "http://backend:" + port + "/logout"               // Выход пользовтаеля из учетной записи
	categoryGet = "http://backend:" + port + "/d3/category?name="    // Получение постов по конкретной категории
	auth = "http://backend:" + port + "/auth"                        // Регистрация и вход пользователя через стороннего провайдера
	UploadImagePosts = "http://backend:" + port + "/d3/upload-image" // Передача данных о загрузке картинки к новому посту пользователя
}

// Точки для работы с провайдерами Google и Github
const (
	googleAuthEndPoint            = "https://accounts.google.com/o/oauth2/auth"
	googleAuthEndPointAccessToken = "https://accounts.google.com/o/oauth2/token"
	googleUserInfoURL             = "https://www.googleapis.com/oauth2/v3/userinfo"

	githubAuthEndPoint       = "https://github.com/login/oauth/authorize"
	githubAuthEndAccessToken = "https://github.com/login/oauth/access_token"
	githubUserInfoURL        = "https://api.github.com/user"

	client_idGIT     = "21c2671efe47648ceedd"
	client_secretGIT = "fbf46e505b7583bd24c5309bd342379f80591e68"
	callbackGIT      = "http://localhost:8082/github/callback"
)
