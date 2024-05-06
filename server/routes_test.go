package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	comment2 "gitea.com/lzhuk/forum/internal/repository/comment"
	post2 "gitea.com/lzhuk/forum/internal/repository/post"
	user2 "gitea.com/lzhuk/forum/internal/repository/user"
	"gitea.com/lzhuk/forum/internal/server"
	"gitea.com/lzhuk/forum/internal/service"
	"gitea.com/lzhuk/forum/internal/service/comment"
	"gitea.com/lzhuk/forum/internal/service/post"
	"gitea.com/lzhuk/forum/internal/service/user"
	"gitea.com/lzhuk/forum/pkg/config"
	"gitea.com/lzhuk/forum/pkg/db/driver"
)

type KeyUser string

func TestStartPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}
	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.Home)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/userd3", nil)
	if err != nil {
		t.Fatal(err)
	}
	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}
	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.Home)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"post_id":1,"user_id":1,"category_name":"Поезда","title":"Иволга 3.0","description":"Новый поезд","create_at":"2024-04-25T13:55:08.10608599+06:00","name":"Леонид Жук","likes":1,"dislikes":0},{"post_id":2,"user_id":1,"category_name":"Другое","title":"Стройка ВСМ","description":"Пока строят","create_at":"2024-04-25T13:55:25.335532621+06:00","name":"Леонид Жук","likes":0,"dislikes":0},{"post_id":3,"user_id":1,"category_name":"Тарифы","title":"Карта Тройка","description":"Система оплаты","create_at":"2024-04-25T13:56:15.285234306+06:00","name":"Леонид Жук","likes":0,"dislikes":0},{"post_id":4,"user_id":1,"category_name":"Другое","title":"МЦД-4","description":"Новые станции","create_at":"2024-04-25T13:56:56.011725766+06:00","name":"Леонид Жук","likes":1,"dislikes":0},{"post_id":5,"user_id":1,"category_name":"Станции","title":"Быково","description":"Когда построят переезд","create_at":"2024-04-25T18:05:06.460383264+06:00","name":"Леонид Жук","likes":0,"dislikes":0},{"post_id":6,"user_id":3,"category_name":"Станции","title":"HAHAHEHE","description":"HEHEHAHA","create_at":"2024-05-06T12:34:37.45018+05:00","name":"Zhalgas Bolatov","likes":0,"dislikes":1}]
`
	if rr.Body.String() != expected {
		println(len(expected))
		println(len(rr.Body.String()))
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRegisterExistingAccount(t *testing.T) {
	payload := map[string]string{
		"name":     "Danial",
		"email":    "danial@gmail.com", // Should be an existing gmail
		"password": "1234512345",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/register", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}
	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.Register)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "{\"status\":500,\"message\":\"Name already exist\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRegisterNewAccount(t *testing.T) {
	payload := map[string]string{
		"name":     "Danial",
		"email":    "danial@gmail.com.somerandom.comx", // Should be random always.
		"password": "1234512345",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/register", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}

	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.Register)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLoginIncorrectCredentials(t *testing.T) {
	payload := map[string]string{
		"email":    "danial@gmail.com",
		"password": "someincorrectpassword",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/login", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}

	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.Login)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	expectedMessage := "Invalid Credentials"
	json.Unmarshal(rr.Body.Bytes(), &response)
	if value, ok := response["message"]; !ok || value != expectedMessage {
		fmt.Println(response)
		t.Errorf("handler returned unexpected body, missing 'name'")
	}
}

func TestLoginCorrectCredentials(t *testing.T) {
	payload := map[string]string{
		"name":     "Danial",
		"email":    "danial@gmail.com.somerandom.comx", // Should be random always.
		"password": "1234512345",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/login", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}

	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.Login)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	if _, ok := response["name"]; !ok {
		fmt.Println(response)
		t.Errorf("handler returned unexpected body, missing 'name'")
	}
	if _, ok := response["email"]; !ok {
		t.Errorf("handler returned unexpected body, missing 'email'")
	}
}

func TestLikePosts(t *testing.T) {
	req, err := http.NewRequest("GET", "/userd3/likeposts", nil)
	if err != nil {
		t.Fatal(err)
	}

	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}

	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.LikePosts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestMyPosts(t *testing.T) {
	req, err := http.NewRequest("GET", "/userd3/myposts", nil)
	if err != nil {
		t.Fatal(err)
	}

	services, rr, err := initServices()
	if err != nil {
		t.Errorf("Error initializing services: %s", err)
	}

	h := server.NewHandler(services)
	handler := http.HandlerFunc(h.PostsUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// func TestCreatePosts(t *testing.T) {
// 	payload := map[string]string{
// 		"cookie_uuid":   "ad07db8a-a88b-4a31-b37e-49c398886756",
// 		"category_name": "Станции",
// 		"title":         "Станция Люберцы",
// 		"discription":   "Новые поезда",
// 	}
// 	payloadBytes, _ := json.Marshal(payload)
// 	req, err := http.NewRequest("POST", "/userd3/posts", bytes.NewReader(payloadBytes))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	user := &model.User{
// 		ID:    1337,
// 		Name:  "Danial",
// 		Email: "danial@example.com",
// 	}

// 	req = setUserContext(req, user)
// 	cookie := &http.Cookie{
// 		Name:  "UserData",
// 		Value: "ad07db8a-a88b-4a31-b37e-49c398886756",
// 	}
// 	req.AddCookie(cookie)
// 	services, rr, err := initServices()
// 	if err != nil {
// 		t.Errorf("Error initializing services: %s", err)
// 	}

// 	h := server.NewHandler(services)
// 	handler := http.HandlerFunc(h.CreatePosts)
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }

func initServices() (service.Service, *httptest.ResponseRecorder, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		return service.Service{}, nil, errors.New("Ошибка загрузки конфигурации")
	}

	db, err := driver.NewDB(cfg)
	if err != nil {
		log.Println("Ошибка инциализации базы данных %w", err)
		return service.Service{}, nil, errors.New("Ошибка загрузки конфигурации")
	}

	usersRepo := user2.NewUserRepo(db)
	sessionRepo := user2.NewSessionRepository(db)
	postsRepo := post2.NewPostsRepo(db)
	likePostRepo := post2.NewLikePostRepository(db)
	commentsRepo := comment2.NewCommentsRepo(db)
	likeCommentsRepo := comment2.NewLikeCommentRepository(db)
	uploadImageRepo := post2.NewUploadImagePostRepository(db)

	usersService := user.NewUserService(usersRepo)
	sessionsService := user.NewSessionService(sessionRepo)
	postsService := post.NewPostsService(postsRepo)
	likePostsService := post.NewLikePostService(likePostRepo)
	commentsService := comment.NewCommentsService(commentsRepo)
	likeCommentsService := comment.NewLikeCommentService(likeCommentsRepo)
	uploadImageService := post.NewUploadImagePostService(uploadImageRepo)

	services := service.NewService(usersService, sessionsService, postsService, commentsService, likePostsService, likeCommentsService, uploadImageService)

	rr := httptest.NewRecorder()
	return services, rr, nil
}

// func setUserContext(r *http.Request, user *model.User) *http.Request {
// 	ctx := context.WithValue(r.Context(), KeyUser("UserData"), user)
// 	return r.WithContext(ctx)
// }
// func addContext(req *http.Request, key, value interface{}) *http.Request {
// 	ctx := context.WithValue(req.Context(), key, value)
// 	return req.WithContext(ctx)
// }
