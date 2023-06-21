package database

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/globals"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/types"
	"github.com/lib/pq"
)

func GetMetaBinariesByDataIDs(dataIDs []uint64) []*types.MetaBinary {
	metaBinaries := make([]*types.MetaBinary, 0, len(dataIDs))

	rows, err := Postgres.Query(`
		SELECT
		data_id,
		owner_pid,
		name,
		data_type,
		meta_binary,
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
		FROM meta_binaries WHERE data_id=ANY($1) ORDER BY updated_time DESC`,
		pq.Array(dataIDs),
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return metaBinaries
	}

	for rows.Next() {
		metaBinary := types.NewMetaBinary()

		metaBinary.CreationTime = nex.NewDateTime(0)
		metaBinary.UpdatedTime = nex.NewDateTime(0)
		metaBinary.ReferredTime = nex.NewDateTime(0)
		metaBinary.ExpireTime = nex.NewDateTime(0)

		var creationTimestamp int64
		var updatedTimestamp int64
		var referredTimestamp int64
		var expireTimestamp int64

		err := rows.Scan(
			&metaBinary.DataID,
			&metaBinary.OwnerPID,
			&metaBinary.Name,
			&metaBinary.DataType,
			&metaBinary.Buffer,
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
			_ = metaBinary.CreationTime.FromTimestamp(time.Unix(creationTimestamp, 0))
			_ = metaBinary.UpdatedTime.FromTimestamp(time.Unix(updatedTimestamp, 0))
			_ = metaBinary.ReferredTime.FromTimestamp(time.Unix(referredTimestamp, 0))
			_ = metaBinary.ExpireTime.FromTimestamp(time.Unix(expireTimestamp, 0))
		}

		metaBinaries = append(metaBinaries, metaBinary)
	}

	return metaBinaries
}
