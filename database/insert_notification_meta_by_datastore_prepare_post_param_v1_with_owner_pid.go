package database

import (
	"time"

	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/lib/pq"
)

func InsertNotificationMetaByDataStorePreparePostParamV1WithOwnerPID(dataStorePreparePostParam *datastore_types.DataStorePreparePostParamV1, pid uint32) (uint32, error) {
	var dataID uint32

	now := time.Now().Unix()
	expireTime := time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC).Unix()

	err := Postgres.QueryRow(`
		INSERT INTO notification_metas (
			owner_pid,
			name,
			data_type,
			permission,
			del_permission,
			flag,
			period,
			tags,
			creation_time,
			updated_time,
			referred_time,
			expire_time
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING data_id`,
		pid,
		dataStorePreparePostParam.Name,
		dataStorePreparePostParam.DataType,
		dataStorePreparePostParam.Permission.Permission,
		dataStorePreparePostParam.DelPermission.Permission,
		dataStorePreparePostParam.Flag,
		dataStorePreparePostParam.Period,
		pq.Array(dataStorePreparePostParam.Tags),
		now,
		now,
		now,
		expireTime,
	).Scan(&dataID)

	return dataID, err
}
