package convertor

import (
	"encoding/json"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/model"
	"lzhuk/clients/pkg/helpers"
	"net/http"
)

func ConvertGetPosts(r *http.Request, resp *http.Response) (*model.GetPostDate, error) {
	getPosts := &model.GetPost{}
	err := json.NewDecoder(resp.Body).Decode(getPosts)
	if err != nil {
		return &model.GetPostDate{}, err
	}

	nicname := cahe.Username[r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Value]

	// Преобразование даты в нужный формат
	formattedStr := getPosts.Post.CreateDate.Format("2006-01-02 15:04:05")

	// Создание структуры для преобразованных данных
	convertDatePosts := &model.GetPostDate{
		Post: &model.PostDate{
			PostId:       getPosts.Post.PostId,
			UserId:       getPosts.Post.UserId,
			CategoryName: getPosts.Post.CategoryName,
			Title:        getPosts.Post.Title,
			Description:  getPosts.Post.Description,
			CreateDate:   formattedStr,
			Author:       getPosts.Post.Author,
			Like:         getPosts.Post.Like,
			Dislike:      getPosts.Post.Dislike,
			PathImage:    cahe.PostImage[getPosts.Post.PostId],
		},
		Comment: make([]*model.CommentDate, len(getPosts.Comments)),
	}

	// Преобразование комментариев
	for i := range convertDatePosts.Comment {
		autorComment := 0
		date := getPosts.Comments[i].CreatedDate
		formattedStr := date.Format("2006-01-02 15:04:05")
		if nicname == getPosts.Comments[i].Name {
			autorComment = 1
		}
		convertDatePosts.Comment[i] = &model.CommentDate{
			ID:           getPosts.Comments[i].ID,
			User:         getPosts.Comments[i].User,
			Post:         getPosts.Comments[i].Post,
			Description:  getPosts.Comments[i].Description,
			CreatedDate:  formattedStr,
			UpdatedDate:  formattedStr,
			Name:         getPosts.Comments[i].Name,
			Like:         getPosts.Comments[i].Like,
			Dislike:      getPosts.Comments[i].Dislike,
			AutorComment: autorComment,
		}
	}
	return convertDatePosts, nil
}

func ConvertGetPostsNoAuth(r *http.Request, resp *http.Response) (*model.GetPostDate, error) {
	getPosts := &model.GetPost{}
	err := json.NewDecoder(resp.Body).Decode(getPosts)
	if err != nil {
		return &model.GetPostDate{}, err
	}

	// Преобразование даты в нужный формат
	formattedStr := getPosts.Post.CreateDate.Format("2006-01-02 15:04:05")

	// Создание структуры для преобразованных данных
	convertDatePosts := &model.GetPostDate{
		Post: &model.PostDate{
			PostId:       getPosts.Post.PostId,
			UserId:       getPosts.Post.UserId,
			CategoryName: getPosts.Post.CategoryName,
			Title:        getPosts.Post.Title,
			Description:  getPosts.Post.Description,
			CreateDate:   formattedStr,
			Author:       getPosts.Post.Author,
			Like:         getPosts.Post.Like,
			Dislike:      getPosts.Post.Dislike,
			PathImage:   cahe.PostImage[getPosts.Post.PostId],
		},
		Comment: make([]*model.CommentDate, len(getPosts.Comments)),
	}
	// Преобразование комментариев
	for i := range convertDatePosts.Comment {
		autorComment := 0
		date := getPosts.Comments[i].CreatedDate
		formattedStr := date.Format("2006-01-02 15:04:05")
		convertDatePosts.Comment[i] = &model.CommentDate{
			ID:           getPosts.Comments[i].ID,
			User:         getPosts.Comments[i].User,
			Post:         getPosts.Comments[i].Post,
			Description:  getPosts.Comments[i].Description,
			CreatedDate:  formattedStr,
			UpdatedDate:  formattedStr,
			Name:         getPosts.Comments[i].Name,
			Like:         getPosts.Comments[i].Like,
			Dislike:      getPosts.Comments[i].Dislike,
			AutorComment: autorComment,
		}
	}
	return convertDatePosts, nil
}
