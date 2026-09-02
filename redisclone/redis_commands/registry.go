package rediscommands

import (
	"Pokemonscraper/redisclone/Datastructure"
	"Pokemonscraper/redisclone/resp"

	"strings"
	"sync"
)

// redis uses a single thread to handle any conncurrent opitions
// In - Memory Storage we will use for commands like hash, set
// we will make sure that only in order go routines will be able to access for it
var (
	cache     = make(map[string]resp.RespValue)
	cacheLock = sync.RWMutex{}                 //the sync.RwMutex keeps safe for concurrent go routines
	HSET      = map[string]map[string]string{} // a hash is map key string
)

type HandlerFunc func(s *Datastructure.Store, args []resp.RespValue) resp.RespValue

// note this will take your commands over resp
// We define this mapping using a simple dictionary:
var Handlers = map[string]func([]resp.RespValue) resp.RespValue{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"DEL":     del,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"ECHO":    echo,
	"HLen":    hlen,
	"Hexists": hexsits,
}

func Execute(s *Datastructure.Store, args []resp.RespValue) resp.RespValue {
	if len(args) == 0 {
		return resp.NewError("Error empty command")
	}
	value := strings.ToUpper(args[0].Bulk)
	// refer to the memory store
	val, ok := cache[value]
	if !ok {
		return resp.NewError("ERR unknown command '%s'")
	}
	return val

}

// Ping will return a pong if there is no argument,
// if the argument is provided then return argument in a bulk string

func ping(args []resp.RespValue) resp.RespValue {
	if len(args) == 0 {
		return resp.NewRespString("PONG")
	}
	return resp.NewBulkString(args[0].Bulk)
} // Set commands in Redis is a key-value pair
// Set and Get use a hashmap in golang to connect. For a set it will stored key to some value, if the key holds a value then it will be overwritten
// Ok response if the set was successfully executed, otherwise it will return an error
func set(args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("Err Wrong number for SET ")
	}
	key := strings.TrimSpace(args[0].Bulk)
	value := args[1]
	// here we lock the go routine
	cacheLock.Lock() // locks the go routine and wait for the unlock
	defer cacheLock.Unlock()
	cache[key] = value
	return resp.NewRespString("OK")
}
func get(args []resp.RespValue) resp.RespValue {
	if len(args) != 1 {
		return resp.NewError("Missing Values for correct GET response ")
	}
	key := strings.TrimSpace(args[0].Bulk)
	// lock the go routines now and then unlock once they've been proccessed
	cacheLock.RLock()
	// close the connection
	defer cacheLock.RUnlock()

	if val, ok := cache[key]; ok {
		return val
	}
	return resp.NewNull()
}

// delete will remove any specified keys. W
func del(args []resp.RespValue) resp.RespValue {
	if len(args) == 0 {
		return resp.NewError("Not enough arguemnts for DEL")
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()

	removed := 0
	for _, arg := range args {
		key := strings.TrimSpace(arg.Bulk)
		if _, ok := cache[key]; ok {
			delete(cache, key)
			removed++ // incerment the
		}
	}

	return resp.NewIntger(removed)
}

// here we will build hashget and hashset, a
//Redis Hash basic commands
//HSET: Sets the value of one or more fields on a hash.
//HGET: Returns the value at a given field.

func hset(args []resp.RespValue) resp.RespValue {
	if len(args) != 3 {
		return resp.NewError("ERR of wrong arguments for hset command")
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk
	cacheLock.Lock()
	if _, ok := HSET[hash]; !ok {
		HSET[hash] = map[string]string{}
	}
	HSET[hash][key] = value
	cacheLock.Unlock()
	return resp.RespValue{}

}
func hget(args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("ERR of wrong arguments for hget command")
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	cacheLock.Lock()

	_, ok := HSET[hash][key]
	cacheLock.RUnlock() // the RUnlock
	if !ok {
		return resp.RespValue{}
	}
	return resp.RespValue{}
}

// return the fields and values in hash
func hgetall(args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("ERR of wrong arguments for hgetall command")
	}
	key := args[1].Bulk
	cacheLock.Unlock()
	// once were done make sure to return the all hash key/elements at the end
	out := make([]resp.RespValue, 0, len(key)*2)
	return resp.NewArray(out)
}

// return the len of the current hash, the current hash will have current elements in the hash
func hlen(args []resp.RespValue) resp.RespValue {
	if len(args) != 2 {
		return resp.NewError("ERR of wrong arguments for hlen command")
	}
	cacheLock.Lock()
	hashlen := args[1].Bulk
	out := len(hashlen)
	return resp.NewIntger(out)

}

// here we check the hash even exisits
func hexsits(args []resp.RespValue) resp.RespValue {
	if len(args) != 3 {
		return resp.NewError("ERR of wrong arguments for hesits command")
	}
	return resp.NewIntger(1)
}

// echo handler
func echo(args []resp.RespValue) resp.RespValue {
	// echo handles two args: cmd + message
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for echo command")
	}
	return resp.NewBulkString(args[1].Bulk)
}
