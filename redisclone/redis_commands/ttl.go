package rediscommands

import (
	"redisclone/Datastructure"
	"redisclone/resp"
	"strconv"
	"time"
)

// a time to live gives a countdown to certain keys
// if I need a recap, just note: 0 seconds unti; expiry, -1, key exisits but no TTL its permanent, -2 key does not exisit

// since were outside the registry.go file we will intialiaze the handler we wrote
// Just to make it easier: because we will have one long file
// Note: we made a dict for all the commands, so ill pass the handlers and wrap it
// var Handlers map[string]func([]resp.RespValue) resp.RespValue

func expireHandler(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 3 {
		return resp.NewError("Err of wrong arguments for expire")
	}
	// for time to live keys we need to understand integer parsing
	// we need to ensure what is the expirey time for each key
	// Note for each key
	tkeys, err := strconv.Atoi(args[2].Bulk)
	if err != nil {
		return resp.NewError("Err value not intger or out of range ")
	}
	// if its not empty then we see if it matches with the args[1]
	if s.Expire(args[1].Bulk, time.Duration(tkeys)*time.Second) {
		return resp.NewIntger(1)
	}
	return resp.NewIntger(0)

}

// make a function to get the remaining time to live of the key
func ttlHandler(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Wrong arguments for 'ttl command'")
	}
	return resp.NewIntger(1)
}

func persistHandler(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments")
	}
	if s.Presist(args[1].Bulk) {
		return resp.NewIntger(1)
	}
	return resp.RespValue{}
}
