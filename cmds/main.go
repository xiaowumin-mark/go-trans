package main

import (
	"bufio"
	"log"
	"os"

	"github.com/xiaowumin-mark/go-trans/youdao"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	log.Print("请输入需要翻译的内容：")
	query, _ := reader.ReadString('\n')

	y := youdao.New(query, "en")
	d, err := y.Translate()
	if err != nil {
		panic(err)
	}

	if d.Parsed.Meta.IsHasSimpleDict == "1" { // 为单词或词组
		log.Println("单词或词组")

		log.Println("结果: ", d.Parsed.WebTrans.WebTranslation[0].Trans[0].Value)

	} else {
		log.Println("句子")
		log.Println("结果: ", d.Parsed.Fanyi.Tran)
	}

}

/*
hello 2
hello world 8

mismatched type with value 3
y 8

yes; 1*/
