package main

import "fmt"

type Any interface{}

type EvalFunc func(Any) (Any, Any)

// 工厂函数BuildLazyEvaluator需要一个函数和一个初始状态作为输入参数,返回一个无参,返回值时生成序列的函数
// 传入的函数需要计算出下一个返回值以及下一个状态参数.
func BuildLazyEvaluator(evalFunc EvalFunc, initState Any) func() Any {
	retValChan := make(chan Any)
	loopFunc := func() {
		var actState Any = initState
		var retVal Any

		// 无限循环
		for {
			retVal, actState = evalFunc(actState)
			retValChan <- retVal
		}
	}
	retFunc := func() Any {
		return <-retValChan
	}

	// go协程
	go loopFunc()

	return retFunc
}

func BuildLazyIntEvaluator(evalFunc EvalFunc, initState Any) func() int {
	ef := BuildLazyEvaluator(evalFunc, initState)
	return func() int {
		return ef().(int)
	}
}

func main() {
	evenFunc := func(state Any) (Any, Any) {
		os := state.(int)

		ns := os + 2
		return os, ns
	}

	even := BuildLazyEvaluator(evenFunc, 0)

	for i := 0; i < 10; i++ {
		fmt.Printf("%vth even:%v\n", i, even())
	}
}
