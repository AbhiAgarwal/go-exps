package main

import (
  "encoding/json"
  "net/http"

  tasksModel "github.com/abhiagarwal/go-exps/jsonapi-template/models/models"

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
    OneHandler: c.tasksHandler,
  }
  oneTask := Actions{
    HandlerName: "one",
    OneHandler: c.taskHandler,
  }
  deleteTask := Actions{
    HandlerName: "delete",
    OneHandler: c.deletetaskHandler,
  }
  updateTask := Actions{
    HandlerName: "put",
    OneHandler: c.updatetaskHandler,
  }
  createTask := Actions{
    HandlerName: "post",
    OneHandler: c.createtaskHandler,
  }
  RegisterAction("tasks", listTask, oneTask, updateTask, deleteTask, createTask)
}

func (c *appContext) tasksHandler(w http.ResponseWriter, r *http.Request) {
  repo := tasksModel.taskRepo{c.db.C("tasks")}
  tasks, err := repo.All()
  if err != nil {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/vnd.api+json")
  json.NewEncoder(w).Encode(tasks)
}

func (c *appContext) taskHandler(w http.ResponseWriter, r *http.Request) {
  params := context.Get(r, "params").(httprouter.Params)
  repo := tasksModel.taskRepo{c.db.C("tasks")}
  task, err := repo.Find(params.ByName("id"))
  if err != nil {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/vnd.api+json")
  json.NewEncoder(w).Encode(task)
}

func (c *appContext) createtaskHandler(w http.ResponseWriter, r *http.Request) {
  body := context.Get(r, "body").(*taskResource)
  repo := tasksModel.taskRepo{c.db.C("tasks")}
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
  body := context.Get(r, "body").(*taskResource)
  body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
  repo := tasksModel.taskRepo{c.db.C("tasks")}
  err := repo.Update(&body.Data)
  if err != nil {
    panic(err)
  }

  w.WriteHeader(204)
  w.Write([]byte("\n"))
}

func (c *appContext) deletetaskHandler(w http.ResponseWriter, r *http.Request) {
  params := context.Get(r, "params").(httprouter.Params)
  repo := tasksModel.taskRepo{c.db.C("tasks")}
  err := repo.Delete(params.ByName("id"))
  if err != nil {
    panic(err)
  }

  w.WriteHeader(204)
  w.Write([]byte("\n"))
}