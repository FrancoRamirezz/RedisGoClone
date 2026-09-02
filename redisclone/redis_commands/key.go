package rediscommands

import (
	"redisclone/Datastructure"
	"redisclone/resp"
)

// key type comes from the keys store
// return the data type stored in a simple string
func typehandler(args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Error wrong number of argumets for type command")
	}
	return resp.NewRespString(args[1].Bulk)
}

// check the key pattern that macthes pattern
func keyhandler(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Error wrong numbers of arguments for key command")
	}
	// we check the key and values pair
	keys := args[1].Bulk
	val := make([]resp.RespValue, len(keys))
	for index, value := range keys {
		val[index] = resp.NewBulkString(value)
	}
	return resp.NewArray(val)
}
func remaininghandler(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 3 {
		return resp.NewError("Error wrong number of rename command")
	}
	_, err := s.remainingKey(args[1].Bulk, args[2].Bulk)
	if err != nil {
		return resp.NewError("Wrong type of commands")
	}
	return resp.NewRespString("OK")
}

// flush the
func flushHandler(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	s.Flushkeys()

}
