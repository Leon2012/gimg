# Gimg

----------
Gimg是[zimg](https://github.com/buaazp/zimg)的golang版本。

完全兼容zimg的文件储存格式。

环境要求：

* 操作系统： ubuntu/debian/osx
* golang版本： >= 1.4
* ImageMagick版本： = 6.8.5-X





# Install

----------
## Ubuntu/Debian

----------
	wget http://www.magickwand.org/download/releases/ImageMagick-6.8.5-10.tar.gz
	tar zxvf ImageMagick-6.8.5-10.tar.gz
	cd ImageMagick-6.8.5-10/
	./configure
	make & make install
	ldconfig /usr/local/lib
----------
## OSX

----------
	brew install ImageMagick
	
----------
## 安装
----------
	go get github.com/gographics/imagick/imagick
	go get code.google.com/p/gcfg
	go get github.com/leon2012/gimg
	cd $GOPATH/gimg/build/
	go build -o gimg
	./gimg --config=./conf/config.ini
	
	
	
	

	
	



