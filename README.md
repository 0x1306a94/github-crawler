# github-crawler 

* 爬取`Github` `Trending`数据,数据存储到`Redis`,在运行之前,请先设置 环境变量 `RedisHost` `RedisPassword`,如果不进行设置,则默认连接到 `127.0.0.1:6379`,`Redis`数据过期时间为1小时
- 如何运行
    * 首先下载本仓库,如果有安装`Go` 环境,则存放 到`GOPATH/src`目录下
    * 源码运行
    
    ```shell
    git clone https://github.com/king129/github-crawler.git $GOPATH/src/github-crawler
    cd $GOPATH/src/github-crawler
    go run app/main.go
    ```
    
    * 可执行文件运行 (注意: 可执行文件只能在 macOS 平台运行)
    
    ```shell
    // 每五分钟分钟执行60次任务 任务为 每个语言的 repo/developer daily weekly monthly 的数据
    cd xxx/github-crawler/app
    chmod +x app
    ./app
    
    cd xxx/github-crawler/server
    chmod +x server
    ./server
    ```
    
    * `Docker`运行

    ```shell
    cd xxx/github-crawler
    docker build -t github-crawler:1.0 .
    docker run github-crawler:1.0
    ```

    ![](https://github.com/king129/github-crawler/blob/master/images/1.png)
    ![](https://github.com/king129/github-crawler/blob/master/images/2.png)
    ![](https://github.com/king129/github-crawler/blob/master/images/3.png)
    ![](https://github.com/king129/github-crawler/blob/master/images/4.png)
    ![](https://github.com/king129/github-crawler/blob/master/images/5.png)

