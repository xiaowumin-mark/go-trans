# Youdao 翻译库使用说明文档

## 概述

这是一个基于有道翻译 API 的 Go 语言库，用于实现中英文单词、词组和句子的翻译功能。该库通过模拟 Web 请求获取有道词典的 JSON 格式翻译结果，并提供结构化数据访问。

## 快速开始

### 安装

```bash
go get github.com/xiaowumin-mark/go-trans/youdao
```

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

## 支持的语言代码

 - 中英 :`en`
 - 中法 :`fr`
 - 中韩 :`ko`
 - 中日 :`ja`

# Bing 翻译库使用说明文档

## 概述

Bing 翻译库是基于微软 Edge 浏览器翻译功能的 API 实现的 Go 语言库。该库提供了批量翻译功能，适用于大量文本的翻译任务。通过调用 Bing 的批量翻译接口，可以实现多种语言之间的文本翻译。

## 快速开始

### 安装

```bash
go get github.com/xiaowumin-mark/go-trans/bing
```

### 基本用法

```go
package main

import (
    "log"
    "github.com/xiaowumin-mark/go-trans/bing"
)

func main() {
    res, err := bing.BatchTranslate([]string{"hello", "Hello world"}, "en", "zh")
    if err != nil {
        panic(err)
    }
    for _, v := range res.Parsed {
        log.Println(v.Translations[0].Text)
    }
}
```

## 功能详解

### BatchTranslate 函数

这是 Bing 翻译库的核心函数，用于执行批量翻译操作。

```go
func BatchTranslate(text []string, from string, to string) (*BatchTranslateResp, error)
```

**参数说明：**
- `text []string`: 需要翻译的文本切片，支持批量翻译
- `from string`: 源语言代码（如："en"、"zh-Hans"、"zh-Hant"、"zh-cn"）
- `to string`: 目标语言代码（如："zh"、"en"、"ja"）

**返回值：**
- `*BatchTranslateResp`: 包含原始响应和解析后结果的结构体
- `error`: 错误信息（如果有）

## 使用示例

### 基本翻译示例

```go
res, err := bing.BatchTranslate([]string{"hello", "Hello world"}, "en", "zh")
if err != nil {
    log.Fatal(err)
}

// 遍历翻译结果
for _, v := range res.Parsed {
    log.Println(v.Translations[0].Text)
}
```

## 支持的语言代码

- `zh`: 中文
- `zh-Hans`: 简体中文
- `zh-Hant`: 繁体中文
- `en`: 英语
- `ja`: 日语
- `ko`: 韩语
- `fr`: 法语
- `de`: 德语
- `es`: 西班牙语

## 依赖库

- `github.com/bytedance/sonic`: 高性能 JSON 解析库
