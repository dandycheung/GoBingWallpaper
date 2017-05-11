package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"
)

const (
	api  = "https://cn.bing.com/HPImageArchive.aspx?format=js&n=1&idx="
	host = "https://cn.bing.com"
)

type Bing struct {
	Images []struct {
		Startdate     string        `json:"startdate"`
		Fullstartdate string        `json:"fullstartdate"`
		Enddate       string        `json:"enddate"`
		URL           string        `json:"url"`
		Urlbase       string        `json:"urlbase"`
		Copyright     string        `json:"copyright"`
		Copyrightlink string        `json:"copyrightlink"`
		Quiz          string        `json:"quiz"`
		Wp            bool          `json:"wp"`
		Hsh           string        `json:"hsh"`
		Drk           int           `json:"drk"`
		Top           int           `json:"top"`
		Bot           int           `json:"bot"`
		Hs            []interface{} `json:"hs"`
	} `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}

var bing Bing
var dir, imgUrl, imgFileName, idx string

func init() {

	userinfo, err := user.Current()
	if err != nil {
		logerr(fmt.Sprint(err))
	}
	defaultDir := userinfo.HomeDir + "/Pictures/BingWallpaper/"

	flag.StringVar(&dir, "dir", defaultDir, "指定图片保存的目录，默认为用户 \"Pictures\" 目录下的 \"BingWallpaper\" 目录，不存在会新建。")
	flag.StringVar(&idx, "idx", "0", "获取几天前的图片，默认为0表示当天，1表示1天前。")

	flag.Parse()
}

func main() {

	checkDir()
	getImgURL()
	getImg()

}

func loginfo(logstr string) {
	//记录日志
	log.Println("info:", logstr)

}
func logerr(logstr string) {
	//记录错误后退出程序。
	log.Fatalln("error:", logstr)

}

func getImgURL() {
	//获取图片地址

	loginfo("从接口获取图片地址……")

	resp, err := http.Get(api + idx)
	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片地址获取失败！")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片地址获取失败！")
	}

	err = json.Unmarshal(body, &bing)

	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片地址解析失败！")
	}

	imgUrl = fmt.Sprintf("%v%v", host, bing.Images[0].URL)
	loginfo(fmt.Sprint("图片地址为：", imgUrl))

	a := strings.Split(bing.Images[0].URL, "/")

	imgFileName = time.Now().Format("20060102-150405_") + a[len(a)-1]
	loginfo(fmt.Sprint("将要保存的图片名称为：", imgFileName))

}

func getImg() {
	//获取图片
	loginfo("开始下载图片……")

	resp, err := http.Get(imgUrl)
	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片获取失败！")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片获取失败！")
	}

	f, err := os.Create(dir + imgFileName)
	defer f.Close()

	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片保存失败！")
	}

	size, err := f.Write(body)

	if err != nil {
		loginfo(fmt.Sprint(err))
		logerr("图片保存失败！")
	}
	loginfo(fmt.Sprintf("图片保存成功：%v bytes。", size))

}

func checkDir() {
	//	检查图像目录是否存在，若不存在则建立。
	//	若出错则程序终止执行，执行完成则表明无问题。

	loginfo("检查图片目录是否正常……")
	_, err := os.Stat(dir)

	if err != nil {
		//		查看目录出错
		loginfo(fmt.Sprint(err))
		if strings.Contains(fmt.Sprint(err), "no such file or directory") {
			loginfo("没有发现图片目录，将新建目录：" + dir)
			err = os.Mkdir(dir, 0755)
			if err != nil {
				//				创建目录报错
				if strings.Contains(fmt.Sprint(err), "no such file or directory") {
					loginfo(fmt.Sprint(err))
					logerr(fmt.Sprintf("%s 目录创建失败，请检查父目录是否存在。", dir))
				} else if strings.Contains(fmt.Sprint(err), "permission denied") {
					loginfo(fmt.Sprint(err))
					logerr(fmt.Sprintf("%s 目录创建失败，请检查父目录权限。", dir))
				} else {
					loginfo(fmt.Sprint(err))
					logerr(fmt.Sprintf("%s 目录创建失败，请手动创建。", dir))
				}
			} else {
				//				检查目录是否创建成功
				loginfo("目录创建完成，检查是否创建成功。")
				_, err := os.Stat(dir)
				if err != nil {
					loginfo(fmt.Sprint(err))
					logerr("目录创建失败，请手动创建。")
				} else {
					loginfo("目录创建未出现错误，似乎创建成功了。")
				}
			}
		} else if strings.Contains(fmt.Sprint(err), "permission denied") {
			logerr("目录检测失败，请检查目录权限。")
		} else {
			logerr("目录检测失败，请检查原因。")
		}
	} else {
		loginfo("目录似乎正常。")
	}

	//创建文件测试目录是否可写
	loginfo("程序将创建测试文件以检查目录是否可写……")
	f, err := os.Create(dir + "test")
	f.WriteString("test file.")
	defer f.Close()

	if err != nil {
		//				创建文件报错
		if strings.Contains(fmt.Sprint(err), "no such file or directory") {
			loginfo(fmt.Sprint(err))
			logerr(fmt.Sprintf("%s 测试文件创建失败，请检查父目录是否存在。", dir+"test"))
		} else if strings.Contains(fmt.Sprint(err), "permission denied") {
			loginfo(fmt.Sprint(err))
			logerr(fmt.Sprintf("%s 测试文件创建失败，请检查父目录权限。", dir+"test"))
		} else {
			loginfo(fmt.Sprint(err))
			logerr(fmt.Sprintf("%s 测试文件创建失败，请手动检查失败原因。", dir+"test"))
		}
	} else {
		err = os.Remove(dir + "test")
		if err != nil {
			loginfo(fmt.Sprint(err))
			logerr("测试文件移除失败。")
		} else {
			loginfo("目录可写。")
		}

	}

}
