# 不需要将一个指向切片的指针传递给函数
切片实际上是一个指向潜在数组的指针,我们常常需要把切片作为一个参数传递给函数是因为;实际上是传递一个指向变量的指针,在函数内可以改变这个变量,而不是传递给数据的拷贝

因此应该这样做:
`func findBiggest( listOfNumber []int) int {}`

而不是:
`func findBiggest( listOfNumbers *[]int) int {}`

当切片作为参数传递时,切记不要解引用切片