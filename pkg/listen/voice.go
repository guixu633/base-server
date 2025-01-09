package listen

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
)

type AudioRecorder struct {
	stream      *portaudio.Stream
	buffer      []float32
	sampleRate  float64
	isRecording bool
}

func NewAudioRecorder() *AudioRecorder {
	// 初始化 PortAudio
	err := portaudio.Initialize()
	if err != nil {
		fmt.Printf("无法初始化 PortAudio: %v\n", err)
		return nil
	}

	return &AudioRecorder{
		sampleRate:  44100,
		buffer:      make([]float32, 0),
		isRecording: false,
	}
}

// StartRecording 开始录音
func (a *AudioRecorder) StartRecording() error {
	if a.isRecording {
		return fmt.Errorf("已经在录音中")
	}

	inputChannels := 1
	framesPerBuffer := make([]float32, 1024)

	// 打开默认输入设备的音频流
	stream, err := portaudio.OpenDefaultStream(inputChannels, 0, a.sampleRate, len(framesPerBuffer), framesPerBuffer)
	if err != nil {
		return fmt.Errorf("无法打开音频流: %v", err)
	}

	err = stream.Start()
	if err != nil {
		stream.Close()
		return fmt.Errorf("无法启动音频流: %v", err)
	}

	a.stream = stream
	a.isRecording = true
	a.buffer = make([]float32, 0)

	// 在后台持续录音
	go func() {
		for a.isRecording {
			// 读取音频数据
			err := stream.Read()
			if err != nil {
				fmt.Printf("读取音频数据错误: %v\n", err)
				continue
			}

			// 将缓冲区数据追加到录音缓冲区
			a.buffer = append(a.buffer, framesPerBuffer...)

			// 可选：打印音频电平，用于调试
			maxLevel := float32(0)
			for _, sample := range framesPerBuffer {
				if abs(sample) > maxLevel {
					maxLevel = abs(sample)
				}
			}
			if maxLevel > 0.01 { // 只打印有声音的部分
				fmt.Printf("当前音频电平: %.2f\n", maxLevel)
			}
		}
	}()

	return nil
}

// 辅助函数：计算绝对值
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// StopRecording 停止录音并保存文件
func (a *AudioRecorder) StopRecording() (string, error) {
	if !a.isRecording {
		return "", fmt.Errorf("没有正在进行的录音")
	}

	a.isRecording = false
	a.stream.Stop()
	a.stream.Close()

	// 终止 PortAudio
	portaudio.Terminate()

	// 生成文件名（使用时间戳）
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("recording_%s.wav", timestamp)

	// 保存音频文件
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("无法创建音频文件: %v", err)
	}
	defer file.Close()

	// 写入 WAV 文件头
	numSamples := len(a.buffer)
	writeWAVHeader(file, numSamples, int(a.sampleRate))

	// 将浮点数据转换为字节并写入文件
	for _, sample := range a.buffer {
		// 将 float32 转换为 16 位整数
		intSample := int16(sample * 32767.0)
		err := binary.Write(file, binary.LittleEndian, intSample)
		if err != nil {
			return "", fmt.Errorf("写入音频数据失败: %v", err)
		}
	}

	// 计算文件的 SHA256 哈希值
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("无法重置文件指针: %v", err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("计算哈希值失败: %v", err)
	}

	hashString := hex.EncodeToString(hash.Sum(nil))

	return hashString, nil
}

// writeWAVHeader 写入 WAV 文件头
func writeWAVHeader(file *os.File, numSamples int, sampleRate int) error {
	// WAV 文件头
	header := struct {
		ChunkID       [4]byte // "RIFF"
		ChunkSize     uint32  // 文件大小 - 8
		Format        [4]byte // "WAVE"
		Subchunk1ID   [4]byte // "fmt "
		Subchunk1Size uint32  // 16 for PCM
		AudioFormat   uint16  // 1 for PCM
		NumChannels   uint16  // 1 for mono
		SampleRate    uint32  // 采样率
		ByteRate      uint32  // SampleRate * NumChannels * BitsPerSample/8
		BlockAlign    uint16  // NumChannels * BitsPerSample/8
		BitsPerSample uint16  // 16 bits
		Subchunk2ID   [4]byte // "data"
		Subchunk2Size uint32  // 数据大小
	}{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   1,
		SampleRate:    uint32(sampleRate),
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
	}

	// 计算文件大小
	header.ByteRate = header.SampleRate * uint32(header.NumChannels) * uint32(header.BitsPerSample) / 8
	header.BlockAlign = header.NumChannels * header.BitsPerSample / 8
	header.Subchunk2Size = uint32(numSamples) * uint32(header.NumChannels) * uint32(header.BitsPerSample) / 8
	header.ChunkSize = 36 + header.Subchunk2Size

	return binary.Write(file, binary.LittleEndian, header)
}
