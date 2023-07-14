## jaeger
官网地址：https://www.jaegertracing.io/docs/1.47/getting-started/

**jaeger client已弃用。**

**jaeger agent已弃用。**

### 端口
| 组件        | 端口    | 作用     |
|-----------|-------|--------|
| query     | 16686 | ui界面查询 |
| collector | 14268 | 直接上报数据 |
| ingester  | 14270 | 数据接收器  |

### collector和ingester
collector可以直接上报存储。也可以上报给ingester，通过ingester上报存储。
数据少的时候可以使用，数据多需要通过ingester，可以异步处理提高吞吐量。

### docker部署
```shell
# 使用es作存储
docker run -d --name jaeger \
--network es \
-e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
-e COLLECTOR_OTLP_ENABLED=true \
-e PROMETHEUS_QUERY_SUPPORT_SPANMETRICS_CONNECTOR=true \
-e SPAN_STORAGE_TYPE=elasticsearch \
-e ES_SERVER_URLS=http://es:9200 \
-e ES_NUM_REPLICAS=0 \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778 \
-p 16686:16686 \
-p 4317:4317 \
-p 4318:4318 \
-p 14250:14250 \
-p 14268:14268 \
-p 14269:14269 \
-p 9411:9411 \
jaegertracing/all-in-one:1.46
  ```
http://localhost:16686 jaeger ui




