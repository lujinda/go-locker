/*
	Author        : tuxpy
	Email         : q8886888@qq.com.com
	Create time   : 2017-04-24 23:31:03
	Filename      : test_locker.go
	Description   :
*/

package locker

import (
	"fmt"
	"testing"
	"time"
)

func TestLocker(t *testing.T) {
	l := New()
	fmt.Println("start lock 1")
	l.Lock("l1", 1)
	fmt.Println("end lock 1")

	time.AfterFunc(time.Second*3, func() {
		fmt.Println("start lock 1")
		l.Unlock("l1")
		fmt.Println("end  unlock 1")
	})

	fmt.Println("start lock 1")
	l.Lock("l1", 1)
	fmt.Println("end lock 1")

	fmt.Println("start lock 2")
	l.Lock("l2", 2)
	fmt.Println("end lock 2")

}
