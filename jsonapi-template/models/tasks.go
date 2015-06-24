package tasks

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

type taskResource struct {
  Data task `json:"data"`
}

type taskRepo struct {
  coll *mgo.Collection
}

func (r *taskRepo) All() (tasksCollection, error) {
  result := tasksCollection{[]task{}}
  err := r.coll.Find(nil).All(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *taskRepo) Find(id string) (taskResource, error) {
  result := taskResource{}
  err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *taskRepo) Create(task *task) error {
  id := bson.NewObjectId()
  _, err := r.coll.UpsertId(id, task)
  if err != nil {
    return err
  }

  task.Id = id

  return nil
}

func (r *taskRepo) Update(task *task) error {
  err := r.coll.UpdateId(task.Id, task)
  if err != nil {
    return err
  }

  return nil
}

func (r *taskRepo) Delete(id string) error {
  err := r.coll.RemoveId(bson.ObjectIdHex(id))
  if err != nil {
    return err
  }

  return nil
}