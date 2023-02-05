# 其他
如何在程序出错时终止程序:
```
if err != nil {
   fmt.Printf("Program stopping with error %v", err)
   os.Exit(1)
}
```

或者:
```
if err != nil { 
    panic("ERROR occurred: " + err.Error())
}
```