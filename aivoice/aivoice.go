package aivoice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"

	"github.com/Logiase/MiraiGo-Template/utils"
	silk "github.com/wdvxdr1123/go-silk"
)

var logger = utils.GetModuleLogger("valbot.aivoice")

const (
	AudioBitRate          = "12.2k" // in Hz
	AudioSamplingRateAMR  = "8000"
)

func WavToSilk(b []byte) (data []byte, err error) {
	md := md5.Sum(b)
	tempName := hex.EncodeToString(md[:])

	wav, pcm := tempName+".wav", tempName+".pcm"
	logger.Infof("writing file: %s", wav)
	err = os.WriteFile(wav, b, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer os.Remove(wav)

	// 2.转换pcm
	logger.Infof("converting to pcm...")
	cmd := exec.Command("ffmpeg", "-i", wav, "-f", "s16le", "-ar", "24000", "-ac", "1", pcm)
	logger.Infof("converted to pcm: %s", pcm)
	if err = cmd.Run(); err != nil {
		return nil, err
	}
	defer os.Remove(pcm)
	pcmByte, err := os.ReadFile(pcm)
	if err != nil {
		return nil, err
	}
	logger.Infof("converting to silk...")
	defer func ()  {
		if err == nil {
			logger.Infof("converted to silk: %s", tempName+".silk")
		}		
	}()
	data, err = silk.EncodePcmBuffToSilk(pcmByte, 24000, 24000, true)
	return
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
	logger.Info("writing file: ", f.Name())
	_, err = f.Write(b)
	f.Close()
	defer removeFile(f.Name())
	if err != nil {
		return nil, err
	}
	logger.Info("converting to amr: ", f.Name())
	// Convert to AMR
	comm := exec.Command("ffmpeg", "-i", "./"+name+".wav", "-ab", AudioBitRate, "-ar", AudioSamplingRateAMR, name+".amr")
	if err = comm.Run(); err != nil {
		return nil, err
	}
	defer removeFile(name + ".amr")
	data, err = os.ReadFile(name + ".amr")
	logger.Info("successfully converted to amr")
	return
}

func removeFile(name string) {
	if err := os.Remove(name); err != nil {
		fmt.Printf("failed to remove file %s: %v\n", name, err)
	}
}
