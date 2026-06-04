## Context

C2 建立寻牛活动的**内容资产层**(见 `apps/lina-plugins/sicau-niu/design.md` 程序级设计的 C2)。C1 已交付玩家身份与院系字典;C2 在其上新增牛、铁牛、卡片、校史金句的数据模型与运营后台录入,作为 C3(地图/激活/发卡)、C4(喂草/铁牛加成)的数据来源。

约束:遵循 `.agents/rules/` 的 `plugin`、`backend-go`、`api-contract`、`database`、`architecture`、`data-permission`、`frontend-ui`、`testing`;单语言(中文),不启用 i18n;插件作用域沿用 C1 的 `platform_only`/`global`。

## Goals / Non-Goals

**Goals:**
- 运营可在「寻牛活动」顶级菜单下管理:牛(含放出节奏)、铁牛标识、卡片(1 牛 1 卡)、校史金句。
- 牛关联院系(复用 C1 院系字典)与主卡;卡片图片通过宿主文件管理上传并以路径引用。
- 提供面向运营的受权限 API(不实现小程序 UI;面向小程序的牛/卡片读取接口在 C3 提供)。

**Non-Goals:**
- 数据看板/统计(C7);牛状态机/激活、地图与可见牛列表(C3);铁牛实时定位与喂草加成(C4);喂草/签到/偷草/送草(C4)。

## Decisions

### D1 内容资产为插件自有表 + 运营 CRUD
新增 4 张插件自有表;运营后台受宿主 `Auth+Tenancy+Permission` 守卫,权限 `sicau-niu:niu:*`、`sicau-niu:iron:*`、`sicau-niu:card:*`、`sicau-niu:quote:*`。

### D2 「1 牛 1 卡」用 `card.niu_id` 唯一约束承载
卡片归属一头牛,`plugin_sicau_niu_card.niu_id` 在未删除集合内唯一。运营流程:建牛 → 建卡并选所属牛 → 该牛即拥有唯一主卡。C3 激活时按 `niu_id` 取主卡发放。
- **备选**:`niu.card_id` 反向引用 —— 否决,卡片是内容主体、按牛归属更自然,且便于「先建牛后补卡」。

### D3 放出节奏建模为牛的字段
`plugin_sicau_niu_niu` 存:`release_stage`(预热/主体/高潮/收尾)、`online_at`(上线时间点,API Unix 毫秒)、可选 `visible_weekdays`(如 `1,3,5`)、`visible_start`/`visible_end`(`HH:MM` 时段)。C2 只存配置,**可见性计算在 C3**。

### D4 稳定枚举用 Go 命名常量(不入宿主字典)
牛类型(普通/特殊)、特殊子类(学院/贡献/校友/精神)、卡片分类(人物/事件/科研/院系/精神)、活动阶段(预热/主体/高潮/收尾)为固定、单语言、运营不改文案的稳定枚举,用 Go 命名类型+常量管理,满足 `backend-go` 枚举规则;不引入宿主字典耦合,与 C1 身份类型策略一致。**记录无字典影响**。

### D5 铁牛标识登记,实时经纬留给 C4
`plugin_sicau_niu_iron` 存:`code`(铁牛/设备标识,唯一)、`name`、`remark`,以及 `last_lat`/`last_lng`/`located_at`(可空,**由 C4 写入**,C2 运营不设置)。C2 只 CRUD 铁牛标识。

### D6 卡片图片用宿主文件管理上传 + 路径引用
卡片图片不由插件自管存储;运营在卡片表单内通过**宿主文件上传**(`system:file`,落 `/uploads`)得到存储路径,卡片表存 `image_path` 字符串,展示走宿主 `/uploads/{path}`。插件无需自有文件存储与清理。

### D7 牛关联院系复用 C1
学院牛的 `college_id` 引用 C1 `plugin_sicau_niu_college`;运营在牛表单从院系下拉(复用 C1 院系数据)选择;删除院系的引用保护已在 C1,C2 不重复。

### D8 列表性能:DB 侧分页 + 批量装配
牛列表需展示关联院系名、是否已绑卡 → 先 DB 侧过滤/排序/分页取当前页,再**批量**收集 `college_id`/`niu_id` 一次性查院系名与卡片归属,内存合并,杜绝 N+1。卡片/金句/铁牛列表同样 DB 侧分页。所有列表有 `maxPageSize` 上限。

### D9 显式依赖注入
各 service 构造函数显式注入 DAO 与跨能力依赖(如牛 service 需院系存在性校验能力);controller 在 `NewV1()` 持有 service 字段;不在请求路径 `New()` 服务。

### 数据模型(C2)
| 表 | 关键列 | 约束/索引 |
|----|-------|-----------|
| `plugin_sicau_niu_niu` | `code`(序号)、`niu_type`、`special_subtype`、`name`、`college_id`、`lat`、`lng`、`release_stage`、`online_at`、`visible_weekdays`、`visible_start`、`visible_end`、`status`(默认 inactive)、软删/时间 | `code` 活跃唯一;`release_stage`/`online_at` 索引(供 C3 可见性查询);`college_id` 索引 |
| `plugin_sicau_niu_iron` | `code`、`name`、`last_lat`、`last_lng`、`located_at`、`remark`、软删/时间 | `code` 活跃唯一 |
| `plugin_sicau_niu_card` | `niu_id`、`category`、`title`、`content`、`image_path`、软删/时间 | `niu_id` 活跃唯一(1 牛 1 卡);`category` 索引 |
| `plugin_sicau_niu_quote` | `content`、`enabled`、软删/时间 | (可选 `enabled` 索引) |

> 软删除/时间走 `deleted_at` 自动维护;DAO/DO/Entity 由 `make dao` 生成。

### API 契约(C2,运营面,`Auth+Permission`,时间 Unix 毫秒)
| 资源 | 方法/路径 | 权限 |
|------|-----------|------|
| 牛 | `GET/POST /admin/niu`、`GET/PUT/DELETE /admin/niu/{id}` | `sicau-niu:niu:list/create/update/delete` |
| 铁牛 | `GET/POST /admin/iron`、`PUT/DELETE /admin/iron/{id}` | `sicau-niu:iron:*` |
| 卡片 | `GET/POST /admin/cards`、`GET/PUT/DELETE /admin/cards/{id}` | `sicau-niu:card:*` |
| 金句 | `GET/POST /admin/quotes`、`PUT/DELETE /admin/quotes/{id}` | `sicau-niu:quote:*` |

菜单挂在 C1 已建的 `plugin:sicau-niu:root`(寻牛活动)顶级目录下。

## Risks / Trade-offs

- [牛字段较多、放出节奏可选项多] → C2 只存配置,可见性/状态机判定推迟到 C3,避免在 C2 实现复杂时序逻辑。
- [1 牛 1 卡的删除一致性] → 删牛时其主卡的处理(级联删卡 or 阻止删牛):本期采用「删牛前要求先解绑/删卡,否则拒绝」或「删牛同时软删其主卡」,在 specs 明确(默认:删牛软删级联其主卡,可恢复)。
- [卡片图片依赖宿主文件管理] → 复用稳定的宿主 `system:file` 能力,不新增插件存储;若宿主文件模块禁用,卡片图片降级为可空。
- [枚举用 Go 常量而非宿主字典] → 单语言、运营不改文案,符合 backend-go;若后续需运营自定义分类再评估字典化。

## Migration Plan

1. 新增当前迭代 `manifest/sql/002-sicau-niu-catalog.sql`(4 表 + 索引)与 `manifest/sql/uninstall/`;`make dao` 生成 DAO/DO/Entity。
2. `plugin.yaml` 在「寻牛活动」下新增 牛/卡片/金句/铁牛 管理菜单与权限按钮。
3. 新增 `api/`(各资源 DTO)→ `make ctrl` 生成骨架 → 填充 controller/service。
4. 新增运营后台前端页面(各资源 vxe-grid + 表单弹窗;卡片表单集成宿主文件上传)。
5. 编译门禁(plugins workspace 构建)+ 服务层单测 + 运营页 E2E;`openspec validate --strict`。
6. 回滚:plugin.yaml 与代码可整体回退;插件未上线、无生产数据。

## Open Questions

1. 删牛时主卡处理策略最终口径(级联软删 vs 阻止),specs 定。
2. 牛「序号 code」是否由运营手填还是系统按上线顺序生成(默认运营手填,唯一)。
3. 放出节奏的「周几/时段」是否本期就要做,还是仅留字段、C3 再用(默认:C2 仅存字段)。

> 影响判断:i18n —— 单语言不启用,**无影响**;字典 —— 稳定枚举用 Go 常量,**无宿主字典影响**(已记录);缓存一致性 —— C2 为运营 CRUD,无新增跨节点缓存,**无影响**;数据权限 —— 内容资产为运营受权限的全量数据,无玩家自隔离维度;开发工具跨平台 —— 仅插件内 SQL/代码,无脚本变更。
