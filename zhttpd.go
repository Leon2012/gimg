package gimg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type ZHttpd struct {
	context *ZContext
	//storage      *ZStorageFile
	storage      ZStorage
	writer       http.ResponseWriter
	request      *http.Request
	contentTypes map[string]string
}

func NewHttpd(c *ZContext) *ZHttpd {
	var s ZStorage
	if c.Config.Storage.Mode == 1 {
		s = NewFileStorage(c)
	}

	return &ZHttpd{context: c, storage: s, contentTypes: genContentTypes()}
}

func genContentTypes() map[string]string {
	types := make(map[string]string)
	types["jpg"] = "image/jpeg"
	types["jpeg"] = "image/jpeg"
	types["png"] = "image/png"
	types["gif"] = "image/gif"
	types["webp"] = "image/webp"

	return types
}

func (z *ZHttpd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	z.writer = w
	z.request = r
	path := r.URL.Path
	method := r.Method

	if "GET" == method {
		if path == "/" {
			z.doDefault()
		} else if path == "/info" {
			z.doInfo()
		} else {
			z.context.Logger.Info("path:" + path)
			md5Sum := path[1:len(path)]
			z.context.Logger.Info("md5Sum:" + md5Sum)

			if is_md5(md5Sum) {
				z.doGet(md5Sum)
			} else {
				http.NotFound(w, r)
			}
		}

	} else if "POST" == method {
		if path == "/upload" {
			z.doUpload()
		} else {
			http.NotFound(w, r)
		}

	} else {
		http.NotFound(w, r)
	}

	return
}

func (z *ZHttpd) doDefault() {
	z.context.Logger.Info("call doDefault function........")

	html := `<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8"/>
    </head>
    <body>
        <form action="/upload" method="POST" enctype="multipart/form-data">

            <label for="field1">file:</label>
            <input name="upload_file" type="file" />
            <input type="submit"></input>

        </form>
    </body>
</html>`
	fmt.Fprint(z.writer, html)

}

func (z *ZHttpd) doInfo() {
	z.context.Logger.Info("call doInfo function........")
	if err := z.request.ParseForm(); err != nil {
		z.context.Logger.Error(err.Error())
		z.doError(err, http.StatusForbidden)
		return
	}

	md5Sum := z.request.Form.Get("md5")
	z.context.Logger.Info("search md5  : %s", md5Sum)

	imgInfo, err := z.storage.InfoImage(md5Sum)
	if err != nil {
		z.context.Logger.Error(err.Error())
		z.doError(err, http.StatusForbidden)
		return
	}

	json, _ := json.Marshal(imgInfo)
	fmt.Fprint(z.writer, string(json))

}

func (z *ZHttpd) doUpload() {
	z.context.Logger.Info("call doUpload function........")

	if err := z.request.ParseMultipartForm(CACHE_MAX_SIZE); err != nil {
		z.context.Logger.Error(err.Error())
		z.doError(err, http.StatusForbidden)
		return
	}

	file, _, err := z.request.FormFile("upload_file")
	if err != nil {
		z.doError(err, 500)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		z.doError(err, 500)
		return
	}

	md5Sum, err := z.storage.SaveImage(data)
	if err != nil {
		z.doError(err, 500)
		return
	}

	fmt.Fprint(z.writer, fmt.Sprintf("upload success! md5 : %s", md5Sum))

}

func (z *ZHttpd) doGet(md5Sum string) {
	z.context.Logger.Info("call doGet function........")
	if err := z.request.ParseForm(); err != nil {
		z.context.Logger.Error(err.Error())
		z.doError(err, http.StatusForbidden)
		return
	}

	imgInfo, err := z.storage.InfoImage(md5Sum)
	if err != nil {
		z.context.Logger.Error(err.Error())
		z.doError(err, http.StatusForbidden)
		return
	}

	var w, h, p, g, x, y, r, s, q int = 0, 0, 0, 0, 0, 0, 0, 0, 0
	var i, f string

	width := z.request.Form.Get("w")
	height := z.request.Form.Get("h")
	gary := z.request.Form.Get("g")
	xx := z.request.Form.Get("x")
	yy := z.request.Form.Get("y")
	rotate := z.request.Form.Get("r")

	w = str2Int(width)
	if w >= imgInfo.Width || w <= 0 {
		w = imgInfo.Width
	}

	h = str2Int(height)
	if h >= imgInfo.Height || h <= 0 {
		h = imgInfo.Height
	}

	g = str2Int(gary)
	if g != 1 {
		g = 0
	}

	x = str2Int(xx)
	if x < 0 {
		x = -1
	}

	// else if x > imgInfo.Width {
	// 	x = imgInfo.Width
	// }

	y = str2Int(yy)
	if y < 0 {
		y = -1
	}
	// else if y > imgInfo.Height {
	// 	y = imgInfo.Height
	// }

	r = str2Int(rotate)

	quality := z.request.Form.Get("q")
	q = str2Int(quality)
	if q <= 0 {
		//q = imgInfo.Quality
		q = z.context.Config.System.Quality //加载默认保存图片质量
	} else if q > 100 {
		q = 100
	}

	save := strings.Trim(z.request.Form.Get("s"), " ")
	if len(save) == 0 {
		s = z.context.Config.Storage.SaveNew
	} else {
		s = str2Int(save)
		if s != 1 {
			s = 0
		}
	}

	format := strings.Trim(z.request.Form.Get("f"), " ")
	if len(format) == 0 {
		//f = "none"
		//f = imgInfo.Format
		f = z.context.Config.System.Format //加载默认保存图片格式
	} else {
		format = strings.ToLower(format)
		formats := strings.Split(z.context.Config.Storage.AllowedTypes, ",")
		isExist := false
		for _, v := range formats {
			if format == v {
				isExist = true
			}
		}
		if !isExist {
			f = z.context.Config.System.Format
		} else {
			f = format
		}
	}
	// if f == strings.ToLower(imgInfo.Format) {
	// 	f = "none"
	// }

	request := &ZRequest{
		Md5:        md5Sum,
		Width:      w,
		Height:     h,
		Gary:       g,
		X:          x,
		Y:          y,
		Rotate:     r,
		Quality:    q,
		Proportion: p,
		Save:       s,
		Format:     f,
		ImageType:  i,
	}

	z.context.Logger.Debug("request params: md5 : %s, width: %d, height: %d, gary: %d, x: %d, y: %d, rotate: %d, quality: %d, proportion: %d, save: %d, format: %s, imageType: %s", request.Md5, request.Width, request.Height, request.Gary, request.X, request.Y, request.Rotate, request.Quality, request.Proportion, request.Save, request.Format, request.ImageType)

	data, err := z.storage.GetImage(request)

	if err != nil {
		z.doError(err, 500)
		return
	}

	headers := z.context.Config.System.Headers
	if len(headers) > 0 {
		arr := strings.Split(headers, ",")
		for i := 0; i < len(arr); i++ {
			header := arr[i]
			kvs := strings.Split(header, ":")
			z.writer.Header().Set(kvs[0], kvs[1])
		}
	}

	//etag := z.context.Config.System.Etag

	imageFormat := strings.ToLower(f)
	if contentType, ok := z.contentTypes[imageFormat]; ok {
		z.writer.Header().Set("Content-Type", contentType)
		z.writer.Write(data)

	} else {
		err = fmt.Errorf("can not found content type!!!")
		z.doError(err, http.StatusForbidden)
		return
	}
}

func (z *ZHttpd) doError(err error, statusCode int) {
	http.Error(z.writer, err.Error(), statusCode)
	return
}

func str2Int(str string) int {
	str = strings.Trim(str, " ")
	if len(str) > 0 {
		i, err := strconv.Atoi(str)
		if err != nil {
			return 0
		} else {
			return i
		}
	} else {
		return 0
	}
}
