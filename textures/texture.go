package textures

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.6-core/gl"
	"image"
	"log"
	"os"
)

func NewTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("não foi possível abrir o arquivo de imagem: %v", err)
	}
	defer func(imgFile *os.File) {
		err := imgFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(imgFile)

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, fmt.Errorf("não foi possível decodificar a imagem: %v", err)
	}

	flippedImage := imaging.FlipV(img)

	rgba := imaging.Clone(flippedImage)
	width, height := int32(rgba.Bounds().Dx()), int32(rgba.Bounds().Dy())

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return textureID, nil
}
