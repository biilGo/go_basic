# 短网址项目简介
你肯定知道浏览器中的地址非常长且/或复杂,在网上又一些将他们转换成简短URL来使用的服务,我们的项目与此类似:它具有2个功能的web服务(webservice)

## 添加(Add)
给定一个较长的URL,会将其转换成较短的版本,比如:
`http://maps.google.com/maps?f=q&source=s_q&hl=en&geocode=&q=tokyo&sll=37.0625,-95.677068&sspn=68.684234,65.566406&ie=UTF8&hq=&hnear=Tokyo,+Japan&t=h&z=9`

- (A)转变为:`http://goto/UrcGq`
- (B)并保存这对数据

## 重定向(Redirect)
短网址被请求时,会把用户重定向到原始的长URL.因此如果你在浏览器输入网址(B),会被重定向到页面(A)