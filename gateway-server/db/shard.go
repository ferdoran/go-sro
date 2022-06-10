package db

import "github.com/ferdoran/go-sro/framework/db"

type Shard struct {
	Id            int
	ContentId     int
	Name          string
	Capacity      int
	OnlinePlayers int
	Status        int
}

const (
	select_shards string = "SELECT ID, CONTENT_ID, NAME, CAPACITY, ONLINE_PLAYERS, STATUS FROM SHARD"
)

func GetShards() []Shard {
	conn := db.OpenConnAccount()
	defer conn.Close()

	queryHandle, err := conn.Query(select_shards)
	db.CheckError(err)

	var shards []Shard
	for queryHandle.Next() {
		var shardId, contentId, capacity, onlinePlayers, status int
		var name string
		err = queryHandle.Scan(&shardId, &contentId, &name, &capacity, &onlinePlayers, &status)
		db.CheckError(err)
		shards = append(shards, Shard{
			Id:            shardId,
			ContentId:     contentId,
			Name:          name,
			Capacity:      capacity,
			OnlinePlayers: onlinePlayers,
			Status:        status,
		})
	}
	return shards
}
