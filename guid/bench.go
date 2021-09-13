package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{
		//		StartTime: time.Now(),
	}
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func main() {
	for i := 0; i < 1000; i++ {
		fmt.Println(GenXID())
	}
	fmt.Printf("%x\n", GenSonyflake())
	fmt.Println(GenSonyflakeStr())
	fmt.Println(GenUUID())
}

func GenXID() string {
	return xid.New().String()
}

func GenSonyflake() uint64 {
	id, err := sf.NextID()
	if err != nil {
		return 0
	}
	return id
}

func GenSonyflakeStr() string {
	id, err := sf.NextID()
	if err != nil {
		return ""
	}
	idbin := make([]byte, 8)
	binary.BigEndian.PutUint64(idbin, id)
	return base32.StdEncoding.EncodeToString(idbin)
}

func GenUUID() string {
	return uuid.NewString()
}

const length = 16

func GenToken() string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
