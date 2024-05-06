package cahe

import "lzhuk/clients/model"

var Username = make(map[string]string) // Хеш-таблица для хранение имени (никнейма) пользователя

var CategoryPosts model.AllPostsConvertDate

var PostImage = make(map[int]string) // Хеш-таблица для хранения данных о привязке постов и картинок к постам