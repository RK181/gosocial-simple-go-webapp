package models

import (
	"github.com/asdine/storm/v3/q"
)

// Modelo que representa una suscripción a un usuario
type UserUserSubscription struct {
	ID          int `storm:"id,increment"`
	SubcriberID int `storm:"index"` // Usuario que se suscribe
	UserID      int `storm:"index"` // Usuario al que se suscribe
}

func (s *UserUserSubscription) SubscribeToUserByID(subcriberID int, userID int) error {
	// Comprobamos si ya existe la suscripcion
	if s.CheckSubscriptionbByUserID(userID, subcriberID) {
		return &SubscriptionError{"Ya estás suscrito a este usuario"}
	}

	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	// Creamos la suscripcion
	var subscription = UserUserSubscription{SubcriberID: subcriberID, UserID: userID}
	err = DBConn.Save(&subscription)
	// Comprobamos el exito
	return err
}

func (s UserUserSubscription) UnsubscribeToUserByID(subcriberID int, userID int) error {
	// Comprobamos si ya existe la suscripcion
	if !s.CheckSubscriptionbByUserID(userID, subcriberID) {
		return nil
	}
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	var subscription = UserUserSubscription{}
	err = DBConn.Select(q.And(q.Eq("UserID", userID), q.Eq("SubcriberID", subcriberID))).First(&subscription)
	if err != nil {
		return err
	}
	err = DBConn.DeleteStruct(&subscription)
	return err
}

func (s UserUserSubscription) CheckSubscriptionbByUserID(userID, subcriberID int) bool {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return false
	}
	defer DBConn.Close()

	// Recuperamos las suscripciones
	var subscriptions = UserUserSubscription{}
	err = DBConn.Select(q.And(q.Eq("UserID", userID), q.Eq("SubcriberID", subcriberID))).First(&subscriptions)
	// Comprobamos el exito
	return err == nil //|| userID == subcriberID
}

type SubscriptionError struct {
	Message string
}

func (e *SubscriptionError) Error() string {
	return e.Message
}
