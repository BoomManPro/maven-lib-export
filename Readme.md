# maven lib export


因为有些公司需要离线开发,但是相关jar包无法迁入内网环境,遂开发此小工具在外网执行使用。



## 设计思路

1. `mvn dependency:copy-dependencies -DoutputDirectory=lib`
将项目的jar第三方依赖导出到目录 `lib`

2. `mvn help:evaluate -Dexpression=settings.localRepository`
使用命令获取maven repository的地址

然后`递归` 搜索`lib`的涉及的文件夹,并且包含parent文件夹 因为spring-boot-parent的pom需要使用


## build windows

```
CGO_ENABLED=0;GOOS=windows;GOARCH=amd64
```


## Use

download update workspace to run