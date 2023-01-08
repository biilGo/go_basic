# 并行化大量数据的计算
假设我们需要处理一些数量巨大且互不相干的数据项,它们从一个in通道被传递进来,当我们处理完以后又要将它们放入另一个out通道,就像工厂流水线一样.处理每个数据项也可能包含许多步骤:Preprocess预处理/StepA步骤A/StepB步骤B/.../PostProcess后处理

一个典型的用于解决按顺序执行每个步骤的顺序流水线算法可以写成下面这样:
```
func SerialProcessData(in <-chan *Data, out chan<- *Data) {
    for data := range in {
        tmpA := PreprocessData(data)
        tmpB := ProcessStepA(tmpA)
        tmpC := ProcessStepB(tmpB)
        out <- PostProcessData(tmpC)
    }
}
```

一次只执行一个步骤,并且按顺序处理每个项目:在第一个项目没有被`PostProcess`并放入out通道之前绝不会处理第二个项目

如果仔细想想,会发现这将造成巨大的时间浪费

一个更高效的计算方式是让每一个处理步骤作为一个协程独立工作,每一个步骤从上一步的输出通道中获得输入数据,这种方式仅有极少数时间浪费,而大部分时间所有的步骤都在一直执行中:
```
func ParallelProcessData (in <-chan *Data, out chan<- *Data) {
    // make channels:
    preOut := make(chan *Data, 100)
    stepAOut := make(chan *Data, 100)
    stepBOut := make(chan *Data, 100)
    stepCOut := make(chan *Data, 100)
    // start parallel computations:
    go PreprocessData(in, preOut)
    go ProcessStepA(preOut,StepAOut)
    go ProcessStepB(StepAOut,StepBOut)
    go ProcessStepC(StepBOut,StepCOut)
    go PostProcessData(StepCOut,out)
}
```

通道的缓冲区大小可以用来进一步优化整个过程