package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := appContext{session.DB("test")}
	authHandler := authenticationHandler("username", "password")
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler, acceptHandler, authHandler)
	router := NewRouter()
	router.Get("/tasks/:id", commonHandlers.ThenFunc(appC.taskHandler))
	router.Put("/tasks/:id", commonHandlers.Append(contentTypeHandler, bodyHandler(taskResource{})).ThenFunc(appC.updatetaskHandler))
	router.Delete("/tasks/:id", commonHandlers.ThenFunc(appC.deletetaskHandler))
	router.Get("/tasks", commonHandlers.ThenFunc(appC.tasksHandler))
	router.Post("/tasks", commonHandlers.Append(contentTypeHandler, bodyHandler(taskResource{})).ThenFunc(appC.createtaskHandler))
	http.ListenAndServe(":8080", router)
}
