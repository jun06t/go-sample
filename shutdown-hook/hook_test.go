package hook

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestInvoke(t *testing.T) {
	Add(func(ctx context.Context) {
		fmt.Println("start hook1")
		time.Sleep(3 * time.Second)
		fmt.Println("end hook1")
	})
	Add(func(ctx context.Context) {
		fmt.Println("start hook2")
		time.Sleep(3 * time.Second)
		fmt.Println("end hook2")
	})
	Add(func(ctx context.Context) {
		fmt.Println("start hook3")
		time.Sleep(10 * time.Second)
		fmt.Println("end hook3")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Invoke(ctx)
	t.Log(err)
}
