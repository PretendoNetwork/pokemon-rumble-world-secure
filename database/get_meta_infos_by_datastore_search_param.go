package database

import (
	"database/sql"
	"time"

	nex_types "github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	"github.com/PretendoNetwork/pokemon-rumble-world/types"
	"github.com/lib/pq"
)

func GetMetaInfosByDataStoreSearchParam(param datastore_types.DataStoreSearchParam) []*types.MetaBinary {
	metaBinaries := make([]*types.MetaBinary, 0, param.ResultRange.Length)

	rows, err := Postgres.Query(`
		SELECT
		data_id,
		owner_pid,
		name,
		data_type,
		permission,
		del_permission,
		flag,
		period,
		tags,
		persistence_slot_id,
		extra_data,
		creation_time,
		updated_time,
		referred_time,
		expire_time
		FROM meta_binaries WHERE data_type=ANY($1) ORDER BY updated_time DESC LIMIT $2`,
		pq.Array(param.DataTypes),
		param.ResultRange.Length,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return metaBinaries
	}

	for rows.Next() {
		metaBinary := types.NewMetaBinary()

		metaBinary.CreationTime = nex_types.NewDateTime(0)
		metaBinary.UpdatedTime = nex_types.NewDateTime(0)
		metaBinary.ReferredTime = nex_types.NewDateTime(0)
		metaBinary.ExpireTime = nex_types.NewDateTime(0)

		var creationTimestamp int64
		var updatedTimestamp int64
		var referredTimestamp int64
		var expireTimestamp int64

		err := rows.Scan(
			&metaBinary.DataID,
			&metaBinary.OwnerPID,
			&metaBinary.Name,
			&metaBinary.DataType,
			&metaBinary.Permission,
			&metaBinary.DeletePermission,
			&metaBinary.Flag,
			&metaBinary.Period,
			pq.Array(&metaBinary.Tags),
			&metaBinary.PersistenceSlotID,
			pq.Array(&metaBinary.ExtraData),
			&creationTimestamp,
			&updatedTimestamp,
			&referredTimestamp,
			&expireTimestamp,
		)

		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		if err == nil {
			metaBinary.CreationTime.FromTimestamp(time.Unix(creationTimestamp, 0))
			metaBinary.UpdatedTime.FromTimestamp(time.Unix(updatedTimestamp, 0))
			metaBinary.ReferredTime.FromTimestamp(time.Unix(referredTimestamp, 0))
			metaBinary.ExpireTime.FromTimestamp(time.Unix(expireTimestamp, 0))
		}

		metaBinaries = append(metaBinaries, metaBinary)
	}

	return metaBinaries
}
