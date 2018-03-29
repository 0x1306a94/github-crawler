# github-crawler 

* 爬取`Github` `Trending`数据
- 如何运行
    * 首先下载本仓库,如果有安装`Go` 环境,则存放 到`GOPATH/src`目录下
    * 源码运行
    
    ```shell
    git clone https://github.com/king129/github-crawler.git $GOPATH/src/github-crawler
    cd $GOPATH/src/github-crawler
    go run app/main.go
    ```
    
    * 可执行文件运行
    
    ```shell
    cd xxx/github-crawler/app
    chmod +x app
    ./app
    ```
    
    * `Docker`运行

    ```shell
    cd xxx/github-crawler
    docker build -t github-crawler:1.0 .
    docker run github-crawler:1.0
    ```


