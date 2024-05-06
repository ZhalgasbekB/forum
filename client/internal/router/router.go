package router

import (
	"lzhuk/clients/internal/controller"
	"lzhuk/clients/pkg/config"
	"net/http"
	"time"
)

func Router(cfg config.Config) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("/", controller.StartPage)
	router.HandleFunc("/userd3", controller.HomePage)
	router.HandleFunc("/register", controller.RegisterPage)
	router.HandleFunc("/login", controller.LoginPage)
	router.HandleFunc("/userd3/posts", controller.CreatePost)
	router.HandleFunc("/userd3/myposts", controller.MyPosts)
	router.HandleFunc("/userd3/post/", controller.GetPost)
	router.HandleFunc("/userd3/comments/", controller.Comments)
	router.HandleFunc("/userd3/createcomment", controller.CreateComment)
	router.HandleFunc("/userd3/updatepost/", controller.UpdatePost)
	router.HandleFunc("/userd3/deletepost/", controller.DeletePost)
	router.HandleFunc("/userd3/votepost", controller.VotePost)
	router.HandleFunc("/userd3/votecomment", controller.VoteComment)
	router.HandleFunc("/userd3/likeposts", controller.LikePost)
	router.HandleFunc("/userd3/updatecomment", controller.UpdateComment)
	router.HandleFunc("/userd3/deletecomment", controller.DeleteComment)
	router.HandleFunc("/logout", controller.LogoutUser)
	router.HandleFunc("/userd3/category", controller.CategoryPosts)

	
	router.HandleFunc("/google/login", controller.Google)
	router.HandleFunc("/google/callback", controller.GoogleCallback)
	router.HandleFunc("/github/login", controller.GitHub)
	router.HandleFunc("/github/callback", controller.GitHubCallback)
	router.HandleFunc("/yandex/auth", controller.YandexAuth)
	router.HandleFunc("/yandex/callback", controller.YandexCallback)

	fileServer := http.FileServer(http.Dir("ui/css"))
	router.Handle("/ui/css/", http.StripPrefix("/ui/css/", fileServer))

	fileServerImg := http.FileServer(http.Dir("ui/img"))
	router.Handle("/ui/img/", http.StripPrefix("/ui/img/", fileServerImg))

	fileServerImagesUser := http.FileServer(http.Dir("./uploads"))
	router.Handle("/uploads/", http.StripPrefix("/uploads/", fileServerImagesUser))

	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,  // время ожидания для чтения данных
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second, // время ожидания для записи данных
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,  // время простоя
		Handler:      router,
	}
	return s
}
