package aivoice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/eric2788/go-silk/multiplat"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	silk "github.com/wdvxdr1123/go-silk"
)

const (
	AudioSamplingRateMP3  = "22050"
	AudioBitRate          = "12.2k" // in Hz
	NumberOfAudioChannels = "1"
	AudioSamplingRateAMR  = "8000"
)

func WavToSilk(b []byte) (data []byte, err error) {
	md := md5.Sum(b)
	tempName := hex.EncodeToString(md[:])

	wav, pcm := tempName+".wav", tempName+".pcm"

	err = os.WriteFile(wav, b, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer os.Remove(wav)

	// 2.转换pcm
	cmd := exec.Command("ffmpeg", "-i", wav, "-f", "s16le", "-ar", "24000", "-ac", "1", pcm)
	multiplat.HideWindow(cmd)
	if err = cmd.Run(); err != nil {
		return nil, err
	}
	defer os.Remove(pcm)
	pcmByte, err := os.ReadFile(pcm)
	if err != nil {
		return nil, err
	}
	return silk.EncodePcmBuffToSilk(pcmByte, 24000, 24000, true)
}

// WavToAmr Wav To Amr file
// Deprecated: use WavToSilk instead
func WavToAmr(b []byte) (data []byte, err error) {
	hash := md5.Sum(b)
	name := hex.EncodeToString(hash[:])
	f, err := os.Create("./" + name + ".wav")
	if err != nil {
		return nil, err
	}
	log.Println("writing file: ", f.Name())
	_, err = f.Write(b)
	f.Close()
	defer removeFile(f.Name())
	if err != nil {
		return nil, err
	}
	log.Println("converting to amr: ", f.Name())
	// Convert to AMR
	comm := ffmpeg_go.Input(name+".wav").Output(name+".amr", ffmpeg_go.KwArgs{"ar": AudioSamplingRateAMR, "ab": AudioBitRate})
	//comm := exec.Command("ffmpeg", "-i", "./"+name+".wav", "-ab", AudioBitRate, "-ar", AudioSamplingRateAMR, name+".amr")
	if err = comm.Run(); err != nil {
		return nil, err
	}
	defer removeFile(name + ".amr")
	data, err = os.ReadFile(name + ".amr")
	log.Println("successfully converted to amr")
	return
}

func removeFile(name string) {
	if err := os.Remove(name); err != nil {
		fmt.Printf("failed to remove file %s: %v\n", name, err)
	}
}
