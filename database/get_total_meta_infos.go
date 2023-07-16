package database

import "github.com/PretendoNetwork/pokemon-rumble-world/globals"

func GetTotalMetaInfos() uint32 {
	var total uint32

	err := Postgres.QueryRow(`SELECT COUNT(*) FROM meta_binaries`).Scan(&total)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return total
}
