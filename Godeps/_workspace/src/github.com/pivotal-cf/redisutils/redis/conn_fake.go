// This file was generated by counterfeiter
// counterfeiter -o redis/conn_fake.go --fake-name ConnFake ~/go/src/github.com/garyburd/redigo/redis/conn.go Conn

package redis

import (
	"sync"

	"github.com/garyburd/redigo/redis"
)

//ConnFake ...
type ConnFake struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	ErrStub        func() error
	errMutex       sync.RWMutex
	errArgsForCall []struct{}
	errReturns     struct {
		result1 error
	}
	DoStub        func(commandName string, args ...interface{}) (reply interface{}, err error)
	doMutex       sync.RWMutex
	doArgsForCall []struct {
		commandName string
		args        []interface{}
	}
	doReturns struct {
		result1 interface{}
		result2 error
	}
	SendStub        func(commandName string, args ...interface{}) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		commandName string
		args        []interface{}
	}
	sendReturns struct {
		result1 error
	}
	FlushStub        func() error
	flushMutex       sync.RWMutex
	flushArgsForCall []struct{}
	flushReturns     struct {
		result1 error
	}
	ReceiveStub        func() (reply interface{}, err error)
	receiveMutex       sync.RWMutex
	receiveArgsForCall []struct{}
	receiveReturns     struct {
		result1 interface{}
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

//NewConnFake is the preferred way to initialise a ConnFake
func NewConnFake() *ConnFake {
	return new(ConnFake)
}

//Close ...
func (fake *ConnFake) Close() error {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	return fake.closeReturns.result1
}

//CloseCallCount ...
func (fake *ConnFake) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

//CloseReturns ...
func (fake *ConnFake) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

//Err ...
func (fake *ConnFake) Err() error {
	fake.errMutex.Lock()
	fake.errArgsForCall = append(fake.errArgsForCall, struct{}{})
	fake.recordInvocation("Err", []interface{}{})
	fake.errMutex.Unlock()
	if fake.ErrStub != nil {
		return fake.ErrStub()
	}
	return fake.errReturns.result1
}

//ErrCallCount ...
func (fake *ConnFake) ErrCallCount() int {
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	return len(fake.errArgsForCall)
}

//ErrReturns ...
func (fake *ConnFake) ErrReturns(result1 error) {
	fake.ErrStub = nil
	fake.errReturns = struct {
		result1 error
	}{result1}
}

//Do ...
func (fake *ConnFake) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	fake.doMutex.Lock()
	fake.doArgsForCall = append(fake.doArgsForCall, struct {
		commandName string
		args        []interface{}
	}{commandName, args})
	fake.recordInvocation("Do", []interface{}{commandName, args})
	fake.doMutex.Unlock()
	if fake.DoStub != nil {
		return fake.DoStub(commandName, args...)
	}
	return fake.doReturns.result1, fake.doReturns.result2
}

//DoCallCount ...
func (fake *ConnFake) DoCallCount() int {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return len(fake.doArgsForCall)
}

//DoArgsForCall ...
func (fake *ConnFake) DoArgsForCall(i int) (string, []interface{}) {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return fake.doArgsForCall[i].commandName, fake.doArgsForCall[i].args
}

//DoReturns ...
func (fake *ConnFake) DoReturns(result1 interface{}, result2 error) {
	fake.DoStub = nil
	fake.doReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

//Send ...
func (fake *ConnFake) Send(commandName string, args ...interface{}) error {
	fake.sendMutex.Lock()
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		commandName string
		args        []interface{}
	}{commandName, args})
	fake.recordInvocation("Send", []interface{}{commandName, args})
	fake.sendMutex.Unlock()
	if fake.SendStub != nil {
		return fake.SendStub(commandName, args...)
	}
	return fake.sendReturns.result1
}

//SendCallCount ...
func (fake *ConnFake) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

//SendArgsForCall ...
func (fake *ConnFake) SendArgsForCall(i int) (string, []interface{}) {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return fake.sendArgsForCall[i].commandName, fake.sendArgsForCall[i].args
}

//SendReturns ...
func (fake *ConnFake) SendReturns(result1 error) {
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

//Flush ...
func (fake *ConnFake) Flush() error {
	fake.flushMutex.Lock()
	fake.flushArgsForCall = append(fake.flushArgsForCall, struct{}{})
	fake.recordInvocation("Flush", []interface{}{})
	fake.flushMutex.Unlock()
	if fake.FlushStub != nil {
		return fake.FlushStub()
	}
	return fake.flushReturns.result1
}

//FlushCallCount ...
func (fake *ConnFake) FlushCallCount() int {
	fake.flushMutex.RLock()
	defer fake.flushMutex.RUnlock()
	return len(fake.flushArgsForCall)
}

//FlushReturns ...
func (fake *ConnFake) FlushReturns(result1 error) {
	fake.FlushStub = nil
	fake.flushReturns = struct {
		result1 error
	}{result1}
}

//Receive ...
func (fake *ConnFake) Receive() (reply interface{}, err error) {
	fake.receiveMutex.Lock()
	fake.receiveArgsForCall = append(fake.receiveArgsForCall, struct{}{})
	fake.recordInvocation("Receive", []interface{}{})
	fake.receiveMutex.Unlock()
	if fake.ReceiveStub != nil {
		return fake.ReceiveStub()
	}
	return fake.receiveReturns.result1, fake.receiveReturns.result2
}

//ReceiveCallCount ...
func (fake *ConnFake) ReceiveCallCount() int {
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	return len(fake.receiveArgsForCall)
}

//ReceiveReturns ...
func (fake *ConnFake) ReceiveReturns(result1 interface{}, result2 error) {
	fake.ReceiveStub = nil
	fake.receiveReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

//Invocations ...
func (fake *ConnFake) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	fake.flushMutex.RLock()
	defer fake.flushMutex.RUnlock()
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	return fake.invocations
}

func (fake *ConnFake) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ redis.Conn = new(ConnFake)
