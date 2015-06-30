package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type task struct {
  Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Name     string        `json:"name"`
  Category string        `json:"category"`
}

type tasksCollection struct {
  Data []task `json:"data"`
}

type TaskResource struct {
  Data task `json:"data"`
}

type TaskRepo struct {
  Coll *mgo.Collection
}

func (r *TaskRepo) All() (tasksCollection, error) {
  result := tasksCollection{[]task{}}
  err := r.Coll.Find(nil).All(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *TaskRepo) Find(id string) (TaskResource, error) {
  result := TaskResource{}
  err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *TaskRepo) Create(task *task) error {
  id := bson.NewObjectId()
  _, err := r.Coll.UpsertId(id, task)
  if err != nil {
    return err
  }

  task.Id = id

  return nil
}

func (r *TaskRepo) Update(task *task) error {
  err := r.Coll.UpdateId(task.Id, task)
  if err != nil {
    return err
  }

  return nil
}

func (r *TaskRepo) Delete(id string) error {
  err := r.Coll.RemoveId(bson.ObjectIdHex(id))
  if err != nil {
    return err
  }

  return nil
}