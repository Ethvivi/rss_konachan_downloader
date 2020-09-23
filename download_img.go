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
	dl := make(chan bool) //创建了信道 dl
	dirname := "IMG"
	temp_dir := "./IMG"
	dir_check(dirname,temp_dir)

	fmt.Println("___________________________________________download will go___________________________________________")

	// 将main阻塞 直到其他协程结束完
	go mutual_download("1d",dirname,dl)
	dl <- true
	go mutual_download("1w",dirname,dl)
	dl <- true
	go mutual_download("1m",dirname,dl)
	dl <- true

	fmt.Println("___________________________________________Finished___________________________________________")
}

func dir_check(dirname,temp_dir string)  {

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
}

func mutual_download(day,dirname string, dl chan bool)  {

	resp, err := http.Get("https://rsshub.ioiox.com/konachan.net/post/popular_recent/"+day)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
	}

	re := regexp.MustCompile(`"https://.*?"`)
	var link []string
	link = re.FindAllString(string(body),-1) //n=-1代表不限定次数 正数为代表限定次数

	var imgurl string

	//偶数位为图片的pixiv链接 跳过直接下载konachan原图链接
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
	<- dl // 读取信道的数据
}
