package gimg

import (
	"fmt"
	"os"
	"sync"
)

type ZBaseStorage struct {
	path string
	sync.RWMutex
	context *ZContext
}

func (z *ZBaseStorage) SaveImage(data []byte) (string, error) {
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

func (z *ZBaseStorage) NewImage(saveName string, data []byte) error {
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
