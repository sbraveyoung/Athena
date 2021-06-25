**ast** is a library that parse a logic string like golang code, and return a boolean to report true or false.

### 支持的判断规则如下：

1. 条件字符串必须返回 boolean 类型，否则将报错。语法与 golang 条件判断语句一致；
2. 分严格模式和兼容模式。严格模式下语法与 golang 完全一致。兼容模式下，= 视同于 ==，& 视同于 &&，| 视同于 ||，字符串可不用带双引号（可能会和关键字冲突，请谨慎使用）；
3. 支持的运算符：
   1. ==：判断是否相等，支持 int、string 类型，两个操作数须为同一类型。string 类型带不带双引号均可。eg: app == media_std，app == "media_std"，"media_std" == app，1 == 2
   2. !=：判断不等，规则同上
   3. %：取模操作，支持 int 类型。eg: uid % 10
   4. \>：大于，支持 int 类型。eg: uid % 10 < 5
   5. <：小于，规则同上
   6. \>=： 大于等于，规则同上
   7. <=：小于等于，规则同上
   8. &&：逻辑与，当且仅当左右操作数分别为 boolean 类型时有效，左右操作数可以分别为 true/false 或 子逻辑表达式。eg: app == "media_std" && uid % 10 < 5，false && (app == "media_std" && uid % 10 >= 5)
   9. ||：逻辑或，规则同上
   10. !：逻辑非，一元运算符，操作数须为 boolean 类型。eg: !(app == "media_std")
   11. map[index ]：判断某元素是否存在某集合内。判断规则由以 map 为函数名，以 index 为入参, 以 bool 类型值为出参的函数给出。eg: cls_whitelist[uid]，cls_blacklist[uid] == true，gpu_whitelist[gpu] == false
   12. str[2:3]：字符串截取，只支持 string 类型。eg: cv[2:3] == "7"，cv[8:] == "Iphone"
4. 预定义函数：
   1. contains(s, substr)：判断字符串 s 中是否包含 substr，用法及出入参数与 golang 标准库中 strings.Contains 函数一致。eg: contains("IK7.8.9_Iphone", "Iphone")，contains(cv, "Android")
   2. mod(str, 10)：对以字符串类型标识的数字取模，第二个参数须为 int 类型，返回值为 int。eg: mod(uid, 10) < 5
