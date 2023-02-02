# Go中的密码学
通过网络传输的数据必须加密,以防止被hacker读取或篡改,并且保证发出的数据和接收的数据检验和一致鉴于go母公司的业务,go的标准库为该领域提供了超过30个的包

- `hash`包:实现了 `adler32`、`crc32`、`crc64` 和 `fnv` 校验；
- `cryto`包:实现了其它的 hash 算法，比如 `md4`、`md5`、`sha1` 等。以及完整地实现了`aes`、`blowfish`、`rc4`、`rsa`、`xtea` 等加密算法

通过调用 sha1.New() 创建了一个新的 hash.Hash 对象，用来计算 SHA1 校验值。Hash 类型实际上是一个接口，它实现了 io.Writer 接口：

```
type Hash interface {
    // Write (via the embedded io.Writer interface) adds more data to the running hash.
    // It never returns an error.
    io.Writer
    // Sum appends the current hash to b and returns the resulting slice.
    // It does not change the underlying hash state.
    Sum(b []byte) []byte
    // Reset resets the Hash to its initial state.
    Reset()
    // Size returns the number of bytes Sum will return.
    Size() int
    // BlockSize returns the hash's underlying block size.
    // The Write method must be able to accept any amount
    // of data, but it may operate more efficiently if all writes
    // are a multiple of the block size.
    BlockSize() int
}
```

通过`io.WriteString`或`hasher.Write`将给定的`[]byte`附件到dangqian的`hash.Hash`对象中