package database

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
)

func GetNotificationsByPIDAndParam(pid uint32, param datastore_types.DataStoreGetNewArrivedNotificationsParam) types.List[datastore_types.DataStoreNotificationV1] {
	var notifications types.List[datastore_types.DataStoreNotificationV1] = make([]datastore_types.DataStoreNotificationV1, 0, param.Limit)

	rows, err := Postgres.Query(`
		SELECT
		notification_id,
		data_id
		FROM notifications WHERE pid=$1 ORDER BY notification_id DESC LIMIT $2`,
		pid,
		param.Limit,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return notifications
	}

	for rows.Next() {
		notification := datastore_types.NewDataStoreNotificationV1()

		err := rows.Scan(
			&notification.NotificationID,
			&notification.DataID,
		)

		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		if err == nil {
			// * Cleanup previous notifications already read
			if param.LastNotificationID == notification.NotificationID {
				notifications = make([]datastore_types.DataStoreNotificationV1, 0, param.Limit)
				continue
			}

			notifications = append(notifications, notification)
		}
	}

	return notifications
}
