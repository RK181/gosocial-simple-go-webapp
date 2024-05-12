package models

import (
	"module/utils"
)

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

// Modelo que representa un usuario
type User struct {
	ID              int    `storm:"id,increment"`
	SessionToken    []byte `storm:"index"`
	Username        string
	Password        []byte
	Salt            []byte
	Email           string `storm:"unique"`
	PhoneNumber     string
	ShowPhoneNumber bool
	ShowEmail       bool
}

func GetUserByEmail(email string) (*User, error) {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return nil, err
	}
	defer DBConn.Close()
	var user User
	// Recuperamos al Usuario de la Base de Datos
	err = DBConn.One("Email", email, &user)
	// Comprobamos el exito
	return &user, err
}

/*
GetUserBySessionToken - Funcion que permite recuperar un usuario por su token de sesion
*/
func (u *User) GetUserBySessionToken(userSessionToken string) bool {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return false
	}
	defer DBConn.Close()
	// Recuperamos al Usuario de la Base de Datos
	err = DBConn.One("SessionToken", utils.Decode64(userSessionToken), u)
	// Comprobamos el exito
	return err == nil
}

/*
GetUserByID - Funcion que permite recuperar un usuario por su ID
*/
func (u *User) GetUserByID(id int) error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()
	// Recuperamos al Usuario de la Base de Datos
	err = DBConn.One("ID", id, u)
	// Comprobamos el exito
	return err
}

/*
RegisterUser - Funcion que permite registrar a un usuario en la Base de Datos
*/
func (u *User) RegisterUser() error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	// Creamos al Usuario
	// Generamos sal de 128 bits
	salt := utils.GenRandByteSlice(16)
	// Aplicamos PBKDF+Sal
	passArgon2Salt := utils.ApplyArgon2Salt(u.Password, salt)
	// Guardamos la contraseña y la sal
	u.Password = passArgon2Salt
	u.Salt = salt
	u.ShowEmail = false
	u.ShowPhoneNumber = false

	// Insertamos al Usuario en la Base de Datos
	err = DBConn.Save(u)
	return err
}

/*
LoginUser - Funcion que permite loguear a un usuario generando un token de sesion
*/
func (u *User) LoginUser() error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	password := u.Password
	// Recuperamos al Usuario de la Base de Datos
	err = DBConn.One("Email", u.Email, u)
	// Comprobamos el exito
	if err != nil {
		return err
	}

	// Comprobamos la contraseña
	if !utils.CheckArgon2Salt(password, u.Salt, u.Password) {
		return &UserError{"Invalid Password"}
	}

	// Generamos Token de Sesion
	token := utils.GenRandByteSlice(16)
	data := append([]byte(u.Email), token...)
	// Lo guardamos en el usuario
	u.SessionToken = utils.Hash256(data)

	// Lo actualizamos
	err = DBConn.Update(u)
	return err
}

/*
Logout - Funcion que permite desloguear a un usuario
*/
func (u *User) Logout() error {
	if u.ID > 0 {
		return &UserError{"User need ID to logout"}
	}

	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	// Limpiamos el token de sesion
	u.SessionToken = nil
	// Lo actualizamos
	err = DBConn.Update(u)
	return err
}

/*
UpdateUser - Funcion que permite actualizar un usuario en la Base de Datos.
Solo lo puede actualizar al propietario de la cuenta.
*/
func (u *User) UpdateUser() error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	// Actualizamos al Usuario en la Base de Datos
	err = DBConn.Update(u)
	return err
}

/*
DeleteUser - Funcion que permite eliminar un usuario de la Base de Datos.
El User tiene que tener un ID valido.
*/
func (u *User) DeleteUser() error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()
	// Eliminamos al Usuario de la Base de Datos
	err = DBConn.DeleteStruct(u)
	return err
}
