package rediscommands

import (
	"Pokemonscraper/redisclone/Datastructure"
	"Pokemonscraper/redisclone/resp"
	"strconv"
)

// Left push must have a key, value pair [value]
// Moves more values to the left the head of your list
// develop a list an we will pop from the left, it follows the que first in/first out
func lPush(_ *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) < 3 {
		return resp.NewError("Error wrong number of arguments for LeftPush command")
	}
	// we push into the bulk string
	values := make([]string, len(args)-2)
	// iterate through the args
	for k, v := range args[2:] {
		values[k] = v.Bulk
	}
	// implement the que method we can also use the time to live key
	// q, err := s.LPush(args[1].Bulk, values) we can use this just replace the _ to s points to the struct
	q, err := args[1].Bulk, values
	if err != nil {
		return resp.NewError("error occured wrong type")
	}
	return resp.NewBulkString(q)
}

// Rpush the key value, we will append to the tail of the list
func rPush(_ *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) < 3 {
		return resp.NewError("Error wrong number of arguments for RightPush command")
	}
	// Apply them to a bulk string
	values := make([]string, len(args)-2)
	for k, v := range args[2:] {
		values[k] = v.Bulk
		// we will add the values to the tail
		// l, err := s.RPush(args.Bulk,values)
		//t == tail
	}
	tail, err := args[1].Bulk, values
	if err != nil {
		return resp.NewError("error occured wrong type")
	}
	return resp.NewBulkString(tail)
}

// now we follow the pop methods which means we remove from the list/que
func lPop(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Error wrong number of arguments for Left Pop command")
	}
	val, found, err := s.LPop(args[1].Bulk)
	if err != nil {
		return resp.NewError("error type")
	}
	// if the element does not exisit
	if !found {
		return resp.NewNull()
	}
	return resp.NewBulkString(val)

}
func rPop(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Error wrong number of arguments for Right Pop command")
	}
	val, found, err := s.RPop(args[1].Bulk)
	if err != nil {
		return resp.NewError("error type")
	}
	// if the element does not exisit
	if !found {
		return resp.NewNull()
	}
	return resp.NewBulkString(val)
}

// find the number of elements in the list, and check the edge case
func lLen(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Error wrong number of argumets for LLen command")
	}
	elem, err := s.LLen(args[1].Bulk)
	if err != nil {
		return resp.NewError("Error type")
	}
	// follow redis logic about returning 0 if the current element does not exisit
	if !elem {
		return resp.NewIntger(0)
	}
	return resp.NewIntger(elem)
}

// We return the a sub list from index start

func lRange(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) != 4 {
		return resp.NewError("Error wrong number of argumnets for LLRange command")
	}
	// parse the in line intger
	start, err := strconv.Atoi(args[2].Bulk)
	if err != nil {
		return resp.NewError("Element is out of range")
	}
	end, err := strconv.Atoi(args[3].Bulk)
	if err != nil {
		return resp.NewError("Element is out of range")
	}
	elem, err := s.LRange(args[1].Bulk, start, end)
	if err != nil {
		return resp.NewError("Element does not exsit")
	}
	// make the line list
	output := make([]resp.RespValue, len(elem))
	for k, v := range elem {
		output[k] = resp.NewBulkString(v)
	}
	return resp.NewArray(output)
}
