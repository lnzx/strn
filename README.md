# strn

部署在vercel的saturn节点查询工具

## 环境变量
`DATABASE_URL`

postgres://postgres:PASSWORD@IP:5432/postgres?sslmode=disable

## cron
定时抓取数据分析,15分钟请求一次
```shell
https://strn.app/api?cron=1
```

## api
查询官网节点
```shell
https://strn.app/api
```

## 可选
api 接口
```
VITE_API_BASE_URL
```



