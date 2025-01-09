package listen

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVoice(t *testing.T) {
	recorder := NewAudioRecorder()
	if recorder == nil {
		t.Fatal("无法创建录音器")
	}

	err := recorder.StartRecording()
	assert.NoError(t, err)
	fmt.Println("开始录音...请说话测试...")

	time.Sleep(5 * time.Second)

	hash, err := recorder.StopRecording()
	assert.NoError(t, err)
	fmt.Printf("录音完成，文件哈希值: %s\n", hash)
}
