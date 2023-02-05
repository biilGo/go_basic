# 字符串

1. 如何修改字符串中的一个字符
```
str:="hello"
c:=[]byte(str)
c[0]='c'
s2:= string(c) // s2 == "cello"
```

2. 如何获取字符串的子串:`substr := str[n:m]`

3. 如何使用`for`或者`for-range`遍历一个字符串
```
// gives only the bytes:
for i:=0; i < len(str); i++ {
… = str[i]
}
// gives the Unicode characters:
for ix, ch := range str {
…
}
```

4. 如何获取一个字符串的字节数:`len(str)`

如果获取一个字符串的字符数:

最快速:`utf8.RuneCountInString(str)`,`len([]rune(str))`

5. 如何连接字符串:

最快速:`with a bytes.Buffer`,`Strings.Join()`

使用`+=`
```
 str1 := "Hello " 
 str2 := "World!"
 str1 += str2 //str1 == "Hello World!"
```