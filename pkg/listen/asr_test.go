package listen

import (
	"fmt"
	"os"
	"testing"
)

func TestAsr(t *testing.T) {
	appid := ""
	token := ""
	cluster := ""
	audioPath := ""      // 本地音频路径
	audioFormat := "wav" // wav 或 mp3，根据音频格式设置

	client := buildAsrClient()
	client.Appid = appid
	client.Token = token
	client.Cluster = cluster
	client.Format = audioFormat

	audioData, err := os.ReadFile(audioPath)
	if err != nil {
		fmt.Println("fail to read audio file", err.Error())
		return
	}
	asrResponse, err := client.requestAsr(audioData)
	if err != nil {
		fmt.Println("fail to request asr, ", err.Error())
		return
	}
	fmt.Println(asrResponse.Reqid, asrResponse.Code, asrResponse.Message, asrResponse.Sequence)
	fmt.Println(asrResponse.Results)
}
