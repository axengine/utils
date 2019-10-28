# codec []byte编解码方法

使用基本的位运算对数据进行编解码，编码后的数据是同大小的[]byte。提供三种方式，可以自行选择和改造。

codec提供了一种简单的思路对数据进行混淆，可以使用在app与server时间数据混淆，前提是app代码不被逆向和重放。

- ENCODE_BIT_NOT：按位取反
- ENCODE_BYTE_RVS：逆序
- ENCODE_LOOP_XOR：环形异或