package shardingmap

import (
	"fmt"
	"testing"

	"github.com/x-junkang/connected/internal/connect"
)

func TestBucketGet(t *testing.T) {
	bucket, err := NewBucket(16)
	if err != nil {
		t.Fatalf("new bucket fail, err is %s", err.Error())
		return
	}
	conn := connect.NewConnectionTcp(nil, nil, 14, nil)
	bucket.Set(16, conn)
	value, ok := bucket.Get(16)
	if !ok {
		t.Fatal()
	}
	if value.GetConnID() != 14 {
		t.Fatal("结果出错了")
	}
	fmt.Println(value.GetConnID())
}
