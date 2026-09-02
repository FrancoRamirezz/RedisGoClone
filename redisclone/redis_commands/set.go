package rediscommands

import (
	"redisclone/Datastructure"
	"redisclone/resp"
)

// the resp value will return any og the
func sadd(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) < 3 {
		return resp.NewError("Error wrong number of arguments for sadd commands")
	}
	// pass the elements to the bulkargs
	elem := make([]string, len(args)-2)
	for index, val := range args[2:] {
		elem[index] = val.Bulk
	}
	added, err := s.SADD(args[1].Bulk, elem)
	if err != nil {
		return resp.NewError(err.Error())
	}
	return resp.NewIntger(added)
}

func srem(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) > 3 {
		return resp.NewError("Error wrong number of arguments for srem command")
	}
	elem := make([]string, len(args)-2)
	for index, val := range args[2:] {
		elem[index] = val.Bulk
	}
	removed, err := s.SRem(args[1].Bulk, elem)
	if err != nil {
		return resp.NewError(err.Error())
	}
	return resp.NewIntger(removed)
}

func smembers(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Err wrong number of arguments for smembers command")
	}
	elem, err := s.SMembers(args[1].Bulk)
	if err != nil {
		return resp.NewError(err.Error())
	}
	// we return the entire memebers into a bulk string
	Bstr := make([]resp.RespValue, len(elem))
	for mem, bulk := range elem {
		Bstr[mem] = resp.NewBulkString(bulk)
	}
	return resp.NewArray(Bstr)
}

// now we check if the members/elem are currently present
func sismember(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 3 {
		return resp.NewError("Error wrong number of arguments for sismember command")
	}
	exists, err := s.SIsMember(args[1].Bulk, args[2].Bulk)
	if err != nil {
		return resp.NewError(err.Error())
	}
	if exists {
		return resp.NewIntger(1)
	}
	return resp.NewIntger(0)
}

// we must return the number of members/elements in the set
// return zero if the key does not exsit
func scard(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Error wrong numbers of arguments for scard command")
	}
	elem, err := s.SCard(args[1].Bulk, args[2].Bulk)
	if err != nil {
		return resp.NewError(err.Error())
	}
	return resp.NewIntger(elem)

}
