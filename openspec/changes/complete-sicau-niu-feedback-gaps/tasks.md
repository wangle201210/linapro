## Feedback

- [x] **FB-1**: 第 2 点，激活海报和电子证书需要输出最终可展示图片，而不是基础占位 PNG
- [x] **FB-2**: 第 3 点，基础互动规则和排行榜统计规则需要运营端可查看、可维护并参与运行时生效
- [x] **FB-3**: 第 4 点，运营管理端需要独立的排行榜数据查看入口
- [x] **FB-4**: 运营结算页面视觉层级不足，规则、榜单、风控和动作区块平铺后不够清晰美观
- [x] **FB-5**: 运营动作区“结算归档”说明文案被输入框和按钮挤压成逐字竖排

### 根因记录

- `FB-1`根因：首版渲染器只验证`PNG`可解码和不同玩家输出有差异，没有把海报/证书作为最终视觉产物处理，导致输出缺少可读活动信息和个人信息。
- `FB-2`根因：签到、偷草、送草、排行榜和风控阈值通过插件启动配置读取，配置服务本身为只读能力，运营端无法在管理工作台查看或修改，也无法在请求路径实时生效。
- `FB-3`根因：排行榜服务已提供玩家端喂草榜、院系榜和好友榜，但路由和前端仅面向玩家端，运营结算页没有管理端表格投影。
- `FB-4`根因：当前结算页把看板、活跃、规则、排行榜、动作、风控和归档全部按相同`Card`平铺，缺少首屏概览、模块主次和表单分组；规则字段数量增加后视觉密度失衡，排行榜也缺少运营视角的标题和摘要。
- `FB-5`根因：运营动作卡片位于右侧窄栏，但`.action-item`仍使用“两列网格 + 右侧控件自动宽度”布局；“结算归档”的输入框和按钮占据过多横向空间后，左侧说明文案剩余宽度接近 0，导致中文逐字换行。

### 影响分析

- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`。
- 插件本地规范：`apps/lina-plugins/sicau-niu/AGENTS.md`不存在，按仓库顶层规范和规则文件执行。
- `i18n`影响：`sicau-niu`插件未在`plugin.yaml`启用`i18n`，新增单语言中文页面文案、菜单按钮和接口文档源文本不要求维护插件`manifest/i18n`资源。
- 缓存一致性影响：本次不新增缓存；规则配置以数据库表为权威数据源，更新后无缓存失效路径。
- 数据权限影响：新增运营端读取和更新接口受宿主认证、租户和权限中间件保护；玩家端数据仍只经玩家 token 访问。新增排行榜运营投影为活动运营全局统计，按平台级插件能力处理。
- 开发工具跨平台影响：不修改脚本、`Makefile`、`make.cmd`、`CI`或`linactl`；数据库生成使用既有插件`make dao`入口。
- DI 来源检查：新增规则服务由`backend/plugin.go`在路由装配期创建一次，并显式传入依赖的业务服务和控制器，共享同一数据库后端，不在请求路径临时构建服务图。

### 执行记录

- `FB-1`修复：升级`posterrender`和`certrender`默认实现，输出带活动标题、玩家/奖项字段、编号、纪念标识、版式边框和个性化装饰的`PNG`；新增`rendertext`工具，优先使用`SICAU_NIU_RENDER_FONT`或系统`CJK`字体，缺失字体时回退基础字体并保持可解码输出；补充海报和证书渲染器像素级单测，防止退回纯占位图。
- `FB-2`修复：新增`plugin_sicau_niu_rule_config`迁移、卸载 SQL、DAO/DO/Entity 生成入口和`rules`服务；运营端新增规则查询/更新接口，关键玩家和运营路径读取同一个运行规则服务，覆盖`LBS`判距、海报/证书标识、签到草量、偷草/送草限制、铁牛加成阈值、排行榜`TopN`、异常阈值和小程序回跳地址。
- `FB-3`修复：运营端新增喂草榜、好友榜、院系榜接口和结算页“排行榜数据”表格区；前端客户端、页面对象和`TC004`可见性断言同步更新。
- `FB-4`修复：将运营结算页从同权重`Card`平铺重排为运营工作台结构；新增顶部“寻牛活动运营概览”、激活进度、风险状态标签和关键指标区；活跃、规则、关键互动、运营动作、排行榜、风控和归档按运营任务分区展示；规则表单按“激活与展示”“签到与送草”“偷草策略”“榜单与风控”分组，排行榜增加榜首、记录数和领先院系摘要；保留既有接口、权限码和`data-testid`测试锚点。
- `FB-5`修复：运营动作区改为纵向动作项布局，不再让说明文案和右侧输入控件竞争同一行宽；批量发证与结算归档控件在动作项内部按可用宽度排列，归档输入框使用`minmax`自适应，说明文案保持正常横排可读；同步补充`TC004`对归档说明文本宽度的布局断言。

### 验证记录

- 已执行并通过：`make db.init confirm=init`。
- 已执行并通过：`psql postgresql://postgres:postgres@127.0.0.1:5432/linapro?sslmode=disable -v ON_ERROR_STOP=1 -f manifest/sql/007-sicau-niu-rule-config.sql`。
- 已执行并通过：`make dao`（在`apps/lina-plugins/sicau-niu`内生成`rule_config`相关`DAO/DO/Entity`）。
- 已执行并通过：`GOWORK=off go test ./... -count=1`（在`apps/lina-plugins/sicau-niu`内）。
- 已新增：`rules`服务数据库门控测试；普通环境跳过，设置`LINA_TEST_PGSQL_LINK`时覆盖规则配置写入、读取和窄访问器生效路径。
- 已执行并通过：`pnpm -C apps/lina-vben/apps/web-antd typecheck`。
- 已执行并通过：`pnpm -C hack/tests test:validate`。
- 已执行并通过（`FB-4`）：`pnpm -C apps/lina-vben/apps/web-antd typecheck`。
- 已执行并通过（`FB-4`）：`pnpm -C hack/tests test:validate`。
- 已执行并通过（`FB-4`）：`pnpm -C hack/tests test:module -- plugin:sicau-niu -- --grep TC-4a`。
- 已执行并通过（`FB-4`）：`openspec validate complete-sicau-niu-feedback-gaps --strict`。
- 已执行并通过（`FB-5`）：`pnpm -C apps/lina-vben/apps/web-antd typecheck`。
- 已执行并通过（`FB-5`）：`pnpm -C hack/tests test:validate`。
- 已执行并通过（`FB-5`）：`pnpm -C hack/tests test:module -- plugin:sicau-niu -- --grep TC-4a`。
- 已执行并通过（`FB-5`）：`openspec validate complete-sicau-niu-feedback-gaps --strict`。
