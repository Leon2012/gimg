package gimg

import (
	"fmt"
	"github.com/gographics/imagick/imagick"
	"os"
	"sync"
)

type ZFileStorage struct {
	path string
	sync.RWMutex
	*ZBaseStorage
}

func NewFileStorage(c *ZContext) *ZFileStorage {
	imgPath := c.Config.Storage.ImgPath

	z := new(ZFileStorage)
	z.path = imgPath

	z.ZBaseStorage = new(ZBaseStorage)
	z.ZBaseStorage.context = c
	return z
}

func (z *ZFileStorage) SaveImage(data []byte) (string, error) {
	var result error
	result = nil

	md5Sum := gen_md5_str(data)
	var savePath string
	var saveName string

	lvl1 := str_hash(md5Sum)
	lvl2 := str_hash(string(md5Sum[3:]))

	savePath = fmt.Sprintf("%s/%d/%d/%s", z.path, lvl1, lvl2, md5Sum)
	z.context.Logger.Debug("save path: ", savePath)

	if is_dir(savePath) {
		//z.context.Logger.Info("Check File Exist. Needn't Save.")
		result = fmt.Errorf("Check File Exist. Needn't Save.")
		goto cache
	}

	if !mk_dir(savePath) {
		result = fmt.Errorf("save path[%s] Create Failed!", savePath)
		goto done
	}

	z.context.Logger.Debug("save path[%s] Create Finish.\n", savePath)
	saveName = fmt.Sprintf("%s/0*0", savePath)
	z.context.Logger.Debug("save name-->: %s", saveName)

	result = z.NewImage(saveName, data)
	goto done

cache:
	if len(data) < CACHE_MAX_SIZE {
		z.context.Cache.SetCacheBin(md5Sum, data)
		z.context.Logger.Info("save " + md5Sum + " into cache")
	}
	//result = nil
	return md5Sum, result

done:
	return md5Sum, result
}

func (z *ZFileStorage) NewImage(saveName string, data []byte) error {
	var result error
	result = nil

	z.context.Logger.Info("Start to Storage the New Image...")

	f, result := os.OpenFile(saveName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if result != nil {
		return result
	}

	z.Lock() //文件读写锁
	defer z.Unlock()
	defer f.Close()

	_, err := f.Write(data)
	if err != nil {
		result = err
	}

	return result
}

func (z *ZFileStorage) GetImage(request *ZRequest) ([]byte, error) {
	var result error
	var data []byte
	var rspCachekey string
	result = nil
	data = nil
	toSave := true
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	z.context.Logger.Info("get_img() start processing zimg request...")

	md5Sum := request.Md5
	lvl1 := str_hash(md5Sum)
	lvl2 := str_hash(string(md5Sum[3:]))

	wholePath := fmt.Sprintf("%s/%d/%d/%s", z.path, lvl1, lvl2, md5Sum)
	z.context.Logger.Debug("whole path : ", wholePath)

	if !is_dir(wholePath) {
		result = fmt.Errorf("Image %s is not existed!", md5Sum)
		//goto ErrHandle
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

	data, err := z.context.Cache.FindCacheBin(rspCachekey)
	if err == nil {
		z.context.Logger.Debug("Hit Cache[Key: %s].\n", rspCachekey)
		toSave = false

		//goto Done
		return data, nil
	}

	z.context.Logger.Info("Start to Find the Image...")
	origPath := fmt.Sprintf("%s/0*0", wholePath)
	z.context.Logger.Debug("0rig File Path: %s \n", origPath)

	rspPath := ""
	if len(request.ImageType) > 0 {
		rspPath = fmt.Sprintf("%s/t_%s", wholePath, request.ImageType)
	} else {
		name := fmt.Sprintf("%d*%d_p%d_g%d_%d*%d_r%d_q%d.%s", request.Width,
			request.Height, request.Proportion, request.Gary, request.X, request.Y, request.Rotate, request.Quality, request.Format)

		if request.Proportion == 0 && request.Width == 0 && request.Height == 0 {
			z.context.Logger.Info("Return original image.")
			rspPath = origPath
		} else {
			rspPath = fmt.Sprintf("%s/%s", wholePath, name)
		}
	}

	z.context.Logger.Debug("Got the rsp_path: %s \n", rspPath)

	f, err := os.OpenFile(rspPath, os.O_RDONLY, 0755)
	defer f.Close()

	if err != nil { //读取不到文件
		origData, err := z.context.Cache.FindCacheBin(md5Sum)
		if err == nil {
			z.context.Logger.Debug("Hit Orignal Image Cache[Key: %s].", md5Sum)
			err = mw.ReadImageBlob(origData)
			if err != nil {
				z.context.Logger.Error("Open Original Image From Blob Failed! Begin to Open it From Disk.")
				z.context.Cache.DelCache(md5Sum)

				// origFile, err := os.Open(origPath)
				// if err != nil {
				// 	result = fmt.Errorf("Open Original Image From Disk Failed!")
				// 	//goto ErrHandle
				// 	return nil, result
				// }

				//err = mw.ReadImageFile(origFile)
				err = mw.ReadImage(origPath)
				if err != nil {
					result = fmt.Errorf("Open Original Image From Disk Failed!")
					//goto ErrHandle
					return nil, result
				} else {
					mw.ResetIterator()
					newBuff := mw.GetImageBlob()
					if newBuff == nil {
						result = fmt.Errorf("Webimg Get Original Blob Failed!")
						//goto ErrHandle
						return nil, result
					}

					if len(newBuff) < CACHE_MAX_SIZE {
						z.context.Cache.SetCacheBin(md5Sum, newBuff)
						z.context.Logger.Info("save " + md5Sum + " into cache")
					}

				}

			}

		} else {
			z.context.Logger.Debug("Not Hit Original Image Cache. Begin to Open it.")
			// origFile, err := os.Open(origPath)
			// if err != nil {
			// 	result = fmt.Errorf("Open Original Image From Disk Failed!")
			// 	//goto ErrHandle
			// 	return nil, result
			// }

			// err = mw.ReadImageFile(origFile)
			err = mw.ReadImage(origPath)
			if err != nil {
				result = fmt.Errorf("Open Original Image From Disk Failed!")
				//goto ErrHandle
				return nil, result
			} else {
				mw.ResetIterator()
				newBuff := mw.GetImageBlob()
				//z.context.Logger.Debug("get image blob length1111 : %d", len(newBuff))

				if newBuff == nil {
					result = fmt.Errorf("Webimg Get Original Blob Failed!")
					//goto ErrHandle
					return nil, result
				}
				if len(newBuff) < CACHE_MAX_SIZE {
					z.context.Cache.SetCacheBin(md5Sum, newBuff)
					z.context.Logger.Info("save " + md5Sum + " into cache")
				}
			}
		}

		// newBuff1 := mw.GetImageBlob()
		// z.context.Logger.Debug("get image blob length2222 : %d", len(newBuff1))

		result = convert(mw, request)
		if result != nil {
			return nil, result
		}

		data = mw.GetImageBlob()
		z.context.Logger.Debug("get image blob length : %d", len(data))

		if data == nil || len(data) == 0 {
			result = fmt.Errorf("Webimg Get Blob Failed!")
			//goto ErrHandle
			return nil, result
		}

	} else {
		toSave = false
		fi, err := f.Stat()
		flen := int(fi.Size())
		if flen <= 0 {
			//z.context.Logger.Debug("File[%s] is Empty.", rspPath)
			result = fmt.Errorf("File[%s] is Empty.", rspPath)
			//goto ErrHandle
			return nil, result
		}

		data = make([]byte, flen)
		rlen, err := f.Read(data)
		if err != nil {
			result = err
			//goto ErrHandle
			return nil, result
		}

		if rlen < flen {
			result = fmt.Errorf("File[%s] Read Not Compeletly. file len : %d, read len :%d", rspPath, flen, rlen)
			//goto ErrHandle
			return nil, result
		}
	}

	saveNew := 0
	if toSave {
		if request.Save == 1 || z.context.Config.Storage.SaveNew == 1 {
			saveNew = 1
		}
	}

	if saveNew == 1 {
		z.context.Logger.Debug("Image[%s] is Not Existed. Begin to Save it.", rspPath)
		if err = z.NewImage(rspPath, data); err != nil {
			z.context.Logger.Debug("New Image[%s] Save Failed!", rspPath)
			z.context.Logger.Warning("fail save %s", rspPath)
		}
	} else {
		z.context.Logger.Debug("Image [%s] Needn't to Storage.", rspPath)
	}

	if len(data) < CACHE_MAX_SIZE {
		z.context.Cache.SetCacheBin(rspCachekey, data)
		z.context.Logger.Info("save " + rspCachekey + " into cache")
	}

	//fmt.Println("image data :", data)

	return data, nil

	// ErrHandle:
	// 	return nil, result

	// Done:
	// 	return data, nil
}

func (z *ZFileStorage) InfoImage(md5 string) (*ZImageInfo, error) {
	var result error
	result = nil
	z.context.Logger.Info("info_img() start processing info request...")

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	md5Sum := md5
	lvl1 := str_hash(md5Sum)
	lvl2 := str_hash(string(md5Sum[3:]))

	wholePath := fmt.Sprintf("%s/%d/%d/%s", z.path, lvl1, lvl2, md5Sum)
	z.context.Logger.Debug("whole_path: %s", wholePath)

	if !is_dir(wholePath) {
		result = fmt.Errorf("Image %s is not existed!", md5Sum)

		return nil, result
	}

	origPath := fmt.Sprintf("%s/0*0", wholePath)
	z.context.Logger.Debug("0rig File Path: %s", origPath)

	err := mw.ReadImage(origPath)

	// origFile, err := os.Open(origPath)
	// if err != nil {
	// 	result = err
	// 	return nil, result
	// }

	// err = mw.ReadImageFile(origFile)
	if err != nil {
		result = fmt.Errorf("Open Original Image From Disk Failed!")
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
