## ADDED Requirements

### Requirement: 卡片运营维护
系统 SHALL 在运营后台提供校史卡片的增删改查,受宿主统一权限校验。卡片包含分类(人物/事件/科研/院系/精神,后端命名常量约束)、标题、文案、图片路径。图片通过宿主文件管理上传得到存储路径后引用,本能力不自管文件存储。响应时间点字段 MUST 返回 Unix 毫秒。

#### Scenario: 新增卡片
- **WHEN** 运营携带 `sicau-niu:card:create` 权限提交合法分类、标题、文案,并可选图片路径
- **THEN** 系统创建该卡片

#### Scenario: 非法分类
- **WHEN** 提交的卡片分类不在允许枚举内
- **THEN** 系统返回 `bizerr` 参数错误,不创建

#### Scenario: 卡片列表分页
- **WHEN** 运营携带 `sicau-niu:card:list` 权限分页查询卡片
- **THEN** 系统在数据库侧过滤/排序/分页返回有界结果

#### Scenario: 无权限访问被拒绝
- **WHEN** 调用方不具备对应 `sicau-niu:card:*` 权限
- **THEN** 系统由统一权限中间件拒绝访问

### Requirement: 1 牛 1 卡约束
系统 SHALL 保证每头牛至多绑定一张主卡:卡片归属的牛在活跃集合内唯一。

#### Scenario: 为未绑卡的牛绑定主卡
- **WHEN** 运营创建/更新一张卡片并将其归属到一头尚无主卡的牛
- **THEN** 系统保存该归属

#### Scenario: 重复绑定同一头牛
- **WHEN** 运营将卡片归属到一头已拥有主卡的牛
- **THEN** 系统返回 `bizerr` 业务错误,拒绝重复绑定

#### Scenario: 归属到不存在的牛
- **WHEN** 运营将卡片归属到不存在的牛
- **THEN** 系统返回 `bizerr` 业务错误,不保存

### Requirement: 校史金句运营维护
系统 SHALL 在运营后台提供校史金句的增删改查,作为后续喂草/点击随机播放的金句池,受宿主统一权限校验。

#### Scenario: 新增金句
- **WHEN** 运营携带 `sicau-niu:quote:create` 权限提交非空金句文案
- **THEN** 系统创建该金句

#### Scenario: 金句列表分页
- **WHEN** 运营携带 `sicau-niu:quote:list` 权限分页查询金句
- **THEN** 系统在数据库侧分页返回有界结果

#### Scenario: 空金句被拒绝
- **WHEN** 提交的金句文案为空
- **THEN** 系统返回 `bizerr` 参数错误,不创建
