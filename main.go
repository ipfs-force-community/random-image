package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/image/draw"
)

func main() {
	app := cli.NewApp()
	app.Name = "Image Generator"
	app.Usage = "Generate images with noise"
	app.Version = "1.0.0"
	app.Commands = []*cli.Command{
		generateCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return
	}
}

var generateCommand = &cli.Command{
	Name:    "generate",
	Aliases: []string{"g"},
	Usage:   "Generate random images",
	Flags: []cli.Flag{
		// &cli.IntFlag{
		// 	Name:    "count",
		// 	Aliases: []string{"c"},
		// 	Usage:   "Number of images to generate",
		// 	Value:   100,
		// },
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output directory for generated images",
			Value:   "./data",
		},
		&cli.StringFlag{
			Name:     "source-dir",
			Usage:    "Directory containing source images for generation",
			Required: true,
		},
		&cli.IntFlag{
			Name:    "width",
			Aliases: []string{"w"},
			Usage:   "Width and height of the generated images",
			Value:   3000,
		},
	},
	Action: func(c *cli.Context) error {
		// count := c.Int("count")
		outputDir := c.String("output")
		sourceDir := c.String("source-dir")
		width := c.Int("width")
		height := width // Assuming square images

		// 设置随机种子
		rand.New(rand.NewSource(time.Now().UnixNano()))

		fmt.Printf("output directory: %s, source directory: %v\n", outputDir, sourceDir)

		fi, err := os.Stat(outputDir)
		if err != nil {
			if os.IsNotExist(err) {
				// If the directory does not exist, create it
				err = os.MkdirAll(outputDir, 0755)
				if err != nil {
					return fmt.Errorf("failed to create output directory: %v", err)
				}
			} else {
				return fmt.Errorf("failed to access output directory: %v", err)
			}
		} else if !fi.IsDir() {
			return fmt.Errorf("output path is not a directory: %s", outputDir)
		}

		// 从源目录中随机选择一张图片进行处理
		files, err := os.ReadDir(sourceDir)
		if err != nil {
			return fmt.Errorf("failed to read source directory: %v", err)
		}
		if len(files) == 0 {
			return fmt.Errorf("source directory is empty: %s", sourceDir)
		}
		for _, f := range files {
			// 生成随机图片
			outputPath := fmt.Sprintf("%s/%s", outputDir, f.Name())
			inputPath := fmt.Sprintf("%s/%s", sourceDir, f.Name())
			// fmt.Printf("Processing source image: %s\n", inputPath)
			err = resizeImage(inputPath, outputPath, randWidth(width), randWidth(height))
			if err != nil {
				return fmt.Errorf("failed to generate image from source: %v", err)
			}
		}
		fmt.Printf("Successfully generated %d random images in directory: %s\n", len(files), outputDir)

		return nil
	},
}

func randWidth(in int) int {
	// 生成一个随机宽度，范围在原始宽度的±50%之间
	min := in - in/2
	max := in + in/2
	return rand.Intn(max-min+1) + min
}

func resizeImage(inputPath, outputPath string, targetWidth, targetHeight int) error {
	// 打开原始图片
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input image: %v", err)
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// 获取原始图片尺寸
	bounds := img.Bounds()
	// origWidth, origHeight := bounds.Dx(), bounds.Dy()
	// fmt.Printf("原始图片尺寸: %dx%d\n", origWidth, origHeight)

	// 创建中间图像，应用变换（添加随机噪声）
	intermediateImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取原始像素
			r, g, b, a := img.At(x, y).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8 // 转换为0-255

			// 添加随机噪声（±50范围，保留原始色调）
			rNoise := int(r) + rand.Intn(100) - 50
			gNoise := int(g) + rand.Intn(100) - 50
			bNoise := int(b) + rand.Intn(100) - 50

			// 限制颜色值在0-255
			if rNoise < 0 {
				rNoise = 0
			} else if rNoise > 255 {
				rNoise = 255
			}
			if gNoise < 0 {
				gNoise = 0
			} else if gNoise > 255 {
				gNoise = 255
			}
			if bNoise < 0 {
				bNoise = 0
			} else if bNoise > 255 {
				bNoise = 255
			}

			// 设置新像素
			intermediateImg.Set(x, y, color.RGBA{uint8(rNoise), uint8(gNoise), uint8(bNoise), uint8(a)})
		}
	}

	// 目标尺寸（8000x8000，RGBA ≈ 244MB，实际PNG约50-100MB）
	// targetWidth, targetHeight := 8000, 8000

	// 创建目标大图片
	largeImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// 缩放中间图像到目标尺寸
	draw.CatmullRom.Scale(largeImg, largeImg.Bounds(), intermediateImg, intermediateImg.Bounds(), draw.Over, nil)

	// 保存为PNG（无压缩）
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	encoder := &png.Encoder{CompressionLevel: png.NoCompression}
	if err := encoder.Encode(outFile, largeImg); err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}

	// 检查文件大小
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024) // 转换为MB
	fmt.Printf("图片已生成，路径：%s，大小：%.2f MB\n", outputPath, fileSizeMB)

	return nil
}

func genRandomImage(outputPath string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))
			a := uint8(255)
			img.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	// 保存为PNG（无压缩）
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := &png.Encoder{CompressionLevel: png.NoCompression}
	if err := encoder.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}

	// 检查文件大小
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024) // 转换为MB
	fmt.Printf("图片已生成，路径：%s，大小：%.2f MB\n", outputPath, fileSizeMB)

	return nil
}
