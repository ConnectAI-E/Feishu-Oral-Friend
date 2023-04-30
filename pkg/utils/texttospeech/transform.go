package texttospeech

import (
	"fmt"

	"github.com/go-tts/tts/pkg/speech"
)

func Transform(q string) {
	audioIn,_ := speech.FromText("text", speech.LangEn)
	fmt.Println(audioIn)
}