# gee framework
## introduction
diy golang web framework

## 参考链接
[前缀树算法实现路由匹配原理解析：Go 实现](https://jishuin.proginn.com/p/763bfbd2a2ba)

[七天用Go从零实现系列](https://geektutu.com/post/gee-day4.html)

[Go 数据结构和算法篇（十三）：字符串匹配之 Trie 树](https://geekr.dev/posts/go-trie-tree-algorithm)

## Trie 树的应用
Trie 树适用于那些查找前缀匹配的字符串，比如敏感词过滤和搜索框联想功能。

### 敏感词过滤系统
首先运营在后台手动更新敏感词，底层通过 Tire 树构建敏感词库，然后当商家发布商品时，以商品标题+详情作为主串，
将敏感词库作为模式串，进行匹配，如果模式串和主串有匹配字符，则以此为起点，继续往后匹配，直到匹配出完整字符串，
然后标记为匹配出该敏感词（如果想嗅探所有敏感词，继续往后匹配），
否则将主串匹配起点位置往后移，从下一个字符开始，继续与模式串匹配。

### 搜索框联想功能
另外，搜索框的查询关键词联想功能也是基于 Trie 树实现的：

![google](https://coding-pic.oss-cn-hangzhou.aliyuncs.com/images/2022-05-07-img.png)

Google搜索框联想词, 进而可以扩展到浏览器网址输入自动补全、IDE 代码编辑器自动补全、输入法自动补全功能等。

## 优化点
在ServeHTTP()中，每接收到一个请求都需要逐一匹配前缀来添加中间件，太影响性能了

获取某个请求要执行的中间件，不应该采用循环的方式，可以采用Trie树,或者在请求之前就通过一个map[string][]HandlerFunc来保存
,但在这里需要注意如请求/a/b和/a/b/,所以map的string要存两份

gin是在前缀树的节点中添加中间件的切片，这样在匹配动态路由并解析参数时，就可以同时获得各分组的中间件。




```go
var middlewares []HandlerFunc
for _, routerGroup := range engine.routerGroups {
    if strings.HasPrefix(req.URL.Path, routerGroup.prefix) {
        middlewares = append(middlewares, routerGroup.middlewares...)
    }
}
c.handlers = middlewares
```

## feat
- recover panic
- ![recover panic](https://coding-pic.oss-cn-hangzhou.aliyuncs.com/images/2022-05-07-L0Ov4C.png)

### 
![](https://dong4j-imgs.oss-cn-hangzhou.aliyuncs.com/oDE4Cl.png)
![](https://coding-pic.oss-cn-hangzhou.aliyuncs.com/images/2022-05-07-L0Ov4C.png)