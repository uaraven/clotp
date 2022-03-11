package cli

import (
	"image"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type ScanCmd struct {
	Image  string `arg:"positional,required" help:"Image with the QR Code"`
	Decode bool   `arg:"-d,--decode" help:"Decode otpauth-migration URI retrieved from the QR code"`
	Parse  bool   `arg:"-p,--parse" help:"For decoded migration URIs, parse them and print each part separately"`
}

func ScanQrCode(cmd *ScanCmd) (string, error) {
	file, err := os.Open(cmd.Image)
	if err != nil {
		return "", err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	qrUri := result.GetText()
	if cmd.Decode {
		uris, err := decodeMigrationUri(qrUri, cmd.Parse)
		if err != nil {
			return "", err
		}
		var sb strings.Builder
		for _, uri := range uris {
			sb.WriteString(uri)
		}
		return sb.String(), nil
	}
	return result.GetText(), nil
}
