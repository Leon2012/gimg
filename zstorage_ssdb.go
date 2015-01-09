package gimg

import (
	"fmt"
	"github.com/gographics/imagick/imagick"
)

type ZSSDBStorage struct {
	*ZBaseStorage
}

func NewSSDBStorage(c *ZContext) *ZSSDBStorage {
	z := new(ZSSDBStorage)
	z.ZBaseStorage = new(ZBaseStorage)
	z.ZBaseStorage.context = c
	return z
}

func (z *ZSSDBStorage) SaveImage(data []byte) (string, error) {
	var result error = nil
	md5Sum := gen_md5_str(data)
	z.context.Logger.Info("md5 : %s", md5Sum)

	if z.context.Redis.Exist(md5Sum) {
		result = fmt.Errorf("File Exist, Needn't Save.")
		return "", result
	}

	z.context.Logger.Debug("exist_db not found. Begin to Save File.")
	result = z.context.Redis.Send("SET", md5Sum, data)
	if result != nil {
		return "", result
	}

	return md5Sum, nil
}

func (z *ZSSDBStorage) GetImage(request *ZRequest) ([]byte, error) {
	var result error = nil
	var data []byte = nil
	var rspCachekey string
	toSave := true

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	z.context.Logger.Info("get_img() start processing zimg request...")

	md5Sum := request.Md5
	if !z.context.Redis.Exist(md5Sum) {
		result = fmt.Errorf("Image [%s] is not existed.", md5Sum)
		return nil, result
	}

	if len(request.ImageType) > 0 {
		rspCachekey = fmt.Sprintf("%s:%s", md5Sum, request.ImageType)
	} else {
		if request.Proportion == 0 && request.Width == 0 && request.Height == 0 {
			rspCachekey = md5Sum
		} else {
			rspCachekey = gen_key(md5Sum, request.Width, request.Height, request.Proportion, request.Gary, request.X, request.Y, request.Rotate, request.Quality, request.Format)
		}
	}

	data, result = z.context.Cache.FindCacheBin(rspCachekey)
	if result == nil {
		z.context.Logger.Debug("Hit Cache[Key: %s].\n", rspCachekey)
		toSave = false
		return data, nil
	}

	z.context.Logger.Debug("Start to Find the Image...")

	data, result = z.context.Redis.Get(rspCachekey)
	if result == nil {
		z.context.Logger.Debug("Get image [%s] from backend db succ.", rspCachekey)
		if len(data) < CACHE_MAX_SIZE {
			z.context.Cache.SetCacheBin(rspCachekey, data)
		}
		toSave = false
		return data, nil
	}

	data, result = z.context.Cache.FindCacheBin(md5Sum)
	if result != nil {
		data, result = z.context.Redis.Get(md5Sum)
		if result != nil {
			z.context.Logger.Debug("Get image [%s] from backend db failed.", md5Sum)
			return nil, result
		} else {
			if len(data) < CACHE_MAX_SIZE {
				z.context.Cache.SetCacheBin(md5Sum, data)
			}
		}
	}

	result = mw.ReadImageBlob(data)
	if result != nil {
		z.context.Logger.Debug("Webimg Read Blob Failed!")
		return nil, result
	}

	result = convert(mw, request)
	if result != nil {
		return nil, result
	}

	data = mw.GetImageBlob()
	z.context.Logger.Debug("get image blob length : %d", len(data))

	if data == nil || len(data) == 0 {
		result = fmt.Errorf("Webimg Get Blob Failed!")
		return nil, result
	}

	if len(data) < CACHE_MAX_SIZE {
		z.context.Cache.SetCacheBin(rspCachekey, data)
	}

	saveNew := 0
	if toSave {
		if request.Save == 1 || z.context.Config.Storage.SaveNew == 1 {
			saveNew = 1
		}
	}
	if saveNew == 1 {
		z.context.Logger.Debug("Image [%s] Saved to Storage.", rspCachekey)
		if err := z.context.Redis.Send("SET", rspCachekey, data); err != nil {
			z.context.Logger.Debug("New Image[%s] Save Failed!", rspCachekey)
		}
	} else {
		z.context.Logger.Debug("Image [%s] Needn't to Storage.", rspCachekey)
	}

	return data, nil
}

func (z *ZSSDBStorage) InfoImage(md5 string) (*ZImageInfo, error) {
	var result error
	result = nil
	z.context.Logger.Info("info_img() start processing info request...")

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	md5Sum := md5
	data, err := z.context.Redis.Get(md5Sum)
	if err != nil {
		result = fmt.Errorf("Image [%s] is not existed.", md5Sum)
		return nil, result
	}

	err = mw.ReadImageBlob(data)
	if err != nil {
		result = err
		return nil, result
	}

	size := 0
	width := int(mw.GetImageWidth())
	height := int(mw.GetImageHeight())
	quality := int(mw.GetImageCompressionQuality())
	format := mw.GetImageFormat()

	return &ZImageInfo{Size: size,
		Width:   width,
		Height:  height,
		Quality: quality,
		Format:  format,
	}, nil
}
