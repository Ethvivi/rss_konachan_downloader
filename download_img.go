package main

import (
	"Regexp"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	dirname := "IMG"
	temp_dir := "./IMG"

	_,err := os.Stat(temp_dir)
	if err != nil{
		fmt.Println("dir maybe not exist")
		if os.IsNotExist(err) {
			fmt.Println("DIR IS NOT EXIT WILL CREATE")
			os.Mkdir(dirname, os.ModePerm)
		}
	} else {
		fmt.Println("dir has exist")
	}

	resp, err := http.Get("https://rsshub.ioiox.com/konachan.net/post/popular_recent/1d")
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	re := regexp.MustCompile(`"https://.*?"`)
	var link []string
	link = re.FindAllString(string(body),-1)
	//fmt.Println(link)
	var imgurl string
	for i:=0; i<len(link); i++{
		if i % 2 != 1 || i == 0 {
			continue
		} else {
			imgurl = link[i]
			fmt.Println(imgurl)
			resp, err := http.Get(imgurl[1:len(imgurl)-1])
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			var imglink string = imgurl[len(imgurl)-39:len(imgurl)-1]

			out, err := os.Create(dirname + "/" + imglink)
			if err != nil {
			panic(err)
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				panic(err)
			}
		}


	}
}