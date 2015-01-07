package gimg

import (
	"fmt"
	"github.com/gographics/imagick/imagick"
	"math"
)

func crop(mw *imagick.MagickWand, x, y int, cols, rows uint) error {
	var result error
	result = nil

	imCols := mw.GetImageWidth()
	imRows := mw.GetImageHeight()

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	if uint(x) >= imCols || uint(y) >= imRows {
		result = fmt.Errorf("x, y more than image cols, rows")
		return result
	}

	if cols == 0 || imCols < uint(x)+cols {
		cols = imCols - uint(x)
	}

	if rows == 0 || imRows < uint(y)+rows {
		rows = imRows - uint(y)
	}

	fmt.Print(fmt.Sprintf("wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows))

	result = mw.CropImage(cols, rows, x, y)

	return result
}

func proportion(mw *imagick.MagickWand, proportion int, cols uint, rows uint) error {
	var result error
	result = nil

	imCols := mw.GetImageWidth()
	imRows := mw.GetImageHeight()

	if proportion == 0 {
		fmt.Sprintf("p=0, wi_scale(im, %d, %d)\n", cols, rows)
		result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)

	} else if proportion == 1 {

		if cols == 0 || rows == 0 {
			if cols > 0 {
				rows = uint(round(float64((cols / imCols) * imRows)))
			} else {
				cols = uint(round(float64((rows / imRows) * imCols)))
			}
			fmt.Sprintf("p=1, wi_scale(im, %d, %d)\n", cols, rows)
			result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
		} else {
			var x, y, sCols, sRows uint
			x, y = 0, 0

			colsRate := cols / imCols
			rowsRate := rows / imRows

			if colsRate > rowsRate {
				sCols = cols
				sRows = uint(round(float64(colsRate * imRows)))
				y = uint(math.Floor(float64((sRows - rows) / 2.0)))
			} else {
				sCols = uint(round(float64(rowsRate * imCols)))
				sRows = rows
				x = uint(math.Floor(float64((sCols - cols) / 2.0)))
			}

			fmt.Sprintf("p=2, wi_scale(im, %d, %d)\n", sCols, sRows)
			result = mw.ResizeImage(sCols, sRows, imagick.FILTER_UNDEFINED, 1.0)

			fmt.Sprintf("p=2, wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows)
			result = mw.CropImage(cols, rows, int(x), int(y))
		}

	} else if proportion == 2 {
		x := int(math.Floor(float64((imCols - cols) / 2.0)))
		y := int(math.Floor(float64((imRows - rows) / 2.0)))
		fmt.Sprintf("p=3, wi_crop(im, %d, %d, %d, %d)\n", x, y, cols, rows)
		result = mw.CropImage(cols, rows, x, y)

	} else if proportion == 3 {
		if cols == 0 || rows == 0 {
			var rate uint
			if cols > 0 {
				rate = cols
			} else {
				rate = rows
			}
			rows = uint(round(float64(imRows * rate / 100)))
			cols = uint(round(float64(imCols * rate / 100)))
			fmt.Sprintf("p=3, wi_scale(im, %d, %d)\n", cols, rows)
			result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
		} else {
			rows = uint(round(float64(imRows * rows / 100)))
			cols = uint(round(float64(imCols * cols / 100)))
			fmt.Sprintf("p=3, wi_scale(im, %d, %d)\n", cols, rows)
			result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
		}

	} else if proportion == 4 {
		var rate float64
		rate = 1.0
		if cols == 0 || rows == 0 {
			if cols > 0 {
				rate = float64(cols / imCols)
			} else {
				rate = float64(rows / imRows)
			}
		} else {
			rateCol := cols / imCols
			rateRow := rows / imRows
			if rateCol < rateRow {
				rate = float64(rateCol)
			} else {
				rate = float64(rateRow)
			}
		}

		cols = uint(round(float64(float64(imCols) * rate)))
		rows = uint(round(float64(float64(imRows) * rate)))
		fmt.Sprintf("p=4, wi_scale(im, %d, %d)\n", cols, rows)
		result = mw.ResizeImage(cols, rows, imagick.FILTER_UNDEFINED, 1.0)
	}

	return result

}

func convert(mw *imagick.MagickWand, request *ZRequest) error {

	fmt.Println("call convert function......")

	var result error
	result = nil
	mw.ResetIterator()
	mw.SetImageOrientation(imagick.ORIENTATION_TOP_LEFT)

	x := request.X
	y := request.Y
	cols := uint(request.Width)
	rows := uint(request.Height)

	fmt.Sprintf("image cols %d, rows %d \n", cols, rows)

	if !(cols == 0 && rows == 0) {
		fmt.Println("call crop&scal function......")

		/* crop and scale */
		if x == -1 && y == -1 {
			fmt.Println("call crop&scal function......")

			fmt.Print(fmt.Sprintf("proportion(im, %d, %d, %d) \n", request.Proportion, cols, rows))
			result = proportion(mw, request.Proportion, cols, rows)
			if result != nil {
				return result
			}
		} else {

			fmt.Print(fmt.Sprintf("crop(im, %d, %d, %d, %d) \n", x, y, cols, rows))

			result = crop(mw, x, y, cols, rows)
			if result != nil {
				return result
			}
		}
	}

	/* rotate image */
	if request.Rotate != 0 {
		fmt.Print(fmt.Sprintf("wi_rotate(im, %d) \n", request.Rotate))

		background := imagick.NewPixelWand()
		if background == nil {
			result = fmt.Errorf("init new pixelwand faile.")
			return result
		}
		defer background.Destroy()
		isOk := background.SetColor("#ffffff")
		if !isOk {
			result = fmt.Errorf("set background color faile.")
			return result
		}

		result = mw.RotateImage(background, float64(request.Rotate))
		if result != nil {
			return result
		}
	}

	/* set gray */
	if request.Gary == 1 {
		fmt.Print(fmt.Sprintf("wi_gray(im) \n"))
		result = mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE)
		if result != nil {
			return result
		}
	}

	/* set quality */
	fmt.Print(fmt.Sprintf("wi_set_quality(im, %d) \n", request.Quality))
	result = mw.SetImageCompressionQuality(uint(request.Quality))
	if result != nil {
		return result
	}

	/* set format */
	if "none" != request.Format {
		fmt.Print(fmt.Sprintf("wi_set_format(im, %s) \n", request.Format))
		result = mw.SetImageFormat(request.Format)
		if result != nil {
			return result
		}
	}

	fmt.Print(fmt.Sprintf("convert(im, req) %s \n", result))

	return result
}
