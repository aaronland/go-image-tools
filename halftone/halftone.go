package halftone

// https://maxhalford.github.io/blog/halftoning-1/

import (
	"errors"
	"github.com/MaxHalford/halfgone"
	"github.com/nfnt/resize"
	"image"
)

type HalftoneOptions struct {
	Mode        string
	ScaleFactor float64
}

func NewDefaultHalftoneOptions() HalftoneOptions {

	opts := HalftoneOptions{
		Mode:        "atkinson",
		ScaleFactor: 1.0,
	}

	return opts
}

func Halftone(im image.Image, opts HalftoneOptions) (image.Image, error) {

	// see notes below (20180927/thisisaaronland)

	var grey *image.Gray
	var w uint
	var h uint

	if opts.ScaleFactor > 0.0 {

		dims := im.Bounds()
		w = uint(dims.Max.X)
		h = uint(dims.Max.Y)

		scale_w := uint(float64(w) / opts.ScaleFactor)
		scale_h := uint(float64(h) / opts.ScaleFactor)

		thumb := resize.Thumbnail(scale_w, scale_h, im, resize.Lanczos3)
		grey = halfgone.ImageToGray(thumb)

	} else {
		grey = halfgone.ImageToGray(im)
	}

	switch opts.Mode {
	case "atkinson":
		grey = halfgone.AtkinsonDitherer{}.Apply(grey)
	case "threshold":
		grey = halfgone.ThresholdDitherer{Threshold: 127}.Apply(grey)
	default:
		return nil, errors.New("Invalid or unsupported mode")
	}

	// the resize process ends up making a greyscale image - not sure
	// what the best way to deal with this is (20180927/thisisaaonland)

	if opts.ScaleFactor > 0.0 {

		dither := resize.Resize(w, h, grey, resize.Lanczos3)
		grey = halfgone.ImageToGray(dither)

		switch opts.Mode {
		case "atkinson":
			grey = halfgone.AtkinsonDitherer{}.Apply(grey)
		case "threshold":
			grey = halfgone.ThresholdDitherer{Threshold: 127}.Apply(grey)
		default:
			return nil, errors.New("Invalid or unsupported mode")
		}
	}

	return grey, nil
}
