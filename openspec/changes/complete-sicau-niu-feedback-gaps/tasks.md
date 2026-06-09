## Feedback

- [x] **FB-1**: 第 2 点，激活海报和电子证书需要输出最终可展示图片，而不是基础占位 PNG
- [x] **FB-2**: 第 3 点，基础互动规则和排行榜统计规则需要运营端可查看、可维护并参与运行时生效
- [x] **FB-3**: 第 4 点，运营管理端需要独立的排行榜数据查看入口
- [x] **FB-4**: 运营结算页面视觉层级不足，规则、榜单、风控和动作区块平铺后不够清晰美观
- [x] **FB-5**: 运营动作区“结算归档”说明文案被输入框和按钮挤压成逐字竖排
- [x] **FB-6**: 运营动作区“批量发证”只能手输证书荣誉 ID，缺少可执行的证书候选选择，导致运营人员无法完成发证动作
- [x] **FB-7**: 牛只发布状态同时依赖`release_stage`和`online_at`，导致运营维护重复且玩家可见/激活逻辑不一致

### 根因记录

- `FB-1`根因：首版渲染器只验证`PNG`可解码和不同玩家输出有差异，没有把海报/证书作为最终视觉产物处理，导致输出缺少可读活动信息和个人信息。
- `FB-2`根因：签到、偷草、送草、排行榜和风控阈值通过插件启动配置读取，配置服务本身为只读能力，运营端无法在管理工作台查看或修改，也无法在请求路径实时生效。
- `FB-3`根因：排行榜服务已提供玩家端喂草榜、院系榜和好友榜，但路由和前端仅面向玩家端，运营结算页没有管理端表格投影。
- `FB-4`根因：当前结算页把看板、活跃、规则、排行榜、动作、风控和归档全部按相同`Card`平铺，缺少首屏概览、模块主次和表单分组；规则字段数量增加后视觉密度失衡，排行榜也缺少运营视角的标题和摘要。
- `FB-5`根因：运营动作卡片位于右侧窄栏，但`.action-item`仍使用“两列网格 + 右侧控件自动宽度”布局；“结算归档”的输入框和按钮占据过多横向空间后，左侧说明文案剩余宽度接近 0，导致中文逐字换行。
- `FB-6`根因：后端`IssueCertificates`服务和`POST /plugins/sicau-niu/settlement/certificates/issue`接口已存在，但运营结算页只提供`InputNumber`让用户手输`honorId`；结算权限下也没有证书候选投影接口，导致操作者既不知道可发放的证书 ID，也可能误选服务端不支持批量结算的收藏类证书，从用户视角表现为“批量发证无法执行”。
- `FB-7`根因：牛只发布同时维护`release_stage`和`online_at`两套状态来源，玩家地图还要求阶段非空且上线时间已到达；Mock 数据只写了阶段、未写上线时间，导致地图为空。激活接口只按牛只 ID 加锁并做距离校验，没有复用地图可见性判断，玩家可绕过地图直接激活未上线牛只。

### 影响分析

- 已读取规则：`AGENTS.md`、`.agents/rules/openspec.md`、`.agents/rules/documentation.md`、`.agents/rules/architecture.md`、`.agents/rules/data-permission.md`、`.agents/rules/plugin.md`、`.agents/rules/api-contract.md`、`.agents/rules/backend-go.md`、`.agents/rules/database.md`、`.agents/rules/cache-consistency.md`、`.agents/rules/frontend-ui.md`、`.agents/rules/testing.md`、`.agents/rules/i18n.md`。
- 插件本地规范：`apps/lina-plugins/sicau-niu/AGENTS.md`不存在，按仓库顶层规范和规则文件执行。
- `i18n`影响：`sicau-niu`插件未在`plugin.yaml`启用`i18n`，新增单语言中文页面文案、菜单按钮和接口文档源文本不要求维护插件`manifest/i18n`资源。
- 缓存一致性影响：本次不新增缓存；规则配置以数据库表为权威数据源，更新后无缓存失效路径。
- 数据权限影响：新增运营端读取和更新接口受宿主认证、租户和权限中间件保护；玩家端数据仍只经玩家 token 访问。新增排行榜运营投影为活动运营全局统计，按平台级插件能力处理。
- 开发工具跨平台影响：不修改脚本、`Makefile`、`make.cmd`、`CI`或`linactl`；数据库生成使用既有插件`make dao`入口。
- DI 来源检查：新增规则服务由`backend/plugin.go`在路由装配期创建一次，并显式传入依赖的业务服务和控制器，共享同一数据库后端，不在请求路径临时构建服务图。
- `FB-6`补充影响：新增证书候选接口为插件自有运营结算 API，使用既有`settlementSvc`，不新增运行期依赖；候选查询仅读取插件自有荣誉定义表，受`只有发证权限可见`约束，不暴露玩家或租户业务数据；查询为单次投影并限制最多`100`条，避免前端逐项补查；不新增缓存、SQL、DAO、脚本或宿主核心契约。
- `FB-7`补充影响：移除插件自有牛只目录 API、管理页和数据库中的`releaseStage/release_stage`，不修改`lina-core`宿主契约；玩家地图与激活接口统一按`onlineAt`和可见窗口判断。数据权限边界不变，运营牛只管理仍受既有宿主权限码保护，玩家激活仍由玩家 token 隔离。`sicau-niu`为单语言插件，不要求新增`manifest/i18n`资源；不新增缓存或运行期依赖；不修改开发工具或脚本，使用既有根`db.init`、插件 SQL 和`make -C apps/lina-plugins dao p=sicau-niu`入口。

### 执行记录

- `FB-1`修复：升级`posterrender`和`certrender`默认实现，输出带活动标题、玩家/奖项字段、编号、纪念标识、版式边框和个性化装饰的`PNG`；新增`rendertext`工具，优先使用`SICAU_NIU_RENDER_FONT`或系统`CJK`字体，缺失字体时回退基础字体并保持可解码输出；补充海报和证书渲染器像素级单测，防止退回纯占位图。
- `FB-2`修复：新增`plugin_sicau_niu_rule_config`迁移、卸载 SQL、DAO/DO/Entity 生成入口和`rules`服务；运营端新增规则查询/更新接口，关键玩家和运营路径读取同一个运行规则服务，覆盖`LBS`判距、海报/证书标识、签到草量、偷草/送草限制、铁牛加成阈值、排行榜`TopN`、异常阈值和小程序回跳地址。
- `FB-3`修复：运营端新增喂草榜、好友榜、院系榜接口和结算页“排行榜数据”表格区；前端客户端、页面对象和`TC004`可见性断言同步更新。
- `FB-4`修复：将运营结算页从同权重`Card`平铺重排为运营工作台结构；新增顶部“寻牛活动运营概览”、激活进度、风险状态标签和关键指标区；活跃、规则、关键互动、运营动作、排行榜、风控和归档按运营任务分区展示；规则表单按“激活与展示”“签到与送草”“偷草策略”“榜单与风控”分组，排行榜增加榜首、记录数和领先院系摘要；保留既有接口、权限码和`data-testid`测试锚点。
- `FB-5`修复：运营动作区改为纵向动作项布局，不再让说明文案和右侧输入控件竞争同一行宽；批量发证与结算归档控件在动作项内部按可用宽度排列，归档输入框使用`minmax`自适应，说明文案保持正常横排可读；同步补充`TC004`对归档说明文本宽度的布局断言。
- `FB-6`修复：核对`功能模块参数.md`和`需求.md`后确认需求语义为“活动结束后/公示结算阶段支持证书批量发放”，不是每日凌晨自动发证；保留运营复核后的手动统一发证动作，新增结算权限下的可发证书候选接口，仅返回`participation/feed_count/activation_count`等后端支持批量结算的证书；运营结算页将手输`honorId`改为证书下拉选择，并将文案调整为“结算发证/统一发证”，避免误解为日常自动或每日手动任务；`TC004`补充选择证书并执行发证的端到端断言。
- `FB-7`修复：将牛只发布契约收敛为`onlineAt`唯一控制字段；管理端 DTO、Controller、`cattle`服务、Vue 客户端、表单、筛选和列表列均移除`releaseStage`，列表改展示“上线时间”。数据库安装 SQL 不再创建`release_stage`，当前迭代 SQL 幂等删除旧索引和旧列并创建`idx_sicau_niu_niu_online`，Mock 数据改为直接写`online_at`。玩家地图查询只按`online_at IS NOT NULL AND online_at <= now`做数据库过滤，再叠加周几/时段窗口；激活事务在行锁后复用同一可见性判断，未上线或窗口外牛只返回`CodeNiuNotVisible`且不写激活记录。

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
- 已执行并通过（`FB-6`）：`GOWORK=off go test ./backend/internal/service/settlement ./backend/internal/controller/settlement ./backend -count=1`（在`apps/lina-plugins/sicau-niu`内）。
- 已执行并通过（`FB-6`）：`pnpm -C apps/lina-vben/apps/web-antd typecheck`。
- 已执行并通过（`FB-6`）：`pnpm -C hack/tests test:validate`。
- 已执行并通过（`FB-6`）：`openspec validate complete-sicau-niu-feedback-gaps --strict`。
- 已执行并通过（`FB-6`）：`make dev`重新打包插件资源并重启本地后端/前端，确认新路由加载后`GET /x/sicau-niu/api/v1/plugins/sicau-niu/settlement/certificates/options`从`404`变为受权限保护的`401`。
- 已执行并通过（`FB-6`）：`pnpm -C hack/tests test:module -- plugin:sicau-niu -- --grep TC-4c`。
- 已执行并通过（`FB-6`）：`pnpm -C hack/tests test:module -- plugin:sicau-niu -- --grep TC-4`。
- 已执行并通过（`FB-7`）：`make db.init confirm=init`（宿主数据库初始化入口）。
- 已执行并通过（`FB-7`）：按顺序执行`apps/lina-plugins/sicau-niu/manifest/sql/001-*.sql`至`007-*.sql`，确认`007`可幂等删除`release_stage`并创建`idx_sicau_niu_niu_online`。
- 已执行并通过（`FB-7`）：`psql ... -f apps/lina-plugins/sicau-niu/manifest/sql/mock-data/001-sicau-niu-mock-data.sql`。
- 已执行并通过（`FB-7`）：`make -C apps/lina-plugins dao p=sicau-niu`。
- 已执行并通过（`FB-7`）：静态检索确认生产代码不再引用`releaseStage`、`ReleaseStage`或`release_stage`，仅 OpenSpec 说明和幂等删除 SQL 保留旧字段名。
- 已执行并通过（`FB-7`）：`GOWORK=off go test ./backend/internal/service/activation ./backend/internal/service/cattle ./backend/internal/controller/admin ./backend/internal/controller/player ./backend -count=1`（在`apps/lina-plugins/sicau-niu`内）。
- 已执行并通过（`FB-7`）：`GOWORK=off go test ./... -count=1`（在`apps/lina-plugins/sicau-niu`内）。
- 已执行并通过（`FB-7`）：`pnpm -C apps/lina-vben/apps/web-antd typecheck`。
- 已执行并通过（`FB-7`）：`pnpm -C hack/tests test:validate`。
- 已执行并通过（`FB-7`）：`openspec validate complete-sicau-niu-feedback-gaps --strict`。
- 已执行并通过（`FB-7`）：`make dev`重启本地后端/前端后，`pnpm -C hack/tests test:module -- plugin:sicau-niu -- --grep TC-2d`。
- 已执行并通过（`FB-7`）：数据库结构检查确认`plugin_sicau_niu_niu`仅存在`online_at`，不存在`release_stage`，并存在`idx_sicau_niu_niu_online`。
