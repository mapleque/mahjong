package constant

// 返回数据中的status定义
// 本文件中所有常量定义必须以STATUS_开头，值必须大于9999
// STATUS_ERROR_* = 10000 - 19999 处理过程错误
// STATUS_INVALID_* = 20000 - 29999 参数错误、状态错误

const STATUS_ERROR_INVALID_USER = 10001 // 用户不合法

const STATUS_INVALID_MOBILE = 20001 // 手机号不合法
const STATUS_INVALID_TOKEN = 20002  // token非法
const STATUS_INVALID_OP = 20003     // 指定的字段不正确
const STATUS_INVALID_ID = 20004     // id不正确
