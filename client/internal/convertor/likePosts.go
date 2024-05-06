package convertor

import (
	"encoding/json"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных для получения всех понравишихся постов пользователя
func ConvertUserLikePosts(resp *http.Response) ([]model.PostDate, error) {
	postsLikeUser := []model.Post{}
	err := json.NewDecoder(resp.Body).Decode(&postsLikeUser)
	if err != nil {
		return nil, err
	}
	convertDatePosts := make([]model.PostDate, len(postsLikeUser))
	for i := range postsLikeUser {
		date := postsLikeUser[i].CreateDate
		formattedStr := date.Format("2006-01-02 15:04:05")
		convertDatePosts[i].PostId = postsLikeUser[i].PostId
		convertDatePosts[i].UserId = postsLikeUser[i].UserId
		convertDatePosts[i].CategoryName = postsLikeUser[i].CategoryName
		convertDatePosts[i].Title = postsLikeUser[i].Title
		convertDatePosts[i].Description = postsLikeUser[i].Description
		convertDatePosts[i].CreateDate = formattedStr
		convertDatePosts[i].Author = postsLikeUser[i].Author
		convertDatePosts[i].Like = postsLikeUser[i].Like
		convertDatePosts[i].Dislike = postsLikeUser[i].Dislike
		convertDatePosts[i].PathImage = cahe.PostImage[postsLikeUser[i].PostId]
	}
	reverseLikePost(convertDatePosts)
	return convertDatePosts, nil
}

func reverseLikePost(s []model.PostDate) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
