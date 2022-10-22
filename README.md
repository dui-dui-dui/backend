# backend

-----

## 配置示例 v2

```json
{
  "groups": [
    {
      "group_id": "pd",
      "group_index": 0,
      "group_override": false,
      "rules": [
        {
          "group_id": "pd",
          "id": "default",
          "index": 0,
          "override": false,
          "start_key": "",
          "end_key": "",
          "role": "voter",
          "count": 3,
          "location_labels": [
            "zone"
          ]
        },
        {
          "group_id": "pd",
          "id": "meta",
          "index": 1,
          "override": false,
          "start_key": "",
          "end_key": "7480000000000000FF0000000000000000F8",
          "role": "voter",
          "count": 2,
          "location_labels": [
            "zone"
          ]
        }
      ]
    },
    {
      "group_id": "tiflash",
      "group_index": 120,
      "group_override": false,
      "rules": [
        {
          "group_id": "tiflash",
          "id": "happy",
          "index": 0,
          "override": false,
          "start_key": "7480000000000000FF4E00000000000000F8",
          "end_key": "7480000000000000FF5000000000000000F8",
          "role": "learner",
          "count": 1
        }
      ]
    }
  ],
  "labels": [
    {
      "key": "engine",
      "values": [
        "tikv",
        "tiflash"
      ]
    },
    {
      "key": "zone",
      "values": [
        "zone-1",
        "zone-2",
        "zone-3"
      ]
    },
    {
      "key": "disk",
      "values": [
        "ssd",
        "hdd"
      ]
    }
  ],
  "markdown": "group pd\n\n01/01/2022 - 01/08/2022: default\n01/01/2022 - 01/01/2022: meta\ngroup tiflash\n\n01/06/2022 - 01/06/2022: happy\n",
  "schemas": [
    {
      "ts": 1640966400000000,
      "size": 0,
      "left": 0,
      "name": "meta",
      "description": "meta data of tidb cluster",
      "start_key": "",
      "end_key": "7480000000000000FF0000000000000000F8"
    },
    {
      "ts": 1641052800000000,
      "size": 0,
      "left": 0,
      "name": "system",
      "description": "system tables of mysql database",
      "start_key": "7480000000000000FF4400000000000000F8",
      "end_key": "7480000000000000FF4600000000000000F8"
    },
    {
      "ts": 1641139200000000,
      "size": 0,
      "left": 0,
      "name": "foo",
      "description": "test/foo",
      "start_key": "7480000000000000FF4600000000000000F8",
      "end_key": "7480000000000000FF4800000000000000F8"
    },
    {
      "ts": 1641225600000000,
      "size": 0,
      "left": 0,
      "name": "bar",
      "description": "test/bar",
      "start_key": "7480000000000000FF4800000000000000F8",
      "end_key": "7480000000000000FF4A00000000000000F8"
    },
    {
      "ts": 1641312000000000,
      "size": 0,
      "left": 0,
      "name": "baz",
      "description": "test/baz",
      "start_key": "7480000000000000FF4A00000000000000F8",
      "end_key": "7480000000000000FF4C00000000000000F8"
    },
    {
      "ts": 1641398400000000,
      "size": 0,
      "left": 0,
      "name": "happy",
      "description": "hackathon/happy",
      "start_key": "7480000000000000FF4E00000000000000F8",
      "end_key": "7480000000000000FF5000000000000000F8"
    },
    {
      "ts": 1641484800000000,
      "size": 0,
      "left": 0,
      "name": "lucky",
      "description": "hackathon/lucky",
      "start_key": "7480000000000000FF5000000000000000F8",
      "end_key": "7480000000000000FF5200000000000000F8"
    },
    {
      "ts": 1641571200000000,
      "size": 0,
      "left": 0,
      "name": "default",
      "description": "future tables",
      "start_key": "7480000000000000FF5200000000000000F8",
      "end_key": ""
    }
  ]
}
```

## 配置示例

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

```json
{
  "schemas": [
    {
      "name": "meta",
      "description": "meta data of tidb cluster",
      "start_key": "",
      "end_key": "7480000000000000FF0000000000000000F8"
    },
    {
      "name": "system",
      "description": "system tables of mysql database",
      "start_key": "7480000000000000FF4400000000000000F8",
      "end_key": "7480000000000000FF4600000000000000F8"
    },
    {
      "name": "foo",
      "description": "test/foo",
      "start_key": "7480000000000000FF4600000000000000F8",
      "end_key": "7480000000000000FF4800000000000000F8"
    },
    {
      "name": "bar",
      "description": "test/bar",
      "start_key": "7480000000000000FF4800000000000000F8",
      "end_key": "7480000000000000FF4A00000000000000F8"
    },
    {
      "name": "baz",
      "description": "test/baz",
      "start_key": "7480000000000000FF4A00000000000000F8",
      "end_key": "7480000000000000FF4C00000000000000F8"
    },
    {
      "name": "happy",
      "description": "hackathon/happy",
      "start_key": "7480000000000000FF4E00000000000000F8",
      "end_key": "7480000000000000FF5000000000000000F8"
    },
    {
      "name": "lucky",
      "description": "hackathon/lucky",
      "start_key": "7480000000000000FF5000000000000000F8",
      "end_key": "7480000000000000FF5200000000000000F8"
    },
    {
      "name": "default",
      "description": "future tables",
      "start_key": "7480000000000000FF5200000000000000F8",
      "end_key": ""
    }
  ],
  "rule_config": [
    {
      "group_id": "pd",
      "group_index": 0,
      "group_override": false,
      "rules": [
        {
          "group_id": "pd",
          "id": "default",
          "index": 0,
          "override": false,
          "start_schema": "meta",
          "end_schema": "default",
          "start_key": "",
          "end_key": "",
          "role": "voter",
          "count": 3
        },
        {
          "group_id": "pd",
          "id": "meta",
          "index": 1,
          "override": false,
          "start_schema": "meta",
          "end_schema": "meta",
          "start_key": "",
          "end_key": "7480000000000000FF0000000000000000F8",
          "role": "voter",
          "count": 2
        }
      ]
    },
    {
      "group_id": "tiflash",
      "group_index": 120,
      "group_override": false,
      "rules": [
        {
          "group_id": "tiflash",
          "id": "happy",
          "index": 0,
          "override": false,
          "start_schema": "happy",
          "end_schema": "happy",
          "start_key": "7480000000000000FF4E00000000000000F8",
          "end_key": "7480000000000000FF5000000000000000F8",
          "role": "learner",
          "count": 1
        }
      ]
    }
  ],
  "store_labels": [
    {
      "key": "engine",
      "values": [
        "tikv",
        "tiflash"
      ]
    },
    {
      "key": "zone",
      "values": [
        "zone-1",
        "zone-2",
        "zone-3"
      ]
    },
    {
      "key": "disk",
      "values": [
        "ssd",
        "hdd"
      ]
    }
  ]
}
```