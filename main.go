package main

import (
	"fmt"
	"time"
	"sync"
	"sync/atomic"
)

func main() {

	var balance int32 = 0
	count := 100000000
	t1 := time.Now()
	//transLock := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i:=0;i<count;i++ {
		wg.Add(1)
		//go transferLock(&balance,&transLock,&wg)//3.1518197s|32.0352204s
		//go transferCAS(&balance,&wg)//3.126052s|31.8275625s
		go transferFAA(&balance,&wg)//3.1998921s|31.6946224s
	}
	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println(" 耗时: ", elapsed)
	fmt.Println("balance:",balance)
}


//锁实现方式
func transferLock(balance *int32, lock *sync.Mutex,wg *sync.WaitGroup) {
	defer wg.Done()
	lock.Lock()
	*balance = *balance +1
	lock.Unlock()
}

//CAS实现方式 : 先读取数据，再计算，再更新
func transferCAS(balance *int32,wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		//atomic原子操作
		oldbalance := atomic.LoadInt32(balance)
		newbalance := oldbalance + 1
		if atomic.CompareAndSwapInt32(balance, oldbalance, newbalance) {
			break
		}
	}
}


//FAA实现方式：直接更新
func transferFAA(balance *int32,wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt32(balance,1)
}
