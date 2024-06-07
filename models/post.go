package models

import (
	"github.com/asdine/storm/v3/q"
)

// Modelo que representa un post
type Post struct {
	ID            int `storm:"id,increment"`
	UserID        int `storm:"index"`
	Title         string
	Content       string
	ForSubcribers bool
}

func NewPost() Post {
	return Post{}
}

/*
GetPostsByUserID - Funcion que permite recuperar los posts de un usuario por su ID.
*/
func (p Post) GetPostsByUserID(userID int, subscriber bool) ([]Post, error) {
	DBConn, err := dbConnect()
	if err != nil {
		return nil, err
	}
	defer DBConn.Close()

	var posts []Post
	err = DBConn.Select(q.And(q.Eq("UserID", userID), q.Or(q.Eq("ForSubcribers", false), q.Eq("ForSubcribers", subscriber)))).Reverse().Find(&posts)
	return posts, err
}

/*
GetAllPosts - Funcion que permite recuperar todos los posts de la Base de Datos.
*/
func (p Post) GetAllPosts(subsptions []UserUserSubscription) ([]Post, error) {
	DBConn, err := dbConnect()
	if err != nil {
		return nil, err
	}
	defer DBConn.Close()

	// Obtener los ID de los usuarios a los que se sigue
	var usersID []int
	for _, el := range subsptions {
		id := el.UserID
		usersID = append(usersID, id)
	}

	var posts []Post
	//err = DBConn.All(&posts, storm.Reverse())
	err = DBConn.Select(q.Or(q.Eq("ForSubcribers", false), q.In("UserID", usersID))).Reverse().Find(&posts)

	return posts, err
}

/*
GetPostByID - Funcion que permite recuperar un post por su ID.
*/
func (p *Post) GetPostByID(id int) error {
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	err = DBConn.One("ID", id, p)
	return err
}

/*
CreatePost - Funcion que permite crear un post en la Base de Datos.
*/
func (p *Post) CreatePost() error {
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	err = DBConn.Save(p)
	return err
}

/*
UpdatePost - Funcion que permite actualizar un post en la Base de Datos.
El Post tiene que tener un ID valido.
*/
func (p *Post) UpdatePost() error {
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	err = DBConn.Update(p)
	if err != nil {
		return err
	}
	// Actualizar los campos que no se actualizan con el Update
	err = DBConn.UpdateField(p, "ForSubcribers", p.ForSubcribers)
	return err
}

/*
DeletePost - Funcion que permite eliminar un post de la Base de Datos.
El Post iene que tener un ID valido.
*/
func (p *Post) DeletePost() error {
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	err = DBConn.DeleteStruct(p)
	return err
}
