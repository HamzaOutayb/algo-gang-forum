package api

import (
	"database/sql"
	"net/http"
	"time"

	handler "real-time-forum/internal/api/handler"
	"real-time-forum/pkg/ratelimiter"
)

func Routes(db *sql.DB) *http.ServeMux {
	d := handler.NewHandler(db)
	mux := http.NewServeMux()

	FileServer := http.FileServer(http.Dir("./Assets/"))
	mux.Handle("/Assets/", http.StripPrefix("/Assets/", FileServer))
	mux.HandleFunc("/", handler.HomeHandler)
	loginRateLimiter := ratelimiter.LoginLimiter.RateMiddlewareAuth(http.HandlerFunc(d.Signin), 5, time.Minute)
	signupRateLimiter := ratelimiter.SignupLimiter.RateMiddlewareAuth(http.HandlerFunc(d.Signup), 5, time.Minute)
	mux.Handle("/Signin", loginRateLimiter)
	mux.Handle("/Signup", signupRateLimiter)
	mux.HandleFunc("/post", d.InsertPostsHandler)
	mux.HandleFunc("GET /post/{id}", d.GetPostByIdHandler)
	mux.HandleFunc("GET /post", d.GetPostHandler)

	addCommentHandler := ratelimiter.AddCommentsLimter.RateMiddleware(http.HandlerFunc(d.AddCommentHandler), 10, 2*time.Second, db)
	mux.Handle("/comment", addCommentHandler)

	reactionRateLimiter := ratelimiter.ReactionsLimiter.RateMiddleware(http.HandlerFunc(d.ReactionHandler), 10, 500*time.Millisecond, db)
	mux.Handle("/api/reaction", reactionRateLimiter)
	mux.HandleFunc("/chat", d.ChatService)
	mux.HandleFunc("/Lastconversation/", d.Lastconversation)

	go func() {
		for {
			time.Sleep(120 * time.Minute)
			ratelimiter.AddCommentsLimter.RemoveSleepUsers()
			ratelimiter.AddPostLimter.RemoveSleepUsers()
			ratelimiter.LoginLimiter.RemoveSleepUsers()
			ratelimiter.ReactionsLimiter.RemoveSleepUsers()
			ratelimiter.SignupLimiter.RemoveSleepUsers()
		}
	}()

	return mux
}
