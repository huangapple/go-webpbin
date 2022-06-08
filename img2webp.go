package webpbin

import (
	"errors"
	"fmt"
	"io"

	"github.com/nickalie/go-binwrapper"
)

// Img2Webp compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/Img2Webp

type Img2WebpFrame struct {
	Url string
	//Lossless bool
	//Lossy    bool
	D int //指定图像的持续时间， 默认100ms
	//Q        int //0~100 压缩因子， 默认75
	//M        int //0~6 默认4
}

type Img2Webp struct {
	*binwrapper.BinWrapper
	outputFile string
	output     io.Writer

	//kmin   int              //minimum number of frame between key-frames (0=disable key-frames altogether)
	//kmax   int              //maximum number of frame between key-frames (0=only keyframes)
	//loop   int              //循环次数
	frames []*Img2WebpFrame //图片帧参数
	mixed  bool             //use mixed lossy/lossless automatic mode
}

// NewImg2Webp creates new Img2Webp instance.
func NewImg2Webp(optionFuncs ...OptionFunc) *Img2Webp {
	bin := &Img2Webp{
		BinWrapper: createBinWrapper(optionFuncs...),
		mixed:      true,
	}
	bin.ExecPath("img2webp")

	return bin
}

// Version returns Img2Webp version.
func (c *Img2Webp) Version() (string, error) {
	return version(c.BinWrapper)
}

// OutputFile specify the name of the output WebP file.
// Output called before will be ignored.
func (c *Img2Webp) OutputFile(file string) *Img2Webp {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *Img2Webp) Output(writer io.Writer) *Img2Webp {
	c.outputFile = ""
	c.output = writer
	return c
}

// Quality specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
// A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *Img2Webp) SetFrames(frames []*Img2WebpFrame) *Img2Webp {
	c.frames = frames
	return c
}

func (c *Img2Webp) Mixed(mixed bool) *Img2Webp {
	c.mixed = mixed
	return c
}

// Run starts Img2Webp with specified parameters.
func (c *Img2Webp) Run() error {
	defer c.BinWrapper.Reset()

	if c.mixed {
		c.Arg("-mixed")
	}
	err := c.setInput()

	if err != nil {
		return err
	}
	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

func (c *Img2Webp) setInput() error {

	for _, frame := range c.frames {
		c.Arg(frame.Url)
		if frame.D != 0 {
			c.Arg("-d")
			c.Arg(fmt.Sprint(frame.D))
		}
	}

	return nil
}

func (c *Img2Webp) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
