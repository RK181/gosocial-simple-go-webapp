package models

// Modelo que representa una suscripción a un usuario
type UserUserSubscription struct {
	ID            int `storm:"id,increment"`
	SubscriptorID int `storm:"index"` // Usuario que se suscribe
	UserID        int `storm:"index"` // Usuario al que se suscribe
}

func (s UserUserSubscription) SubscribeToUserByID(subscriptorID int, userID int) error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	// Comprobamos si ya existe la suscripcion
	var subscription UserUserSubscription
	err = DBConn.One("SubscriptorID", subscriptorID, &subscription)
	if err == nil {
		return &SubscriptionError{"Ya estás suscrito a este usuario"}
	}

	// Creamos la suscripcion
	subscription = UserUserSubscription{SubscriptorID: subscriptorID, UserID: userID}
	err = DBConn.Save(&subscription)
	// Comprobamos el exito
	return err
}

func (s UserUserSubscription) UnsubscribeToUserByID(subscriptorID int, userID int) error {
	// Conectamos a la Base de Datos
	DBConn, err := dbConnect()
	if err != nil {
		return err
	}
	defer DBConn.Close()

	// Eliminamos la suscripcion
	err = DBConn.DeleteStruct(&UserUserSubscription{SubscriptorID: subscriptorID, UserID: userID})
	// Comprobamos el exito
	return err
}

type SubscriptionError struct {
	Message string
}

func (e *SubscriptionError) Error() string {
	return e.Message
}
