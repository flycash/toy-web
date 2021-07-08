package main

import (
	"fmt"
	"io"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "Hi, this is home page")
}

func readBodyOnce(w http.ResponseWriter, r *http.Request)  {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		// 记住要返回，不然就还会执行后面的代码
		return
	}
	// 类型转换，将 []byte 转换为 string
	fmt.Fprintf(w, "read the data: %s \n", string(body))

	// 尝试再次读取，啥也读不到，但是也不会报错
	body, err = io.ReadAll(r.Body)
	if err != nil {
		// 不会进来这里
		fmt.Fprintf(w, "read the data one more time got error: %v", err)
	}
	fmt.Fprintf(w, "read the data one more time: %s and read nothing\n", string(body))
}

func readBodyMultipleTimes(w http.ResponseWriter, r *http.Request) {
	reader, err := r.GetBody()
	if err != nil {
		fmt.Fprintf(w, "GetBody error: %v", err)
		return
	}
	
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(w, "read data from reader error: %v", err)
		return
	}

	// 类型转换，将 []byte 转换为 string
	fmt.Fprintf(w, "read the data: %s", string(body))

}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/body/once", readBodyOnce)
	http.HandleFunc("/body/multi", readBodyMultipleTimes)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}