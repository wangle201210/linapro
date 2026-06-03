## Context

C6 交付对外传播的 H5 数字纪念墙(见程序级设计 C6)。复用 C3 首发激活(`activation.is_first`)、C1 玩家昵称/身份、C2 卡片/金句。核心是**公开免登录 + 隐私最小暴露**。

约束:遵循 `.agents/rules/` 的 `plugin`、`backend-go`、`api-contract`、`data-permission`、`architecture`、`frontend-ui`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 公开只读纪念墙 API(首发墙/校史精选/统计,免登录);H5 静态页(首发墙展示 + 扫码回跳小程序);只暴露昵称与活动信息。

**Non-Goals:** 看板/导出/结算(C7);荣誉持久发放(C7);小程序 UI;新增数据表。

## Decisions

### D1 公开免登录路由 + 隐私最小暴露
纪念墙 API 绑定在 `plugin.go` 的**公开子路由组**(与玩家登录同级:基础中间件 + `HandlerResponse`,**无 Auth、无 player-auth**)。响应**只含昵称、身份标签、牛名称/序号、到场顺序与公开统计**;**绝不返回手机号、openid、设备指纹**。投影在 DTO 与 service 层显式裁剪。

### D2 首发墙派生自激活记录
首发墙 = `activation` 中 `is_first=1` 的记录,关联 `user`(昵称/身份)与 `niu`(名称/序号),按到场顺序/牛序号排序,有界返回(≤120)。批量装配,杜绝 N+1。

### D3 校史精选与统计
- 校史精选:从已上线 `card`(标题/分类/图片)与 `quote`(文案)取样(随机/前 N),有上限。
- 统计:已激活牛数(`niu.status=active` 计数 或 distinct 首发数)、参与玩家数(`user` 计数或激活玩家数)等公开聚合,DB 侧计数。

### D4 H5 静态页 + public_assets
新增 `frontend/wall/index.html`(纯 HTML+内联 JS 的移动端页面),`plugin.yaml` 声明:
```
public_assets:
  - source: frontend/wall
    mount: /wall
    index: index.html
```
H5 页 `fetch` 公开 API(`/x/sicau-niu/api/v1/plugins/sicau-niu/wall/*`)渲染首发墙/精选/统计,并提供**扫码/链接回跳小程序**入口(小程序跳转 URL/scheme 由配置或占位)。该静态页独立于 `frontend/pages`(admin Vue),不进入 vben 构建。

### D5 数据权限(公开活动展示)
纪念墙为公开读,例外说明:展示昵称与活动信息属公开传播预期,**已在 DTO/service 裁剪隐私字段**;有界返回防滥用。无写操作。

### D6 显式依赖注入
wall service 显式注入 DAO(activation/user/niu/card/quote 读取)或复用既有读取能力;controller 持 service 字段;不在请求路径 `New()`。

### API 契约(C6,公开免登录,时间 Unix 毫秒)
| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/plugins/sicau-niu/wall/first-activators` | 首发玩家名单(昵称/身份/牛名/序号/到场顺序),有界 |
| `GET` | `/plugins/sicau-niu/wall/highlights` | 校史精选(卡片/金句取样) |
| `GET` | `/plugins/sicau-niu/wall/stats` | 活动公开统计(已激活牛数/参与玩家数等) |

（均在公开子路由组,无 permission 标签、无 player-auth。）

## Risks / Trade-offs

- [隐私暴露] → DTO 仅含昵称/身份/活动字段,service 显式裁剪;单测断言响应不含手机号/openid。
- [公开接口滥用] → 有界返回(首发 ≤120、精选/统计有上限);只读无副作用。
- [小程序回跳形态] → 扫码/链接的小程序跳转 URL/scheme 用配置或占位,真实形态后续替换。

## Migration Plan

1. `api/wall` DTO → `make ctrl` → 填 controller/service;`plugin.go` 公开子组绑定。
2. `frontend/wall/index.html`(H5 页);`plugin.yaml` 声明 `public_assets`。
3. 服务层 DB 门控单测(首发派生/隐私裁剪/统计)+ 编译门禁 + 公开 API 冒烟(无 token 可访问);`openspec validate --strict`。
4. 回滚:代码/静态页/manifest 可整体回退;无新增表、无生产数据。

## Open Questions

1. 小程序回跳的 URL Scheme / 二维码内容(占位,后续替换)。
2. 校史精选的取样口径(随机 vs 最新,数量上限默认 N)。

> 影响判断:i18n 单语言不启用,**无影响**;字典 —— 无新增枚举,**无影响**;缓存一致性 —— 无新增缓存,公开 API 为 DB 读取,**无影响**;数据权限 —— 公开活动展示,**已裁剪隐私字段并有界**(例外已记录);跨平台 —— 仅插件内文件,**无影响**。
