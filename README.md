# backend

-----

## 下载运行

下载服务端 https://github.com/dui-dui-dui/backend/releases

运行：

```bash
./backend --tidb 127.0.0.1:10080 --pd 127.0.0.1:2379 --addr :8080
```

开发模式：

```bash
./backend --dev
```

## 接口

### GET /config

获取配置，示例 https://github.com/dui-dui-dui/backend/blob/main/data/config_out.json

**开发模式下，schema/store/store_labels 都是固定的默认值，groups 读取 groups.json 文件，没有此文件时也返回固定默认值**

说明：

* schemas 为 schema 信息（页面顶部显示），数组的顺序即为显示顺序。其中有几项比较特殊，可以换不同的颜色
   * meta 是 tidb 元信息
   * system 是系统表
   * default 放在最右边，是后续创建新表的默认配置
   * 其余的是正常的 user table

* rule_config 是 placement rule 配置，是一个数组，分为若干个 group，每个 group 有一些属性，和若干个 rule
   * group 的属性有用的是名字（用于显示）index（决定竖直方向的顺序）override（影响rule的覆盖关系）
   * rule 的属性有用的是 key（用于显示）index（决定竖直方向的顺序）override (影响rule的覆盖关系)
   * rule 的长度和位置由 start_schema 和 end_schema 决定，比如 start_schema=meta end_schema=system，就是覆盖 meta 和 system 两个 schema

* store_labels 是所有可选的 label，用于配置 rule 的表单使用

### GET /region

获取 region 分布，格式为 `{"markdown":"$MARKDOWN"}`，MARKDOWN 文本示例： https://github.com/dui-dui-dui/backend/blob/main/data/region_out.txt

**开发模式下，每次取到的都是随机分布**

### POST /save

更新配置，格式与 GET /config 一致 https://github.com/dui-dui-dui/backend/blob/main/data/config_out.json

**开发模式下，保存在本地文件 groups.json，如果文件出问题了可以手动删掉使用默认值重新来**

