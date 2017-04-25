/*
	Author        : tuxpy
	Email         : q8886888@qq.com.com
	Create time   : 2017-04-24 22:15:08
	Filename      : main.go
	Description   :
*/

package locker

import (
	"errors"
	"fmt"
)

type LockerCommand uint8

const (
	LOCK   LockerCommand = 1
	UNLOCK LockerCommand = 2
	CLOSE  LockerCommand = 0
)

type Locker interface {
	Lock(name string, concurrent int) error
	Unlock(name string) error
	Close()
}

type _MyLocker struct {
	ticketers         map[string]Ticketer
	running           bool
	cmd_ch            chan Command
	ticketer_max_size int
	closed            bool
}

type Command struct {
	lock_name  string
	concurrent int
	cmd        LockerCommand
	result_ch  chan bool
}

func (locker *_MyLocker) Lock(name string, concurrent int) error {
	if err := locker.CheckStatus(); err != nil {
		return err
	}

	result_ch := make(chan bool)
	locker.cmd_ch <- Command{
		cmd:        LOCK,
		concurrent: concurrent,
		lock_name:  name,
		result_ch:  result_ch,
	}

	<-result_ch
	return nil
}

func (locker *_MyLocker) Unlock(name string) error {
	if err := locker.CheckStatus(); err != nil {
		return err
	}

	result_ch := make(chan bool)
	locker.cmd_ch <- Command{
		cmd:       UNLOCK,
		lock_name: name,
		result_ch: result_ch,
	}
	<-result_ch

	return nil
}

func (locker *_MyLocker) Close() {
	locker.cmd_ch <- Command{
		cmd: CLOSE,
	}
}

/*
	执行锁定操作
	@param command 命令结构体
*/
func (locker *_MyLocker) _ExecLock(command Command) {
	ticketer, ok := locker.ticketers[command.lock_name]
	if !ok {
		ticketer = NewTicketer(command.concurrent)
		locker.ticketers[command.lock_name] = ticketer
	}

	// 防止不同的lock_name阻塞当前goroutine
	go func() {
		ticketer.Take()
		command.result_ch <- true
	}()
}

/*
	执行解锁操作
	@param command 命令结构体
*/
func (locker *_MyLocker) _ExecUnlock(command Command) {
	ticketer, ok := locker.ticketers[command.lock_name]
	if !ok {
		return
	}

	// 防止不同的lock_name阻塞当前goroutine
	go func() {
		ticketer.Return()
		command.result_ch <- true
	}()
}

/*
	如果当前的ticketer超过了ticketer_max_size则需要将空的ticketer删除

*/
func (locker *_MyLocker) Flush() {
	if len(locker.ticketers) < locker.ticketer_max_size {
		return
	}

	for locker_name, ticketer := range locker.ticketers {
		fmt.Println(locker_name, ticketer.Used())
		if ticketer.Used() == 0 {
			delete(locker.ticketers, locker_name)
		}
	}
}

/*
	检查locker的状态

	@return error
*/
func (locker *_MyLocker) CheckStatus() error {
	if locker.closed {
		return errors.New("locker already closed")
	}

	return nil
}

/*
	单独在一个goroutine中通过chan 接受指令
*/

func (locker *_MyLocker) Run() {
	if locker.running {
		return
	}

	locker.running = true

	for {
		command, opened := <-locker.cmd_ch
		if opened == false {
			break
		}

		switch command.cmd {
		case CLOSE:
			if !locker.closed {
				close(locker.cmd_ch)
				locker.closed = true
			}

		case LOCK:
			locker.Flush()
			locker._ExecLock(command)

		case UNLOCK:
			locker._ExecUnlock(command)
		}
	}

}

func New() Locker {
	ticketer_max_size := 20
	my_locker := &_MyLocker{
		ticketers:         make(map[string]Ticketer, ticketer_max_size),
		cmd_ch:            make(chan Command),
		ticketer_max_size: ticketer_max_size,
		closed:            false,
	}

	go my_locker.Run()

	return my_locker
}
