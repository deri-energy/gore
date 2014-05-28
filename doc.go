// Copyright 2014 keimoon. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Gore is a full feature Redis client for go:
  - Convenient command building and reply parsing
  - Pipeline, multi-exec, LUA scripting
  - Pubsub
  - Connection pool
  - Redis sentinel
  - Client implementation of sharding

Connections

Gore only supports TCP connection for Redis. The connection is thread-safe and can be auto-repaired
with or without sentinel.

  conn, err := gore.Dial("localhost:6379", 10 * time.Duration) //Connect to redis server at localhost:6379
  if err != nil {
    return
  }
  defer conn.Close()

Command

Redis command is built with NewCommand

  gore.NewCommand("SET", "kirisame", "marisa") // SET kirisame marisa
  gore.NewCommand("ZADD", "magician", 1337, "alice") // ZADD magician 1337 alice
  gore.NewCommand("HSET", "sdm", "sakuya", 99) // HSET smd sakuya 99

In the last command, the value stored by redis will be the string "99", not the integer 99.
  Integer and float values are converted to string using strconv
  Boolean values are convert to "1" and "0"
  Nil values are stored as zero length string
  Other types are converted to string using standard fmt.Sprint

To efficiently store integer, you can use gore.FixInt or gore.VarInt

Compact integer

Gore supports compacting integer to reduce memory used by redis. There are 2 ways of compacting integer:
  gore.FixInt stores an integer as a fixed 8 bytes []byte.
  gore.VarInt encodes an integer with variable length []byte.

  gore.NewCommand("SET", "fixint", gore.FixInt(1337)) // Set fixint as an 8 bytes []byte
  gore.NewCommand("SET", "varint", gore.VarInt(1337)) // varint only takes 3 bytes

Reply

A redis reply is return when the command is run on a connection

  rep, err := gore.NewCommand("GET", "kirisame").Run(conn)

Parsing the reply is straightforward:

  s, _ := rep.String()  // Return string value if reply is simple string (status) or bulk string
  b, _ := rep.Bytes()   // Return a byte array
  x, _ := rep.Integer() // Return integer value if reply type is integer (INCR, DEL)
  e, _ := rep.Error()   // Return error message if reply type is error
  a, _ := rep.Array()   // Return reply list if reply type is array (MGET, ZRANGE)

Reply converting

Reply support convenient methods to convert to other types

  x, _ := rep.Int()    // Convert string value to int64. This method is different from rep.Integer()
  f, _ := rep.Float()  // Convert string value to float64
  t, _ := rep.Bool()   // Convert string value to boolean, where "1" is true and "0" is false
  x, _ := rep.FixInt() // Convert string value to FixInt
  x, _ := rep.VarInt() // Convert string value to VarInt

To convert an array reply to a slice, you can use Slice method:

  s := []int
  err := rep.Slice(&s) // Convert an array reply to a slice of integer

Gore supports following slice element type:
  - integer (int, int64)
  - float (float64)
  - string and []byte
  - FixInt and VarInt
  - *gore.Pair for converting map data from HGETALL or ZRANGE WITHSCORES

*/
package gore
