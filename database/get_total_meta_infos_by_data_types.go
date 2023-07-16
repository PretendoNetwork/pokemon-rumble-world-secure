package database

import (
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	"github.com/lib/pq"
)

func GetTotalMetaInfosByDataTypes(dataTypes []uint16) uint32 {
	var total uint32

	err := Postgres.QueryRow(`SELECT COUNT(*) FROM meta_binaries WHERE data_type=ANY($1)`, pq.Array(dataTypes)).Scan(&total)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return total
}
