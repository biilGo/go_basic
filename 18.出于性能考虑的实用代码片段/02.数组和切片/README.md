# 数组和切片

创建：
```
arr1 := new([len]type)

slice1 := make([]type, len)
```

初始化:
```
arr1 := [...]type{i1, i2, i3, i4, i5}

arrKeyValue := [len]type{i1: val1, i2: val2}

var slice1 []type = arr1[start:end]
```

1. 如果截断数组或者也切片的最后一个元素:`line = line[:len(line)-1]`

2. 如果使用for或者for-range遍历一个数组或者切片
```
for i:=0; i < len(arr); i++ {
… = arr[i]
}
for ix, value := range arr {
…
}
```

3. 如何在一个二维数组或切片`arr2Dim`中查找一个指定值V:\
```
found := false
Found: for row := range arr2Dim {
    for column := range arr2Dim[row] {
        if arr2Dim[row][column] == V{
            found = true
            break Found
        }
    }
}
```