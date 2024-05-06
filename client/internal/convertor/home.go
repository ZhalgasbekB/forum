package convertor

import (
	"bytes"
	"encoding/json"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных полученных при запросе всех постов
func ConvertAllPosts(req *http.Request, resp *http.Response) (model.AllPostsConvertDate, error) {
	posts := model.AllPosts{}
	err := json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	convertDatePosts := make(model.AllPostsConvertDate, len(posts))
	for i := range posts {
		date := posts[i].CreatedAt
		formattedStr := date.Format("2006-01-02 15:04:05")
		convertDatePosts[i].PostID = posts[i].PostID
		convertDatePosts[i].UserID = posts[i].UserID
		convertDatePosts[i].CategoryName = posts[i].CategoryName
		convertDatePosts[i].Title = posts[i].Title
		convertDatePosts[i].Description = posts[i].Description
		convertDatePosts[i].CreatedAt = formattedStr
		convertDatePosts[i].Author = posts[i].Author
		convertDatePosts[i].Like = posts[i].Like
		convertDatePosts[i].Dislike = posts[i].Dislike
		convertDatePosts[i].PathImage = getPathImage(req, posts[i].PostID)
		cahe.PostImage[posts[i].PostID] = convertDatePosts[i].PathImage // Пишем кеш данных об картинках в постах
	}
	reverseSlice(convertDatePosts)
	return convertDatePosts, nil
}

func reverseSlice(s model.AllPostsConvertDate) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func getPathImage(r *http.Request, postId int) string {
	postsId := model.PostId{
		PostId: postId,
	}

	jsonData, err := json.Marshal(postsId)
	if err != nil {
		return ""
	}

	// Создание GET запроса на получение информации о картинке к новому посту
	req, err := http.NewRequest("GET", "http://localhost:8083/d3/upload-image", bytes.NewBuffer(jsonData))
	if err != nil {
		// errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при формировании GET запроса на получение данных об картинке к посту. Ошибка: %v", err)
		return ""
	}
	// // Добавление из браузера куки в запрос на сервер
	// req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
	req.Header.Set("Content-Type", "application/json")
	// Создаем структуру нового клиента
	client := http.Client{}
	// Отправляем запрос на сервер
	resp, err := client.Do(req)
	if err != nil {
		// errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при передаче GET запроса на сервис forum-api на получение данных об картинке к посту. Ошибка: %v", err)
		return ""
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		pathImagePost := &model.PathImagePost{}
		err := json.NewDecoder(resp.Body).Decode(pathImagePost)
		if err != nil {
			return ""
		}
		return pathImagePost.Path[1:]
	}
	return ""
}
