package main

import (
	"fmt"
	"github.com/schollz/progressbar"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func main() {
	dl := make(chan bool) //创建了信道 dl
	dirname := "IMG"
	temp_dir := "./IMG"
	dir_check(dirname,temp_dir)

	fmt.Println("___________________________________________download will go___________________________________________")

	var wg sync.WaitGroup
	wg.Add(3)//创建3个协程任务池
	// 将主协程阻塞 直到其他协程结束完
	go mutual_download("1d",dirname,dl,&wg)
	go mutual_download("1w",dirname,dl,&wg)
	go mutual_download("1m",dirname,dl,&wg)
	wg.Wait() //等待所有子协程释放

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

func mutual_download(day,dirname string, dl chan bool , wg *sync.WaitGroup)  {

	resp, err := http.Get("https://rsshub.ioiox.com/konachan.net/post/popular_recent/"+day)
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

			bar := progressbar.DefaultBytes(
				resp.ContentLength,
				"downloading",
			)

			_, err = io.Copy(io.MultiWriter(out,bar), resp.Body)
			if err != nil {
				panic(err)
			}
		}
	}
	defer wg.Done() //协程任务递减
	//<-dl // 读取信道的数据
}
