package convertor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных при создании нового поста
func ConvertCreatePost(r *http.Request) ([]byte, error) {
	createPost := model.CreatePost{
		CategoryName: r.FormValue("category"),
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
	}

	jsonData, err := json.Marshal(createPost)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func ConvertCreatePostImage(r *http.Request, postId int) ([]byte, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("Произошла ошибка при получении данных об картинке из формы запроса. Ошибка: %v", err)
		return []byte{}, err
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Произошла ошибка при чтении отправленной картинки пользователем. Ошибка: %v", err)
		return nil, err
	}

	if len(imageData) != 0 {
		err = ioutil.WriteFile("./uploads/"+handler.Filename, imageData, 0666)
		if err != nil {
			log.Printf("Произошла ошибка при записи на сервер картинки. Ошибка: %v", err)
			return nil, err
		}
		uploadPostImage := model.UploadPostImage{
			Path:   fmt.Sprintf("./uploads/" + handler.Filename),
			PostId: postId,
		}
		jsonData, err := json.Marshal(uploadPostImage)
		if err != nil {
			return nil, err
		}
		return jsonData, nil
	}
	return nil, err
}

func ConvertPostId(resp http.Response) (int, error) {
	postId := model.PostId{}
	if err := json.NewDecoder(resp.Body).Decode(&postId); err != nil {
		return 0, err
	}
	return postId.PostId, nil
}
