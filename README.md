# Youdao 翻译库使用说明文档

## 概述

这是一个基于有道翻译 API 的 Go 语言库，用于实现中英文单词、词组和句子的翻译功能。该库通过模拟 Web 请求获取有道词典的 JSON 格式翻译结果，并提供结构化数据访问。

## 安装

```bash
go get github.com/xiaowumin-mark/go-trans/youdao
```

## 快速开始

请参考 [cmds/main.go](https://github.com/xiaowumin-mark/go-trans/blob/main/cmds/main.go)

## 使用方法

### 1. 创建翻译器实例

```go
y := youdao.New("hello", "en")
```

### 2. 执行翻译

```go
result, err := y.Translate()
if err != nil {
    // 处理错误
}
```

### 3. 处理结果

根据 `Meta.IsHasSimpleDict` 判断翻译内容类型：

- `"1"`: 单词或词组，使用 `EC.WebTrans` 获取结果
- `"0"`: 句子，使用 `Fanyi.Tran` 获取结果

## 支持的翻译类型

| 类型 | 判断条件 | 结果字段 |
|------|----------|----------|
| 单词/词组 | `Meta.IsHasSimpleDict == "1"` | `EC.WebTrans[0]` |
| 句子 | `Meta.IsHasSimpleDict == "0"` | `Fanyi.Tran` |

## 依赖库

- `github.com/bytedance/sonic`: 高性能 JSON 解析库
