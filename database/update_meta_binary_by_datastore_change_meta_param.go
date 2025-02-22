package database

import (
	"time"

	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
)

func UpdateMetaBinaryByDataStoreChangeMetaParam(dataStoreChangeMetaParam datastore_types.DataStoreChangeMetaParam) error {
	// TODO - Check DataStoreChangeMetaParam flags for changes

	now := time.Now().Unix()

	tags := make([]string, len(dataStoreChangeMetaParam.Tags))
	for i, tag := range dataStoreChangeMetaParam.Tags {
		tags[i] = string(tag)
	}

	_, err := Postgres.Exec(`
		UPDATE meta_binaries SET (
			name,
			data_type,
			meta_binary,
			period,
			tags,
			updated_time
		) = ($1, $2, $3, $4, $5, $6) WHERE data_id=$7`,
		dataStoreChangeMetaParam.Name,
		dataStoreChangeMetaParam.DataType,
		dataStoreChangeMetaParam.MetaBinary,
		dataStoreChangeMetaParam.Period,
		pq.Array(tags),
		now,
		uint32(dataStoreChangeMetaParam.DataID),
	)

	return err
}
