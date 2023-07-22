package database

func InsertNotificationByDataIDAndPID(dataID, pid uint32) (uint64, error) {
	var notificationID uint64

	err := Postgres.QueryRow(`INSERT INTO notifications (data_id, pid) VALUES ($1, $2) RETURNING notification_id`, dataID, pid).Scan(&notificationID)

	return notificationID, err
}
