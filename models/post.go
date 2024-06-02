package models

import "github.com/asdine/storm/v3"

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

func (p Post) GetPostsByUserID(userID int) ([]Post, error) {
	DBConn, err := dbConnect()
	if err != nil {
		return nil, err
	}
	defer DBConn.Close()

	var posts []Post
	err = DBConn.Find("UserID", userID, &posts, storm.Reverse())
	return posts, err
}

func (p Post) GetAllPosts() ([]Post, error) {
	DBConn, err := dbConnect()
	if err != nil {
		return nil, err
	}
	defer DBConn.Close()

	var posts []Post
	err = DBConn.All(&posts, storm.Reverse())
	return posts, err
}

func (p *Post) GetPostByID(id int) error {
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	err = DBConn.One("ID", id, p)
	return err
}

func (p *Post) CreatePost() error {
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	err = DBConn.Save(p)
	return err
}

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
