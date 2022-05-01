package webpbin

import (
	"errors"
	"fmt"
	"image/gif"
	"io"

	"github.com/nickalie/go-binwrapper"
)

// Gif2WebP compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/Gif2WebP
type Gif2WebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	inputGif   *gif.GIF
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
	crop       *cropInfo
	mixed      bool
}

// NewGif2WebP creates new Gif2WebP instance.
func NewGif2WebP(optionFuncs ...OptionFunc) *Gif2WebP {
	bin := &Gif2WebP{
		BinWrapper: createBinWrapper(optionFuncs...),
		quality:    -1,
	}
	bin.ExecPath("gif2webp")

	return bin
}

// Version returns Gif2WebP version.
func (c *Gif2WebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *Gif2WebP) InputFile(file string) *Gif2WebP {
	c.input = nil
	c.inputGif = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *Gif2WebP) Input(reader io.Reader) *Gif2WebP {
	c.inputFile = ""
	c.inputGif = nil
	c.input = reader
	return c
}

// InputImage sets image to convert.
// InputFile or Input called before will be ignored.
func (c *Gif2WebP) InputGif(img *gif.GIF) *Gif2WebP {
	c.inputFile = ""
	c.input = nil
	c.inputGif = img
	return c
}

// OutputFile specify the name of the output WebP file.
// Output called before will be ignored.
func (c *Gif2WebP) OutputFile(file string) *Gif2WebP {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *Gif2WebP) Output(writer io.Writer) *Gif2WebP {
	c.outputFile = ""
	c.output = writer
	return c
}

// Quality specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
// A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *Gif2WebP) Quality(quality uint) *Gif2WebP {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

func (c *Gif2WebP) Mixed(mixed bool) *Gif2WebP {
	c.mixed = mixed
	return c
}

// Crop the source to a rectangle with top-left corner at coordinates (x, y) and size width x height. This cropping area must be fully contained within the source rectangle.
func (c *Gif2WebP) Crop(x, y, width, height int) *Gif2WebP {
	c.crop = &cropInfo{x, y, width, height}
	return c
}

// Run starts Gif2WebP with specified parameters.
func (c *Gif2WebP) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.Arg("-q", fmt.Sprintf("%d", c.quality))
	}

	if c.crop != nil {
		c.Arg("-crop", fmt.Sprintf("%d", c.crop.x), fmt.Sprintf("%d", c.crop.y), fmt.Sprintf("%d", c.crop.width), fmt.Sprintf("%d", c.crop.height))
	}

	if c.mixed {
		c.Arg("-mixed")
	}

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

// Reset all parameters to default values
func (c *Gif2WebP) Reset() *Gif2WebP {
	c.crop = nil
	c.quality = -1
	return c
}

func (c *Gif2WebP) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputGif != nil {
		r, err := createReaderFromGif(c.inputGif)

		if err != nil {
			return err
		}

		c.Arg("--").Arg("-")
		c.StdIn(r)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *Gif2WebP) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
