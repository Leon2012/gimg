# Gimg

- - - 
Gimg是[zimg](https://github.com/buaazp/zimg)的golang版本。

完全兼容zimg的文件目录储存格式，支持文件和类Redis协议(SSDB)储存。

环境要求：

* 操作系统： ubuntu/debian/osx
* golang版本： >= 1.4
* ImageMagick版本： = 6.8.5-X



# Install

- - -
## Ubuntu/Debian

- - - 
	wget http://www.magickwand.org/download/releases/ImageMagick-6.8.5-10.tar.gz
	tar zxvf ImageMagick-6.8.5-10.tar.gz
	cd ImageMagick-6.8.5-10/
	./configure
	make & make install
	ldconfig /usr/local/lib
- - -
## OSX

- - -
	brew install ImageMagick
	
- - -
## 安装
- - -
	go get github.com/gographics/imagick/imagick
	go get code.google.com/p/gcfg
	go get github.com/garyburd/redigo/redis
	go get github.com/Leon2012/gimg
	cd $GOPATH/src/github.com/Leon2012/gimg/build/
	go build -o gimg
	./gimg --config=./conf/config.ini
	
	
## 演示
- - -
[http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93](http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93)

[http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93?w=100&h=100&x=-1&y=-1](http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93?w=100&h=100&x=-1&y=-1)

[http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93?w=100&h=100&x=-1&y=-1&r=45](http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93?w=100&h=100&x=-1&y=-1&r=45)

[http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93?w=100&h=100&x=-1&y=-1&g=1](http://182.92.189.64:8081/a258607b53444f32208e864f44a06b93?w=100&h=100&x=-1&y=-1&g=1)

	
	
	

	
	



