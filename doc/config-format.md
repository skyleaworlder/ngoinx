# config

## metadata

* 文件类型：`json`；

## structure

```json
{
  "service": [
    {
      "listen": uint16,
      "proxy": [
        {
          "src": string,
          "target": [
            { "dst": string, "weight": uint8 }
          ],
        }
      ]
    }
  ]
}
```
