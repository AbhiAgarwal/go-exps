package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type User struct {
  Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  FirstName   string `json:"first_name"`
  Gender      string `json:"gender"`
  ID          string `json:"id"`
  LastName    string `json:"last_name"`
  Link        string `json:"link"`
  Locale      string `json:"locale"`
  Name        string `json:"name"`
  Timezone    int    `json:"timezone"`
  UpdatedTime string `json:"updated_time"`
  Verified    bool   `json:"verified"`
}

type usersCollection struct {
  Data []User `json:"data"`
}

type UserResource struct {
  Data User `json:"data"`
}

type UserRepo struct {
  Coll *mgo.Collection
}

func (r *UserRepo) All() (usersCollection, error) {
  result := usersCollection{[]User{}}
  err := r.Coll.Find(nil).All(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *UserRepo) Find(id string) (UserResource, error) {
  result := UserResource{}
  err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
  if err != nil {
    return result, err
  }

  return result, nil
}

func (r *UserRepo) Create(user *User) error {
  id := bson.NewObjectId()
  _, err := r.Coll.UpsertId(id, user)
  if err != nil {
    return err
  }

  user.Id = id

  return nil
}

func (r *UserRepo) Update(user *User) error {
  err := r.Coll.UpdateId(user.Id, user)
  if err != nil {
    return err
  }

  return nil
}

func (r *UserRepo) Delete(id string) error {
  err := r.Coll.RemoveId(bson.ObjectIdHex(id))
  if err != nil {
    return err
  }

  return nil
}