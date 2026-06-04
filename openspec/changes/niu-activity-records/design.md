## Context

补齐运营后台的活动行为记录只读查询(喂草/偷草/送草/签到/激活/草账户流水)。复用 C3/C4 行为表与 C1 玩家、C2 牛。核心是**统一的只读分页 + 批量装配杜绝 N+1 + 管理端鉴权**。约束遵循 `.agents/rules/` 的 `backend-go`、`api-contract`、`data-permission`、`architecture`、`frontend-ui`、`testing`;单语言中文不启用 i18n;`platform_only`。

## Goals / Non-Goals

**Goals:** 6 类记录的只读分页查询(可按玩家/牛筛选),DB 侧分页 + 批量昵称/牛名装配,管理端鉴权。

**Non-Goals:** 记录增删改;记录导出;玩家端 UI;新增表。

## Decisions

### D1 统一只读分页范式
每类记录一个列表方法:`model := dao.X.Ctx(ctx)` → 按筛选 `Where/WhereLike` → `Count()` → `Fields(...).Page(pageNum,pageSize).OrderDesc(created_at).Scan(&rows)`。分页默认 1/10、上限 100(复用现有 admin 列表口径)。绝不全量载入内存。

### D2 批量装配杜绝 N+1
列表页内涉及的玩家 ID、牛 ID 先收集去重,再各一条 `WhereIn` 投影查询装配 `id→昵称`、`id→牛名/编码` 映射,内存合并。喂草/激活含牛;偷草/送草含双方玩家;签到/草流水含玩家。**禁止逐行查名**。

### D3 显式依赖注入与分层
新增 `record` service(`New()` 无依赖,读自有表)、`recordController`(持 service 字段,装配期构建)、`api/record` 契约。`plugin.go` 管理端子组绑定 6 个接口。

### D4 数据权限与权限模型
全部为管理端只读查询,经宿主 `Auth+Tenancy+Permission`;6 个 DTO 统一 `permission:"sicau-niu:record:list"`(同一"记录查看"能力)。只读、有界、DB 分页。

### API 契约(管理端鉴权,时间 Unix 毫秒)
| 方法 | 路径 | 权限 | 筛选 |
|------|------|------|------|
| `GET` | `/plugins/sicau-niu/admin/records/feedings` | `sicau-niu:record:list` | niuId, userId |
| `GET` | `/plugins/sicau-niu/admin/records/steals` | `sicau-niu:record:list` | actorUserId, targetUserId |
| `GET` | `/plugins/sicau-niu/admin/records/gifts` | `sicau-niu:record:list` | fromUserId, toUserId |
| `GET` | `/plugins/sicau-niu/admin/records/checkins` | `sicau-niu:record:list` | userId |
| `GET` | `/plugins/sicau-niu/admin/records/activations` | `sicau-niu:record:list` | userId, niuId |
| `GET` | `/plugins/sicau-niu/admin/records/grass-txns` | `sicau-niu:record:list` | userId |

各响应 `{list, total}`;行含玩家昵称/牛名等装配字段与该记录业务字段。

### D5 前端与菜单
后台新增"活动记录"子目录(D,parent=寻牛活动 root)+ 6 个 vxe 查询页(M),筛选项(玩家 ID/牛 ID 等)+ 分页;统一权限 `sicau-niu:record:list`。

## Risks / Trade-offs

- [记录数据量大] → DB 侧分页(默认 10、上限 100)+ 批量装配,杜绝 N+1;按 created_at 倒序有界返回。
- [筛选为 ID 而非名称] → 运营按玩家/牛 ID 精确筛选(从其它页可得 ID);名称在结果中装配展示。模糊名搜索后续可加。

## Migration Plan

1. `api/record` DTO/接口 → `internal/service/record` 6 方法 + 批量装配 → `internal/controller/record` → `plugin.go` 绑定。
2. 前端"活动记录"子目录 + 6 页 + record-client;`plugin.yaml` 菜单/权限。
3. 服务层 DB 门控单测(feeding/steal/grass_txn 分页+筛选+装配)+ 编译门禁 + `openspec validate --strict`。
4. 回滚:纯新增只读路径 + 菜单/权限,可整体回退。

## Open Questions

1. 是否需要按玩家昵称/牛名称的模糊搜索(本期按 ID 筛选)。

> 影响判断:字典 —— 无新增枚举,**无影响**;缓存一致性 —— 实时读取,**无新增缓存**;数据权限 —— 管理端鉴权 + `permission` 标签,只读有界(已记录);跨平台 —— 仅插件内文件,**无影响**;i18n —— 单语言,**无影响**。
