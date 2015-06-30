package main

import (
	"encoding/json"
	"net/http"

	models "github.com/abhiagarwal/go-exps/jsonapi-facebookauth-mongo/models"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type appContext struct {
	db *mgo.Database
}

func (c *appContext) register() {
	listTask := Actions{
		HandlerName: "list",
		OneHandler:  c.tasksHandler,
	}
	oneTask := Actions{
		HandlerName: "one",
		OneHandler:  c.taskHandler,
	}
	deleteTask := Actions{
		HandlerName: "delete",
		OneHandler:  c.deletetaskHandler,
	}
	updateTask := Actions{
		HandlerName: "put",
		OneHandler:  c.updatetaskHandler,
	}
	createTask := Actions{
		HandlerName: "post",
		OneHandler:  c.createtaskHandler,
	}
	RegisterAction("tasks", listTask, oneTask, updateTask, deleteTask, createTask)

	userUser := Actions{
		HandlerName: "list",
		OneHandler:  c.usersHandler,
	}
	oneUser := Actions{
		HandlerName: "one",
		OneHandler:  c.userHandler,
	}
	deleteUser := Actions{
		HandlerName: "delete",
		OneHandler:  c.deleteuserHandler,
	}
	updateUser := Actions{
		HandlerName: "put",
		OneHandler:  c.updateuserHandler,
	}
	createUser := Actions{
		HandlerName: "post",
		OneHandler:  c.createuserHandler,
	}
	RegisterAction("users", userUser, oneUser, updateUser, deleteUser, createUser)
}

func (c *appContext) usersHandler(w http.ResponseWriter, r *http.Request) {
	repo := models.UserRepo{c.db.C("users")}
	users, err := repo.All()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(users)
}

func (c *appContext) userHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := models.UserRepo{c.db.C("users")}
	user, err := repo.Find(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(user)
}

func (c *appContext) createuserHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*models.UserResource)
	repo := models.UserRepo{c.db.C("users")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (c *appContext) updateuserHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*models.UserResource)
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
	repo := models.UserRepo{c.db.C("users")}
	err := repo.Update(&body.Data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *appContext) deleteuserHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := models.UserRepo{c.db.C("users")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *appContext) tasksHandler(w http.ResponseWriter, r *http.Request) {
	repo := models.TaskRepo{c.db.C("tasks")}
	tasks, err := repo.All()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(tasks)
}

func (c *appContext) taskHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := models.TaskRepo{c.db.C("tasks")}
	task, err := repo.Find(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(task)
}

func (c *appContext) createtaskHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*models.TaskResource)
	repo := models.TaskRepo{c.db.C("tasks")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (c *appContext) updatetaskHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*models.TaskResource)
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
	repo := models.TaskRepo{c.db.C("tasks")}
	err := repo.Update(&body.Data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *appContext) deletetaskHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := models.TaskRepo{c.db.C("tasks")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}
