# ngoinx

记不住 `nginx.conf` 里面的诸多参量是极为正常的，大体来说也就三个原因：

* 我花的时间少了；
* 我日常生活很少用；
* 我等级还不够。

现在自我感觉良好，等级应该是够用，不会存在之前那么多无法理解的概念了。

自己有实现一些功能这个想法时，总是觉得一切都会水到渠成。但过程中难免会出各种设计问题。在感慨自己水平太低后，也会耐下心来参考这些优秀的软件是使用怎样的设计思路。

在双手离开键盘的那段时间里，我的大脑不在局限于编辑器框中的那几行代码，很多奇思妙想转瞬即逝，真是十分可惜。

通过写点东西来辅助记忆，这或许是极为低效的一种学习方式。

当然了，写这个还有其他目的，比如想打磨出一个或者两三个我觉得合格的 `interface`，以及我觉得比较恰当的包间关系，还有学习一些包的使用 balabala。

## 现在的进度

* [x] 负载均衡
  * 写是写好了，但是不耐用，还需要更多的学习。
  * [x] 一致性哈希
  * [x] 权重轮询法
  * [x] 日志输出
* [x] 反向代理
  * 不得不说 `golang` 的 `net/http` 真不戳，我自认为已经提供这个功能了。
* [ ] 动静分离
  * 这玩意儿咋整啊？我之前以为是静态文件缓存的意思，但现在发现和之前理解的有些偏差。
