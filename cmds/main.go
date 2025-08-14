package main

import (
	//"bufio"
	"log"

	"github.com/xiaowumin-mark/go-trans/bing"
	//"os"
	//"github.com/xiaowumin-mark/go-trans/youdao"
)

func main() {
	/*reader := bufio.NewReader(os.Stdin)
	log.Print("请输入需要翻译的内容：")
	query, _ := reader.ReadString('\n')

	y := youdao.New(query, "en")
	d, err := y.SimpleTranslate()
	if err != nil {
		panic(err)
	}
	log.Println(d)*/

	res, err := bing.BatchTranslate([]string{"hello", "Hello world"}, "en", "zh")
	if err != nil {
		panic(err)
	}
	for _, v := range res.Parsed {
		log.Println(v.Translations[0].Text)
	}

}

/*
hello 2
hello world 8

mismatched type with value 3
y 8

yes; 1*/
