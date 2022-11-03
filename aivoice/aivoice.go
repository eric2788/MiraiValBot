package aivoice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

const (
	AudioSamplingRateMP3  = "22050"
	AudioBitRate          = "12.2k" // in Hz
	NumberOfAudioChannels = "1"
	AudioSamplingRateAMR  = "8000"
)

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
