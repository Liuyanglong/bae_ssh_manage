bae_ssh_manage
==============

准备用的ssh登陆功能

todo list:

1、选择最优的port时，从对应的proxyserver直接读取passwd中的数量，选择最少；

2、选择最优时，注意全局变量60000的上限；

3、当某个proxyserver挂了时，增加一键恢复脚本；

4、当publicport下有多个proxyserver时，增加一键同步操作；

5、健康检查，对proxyserver list进行check，及时发现可用&不可用机器；
