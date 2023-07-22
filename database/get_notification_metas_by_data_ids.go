package database

import (
	"database/sql"
	"fmt"
	"os"

	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	"github.com/lib/pq"
)

func GetNotificationMetasByDataIDs(dataIDs []uint32) []*datastore_types.DataStoreSpecificMetaInfoV1 {
	notificationMetas := make([]*datastore_types.DataStoreSpecificMetaInfoV1, 0, len(dataIDs))

	rows, err := Postgres.Query(`
		SELECT
		data_id,
		owner_pid,
		data_type
		FROM notification_metas WHERE data_id=ANY($1)`,
		pq.Array(dataIDs),
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return notificationMetas
	}

	for rows.Next() {
		notificationMeta := datastore_types.NewDataStoreSpecificMetaInfoV1()

		err := rows.Scan(
			&notificationMeta.DataID,
			&notificationMeta.OwnerID,
			&notificationMeta.DataType,
		)

		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		if err == nil {
			bucket := os.Getenv("PN_PRW_S3_BUCKET")
			key := fmt.Sprintf("data/%011d", notificationMeta.DataID)

			size, err := globals.S3ObjectSize(bucket, key)
			if err != nil {
				globals.Logger.Error(err.Error())
				continue
			}

			notificationMeta.Size = uint32(size)

			notificationMetas = append(notificationMetas, notificationMeta)
		}
	}

	return notificationMetas
}
