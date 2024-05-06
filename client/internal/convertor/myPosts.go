package convertor

import (
	"encoding/json"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных полученных при запросе всех постов
func ConvertMyPosts(resp *http.Response) (model.AllPostsConvertDate, error) {
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
		convertDatePosts[i].PathImage = cahe.PostImage[posts[i].PostID]
	}
	return convertDatePosts, nil
}
