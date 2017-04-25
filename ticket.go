/*
	Author        : tuxpy
	Email         : q8886888@qq.com.com
	Create time   : 2017-04-24 23:19:50
	Filename      : ticket.go
	Description   :
*/

package locker

type Ticketer interface {
	Return()
	Take()
	Size() int
	Used() int
}

type Ticket uint8

type _MyTicketer struct {
	tickets chan Ticket
	length  int
}

func (ticketer *_MyTicketer) Take() {
	<-ticketer.tickets
}

func (ticketer *_MyTicketer) Used() int {
	return ticketer.Size() - len(ticketer.tickets)
}

func (ticketer *_MyTicketer) Size() int {
	return cap(ticketer.tickets)
}

func (ticketer *_MyTicketer) Return() {
	ticketer.tickets <- Ticket(0)
}

func NewTicketer(concurrent int) Ticketer {
	tickets := make(chan Ticket, concurrent)
	for i := 0; i < concurrent; i++ {
		tickets <- Ticket(0)
	}

	return &_MyTicketer{
		length:  concurrent,
		tickets: tickets,
	}
}
