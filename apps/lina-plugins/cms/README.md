# LinaPro CMS Plugin

The CMS source plugin provides content site capabilities for LinaPro. It owns site settings, category trees, articles, links, slides, visitor messages, and public HTML rendering APIs.

## Scope

- Source plugin under `apps/lina-plugins/cms`
- Plugin-owned `plugin_cms_*` tables
- Authenticated management APIs under `/api/v1/cms/*`
- Public read-only APIs under `/api/v1/cms/public/*`
- Public HTML pages under `/cms-site`
- Vben management page embedded through the plugin runtime
- Plugin-owned E2E tests under `hack/tests`

## Public Templates

Public templates live in `public/templates`. They use CMS-owned template tags with the `cms:` prefix, for example `cms:list`, `cms:nav`, `cms:content`, `cms:prev`, and `cms:next`. Static assets used by the reference site live in `public/assets`.

## Development

The plugin is generated and tested as an independent source plugin. Database code generation uses `backend/hack/config.yaml`, and the generated DAO, DO, and Entity files stay inside the plugin `backend/internal` tree.
