// This file implements host-managed OpenAPI document construction from the
// current host route table and plugin route projections.

package apidoc

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/util/gmeta"

	"lina-core/pkg/logger"
	pluginbridge "lina-core/pkg/pluginbridge/contract"
	"lina-core/pkg/pluginhost"
)

// openAPIAccessMetaKey is the plugin route metadata key used to override the
// document-level authentication default for public source-plugin routes.
const openAPIAccessMetaKey = "access"

// Build builds one host-managed OpenAPI document from the current route table
// and current plugin enablement state.
func (s *serviceImpl) Build(ctx context.Context, server *ghttp.Server) (*goai.OpenApiV3, error) {
	if server == nil {
		return nil, gerror.New("apidoc: host server is nil")
	}

	document := s.newDocument(ctx)
	sourceRouteBindings := s.listSourceRouteBindings()
	sourceRouteKeySet := buildSourceRouteKeySet(sourceRouteBindings)

	if err := s.addHostStaticRoutes(document, server, sourceRouteKeySet); err != nil {
		return nil, err
	}
	s.addEnabledSourceRoutes(ctx, document, sourceRouteBindings)
	if s.pluginSvc != nil {
		if err := s.pluginSvc.ProjectDynamicRoutesToOpenAPI(ctx, document.Paths); err != nil {
			return nil, err
		}
	}
	s.localizeDocument(ctx, document)
	return document, nil
}

// newDocument creates the baseline host-managed OpenAPI document and applies
// configured document metadata and shared security defaults.
func (s *serviceImpl) newDocument(ctx context.Context) *goai.OpenApiV3 {
	document := goai.New()
	if document.Paths == nil {
		document.Paths = goai.Paths{}
	}

	document.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	document.Config.CommonResponseDataField = "Data"
	document.Components.SecuritySchemes = goai.SecuritySchemes{
		"BearerAuth": goai.SecuritySchemeRef{
			Value: &goai.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT Bearer Token Authentication",
				In:           "header",
				Name:         "Authorization",
			},
		},
	}
	document.Security = &goai.SecurityRequirements{
		{"BearerAuth": {}},
	}

	if s == nil || s.configSvc == nil {
		return document
	}
	oaiCfg := s.configSvc.GetOpenApi(ctx)
	if oaiCfg == nil {
		return document
	}
	document.Info.Title = oaiCfg.Title
	document.Info.Description = oaiCfg.Description
	document.Info.Version = oaiCfg.Version
	serverURL := strings.TrimSpace(oaiCfg.ServerUrl)
	if serverURL == "" {
		serverURL = "/"
	}
	document.Servers = &goai.Servers{
		{
			URL:         serverURL,
			Description: oaiCfg.ServerDescription,
		},
	}
	return document
}

// addHostStaticRoutes projects host-owned strict routes that are not shadowed
// by source-plugin bindings into the output OpenAPI document.
func (s *serviceImpl) addHostStaticRoutes(
	document *goai.OpenApiV3,
	server *ghttp.Server,
	sourceRouteKeySet map[string]struct{},
) error {
	if document == nil || server == nil {
		return nil
	}
	for _, route := range server.GetRoutes() {
		if !shouldIncludeHostStaticRoute(route, sourceRouteKeySet) {
			continue
		}
		if err := addHandlerRouteToOpenAPI(
			document, route.Route, route.Method, route.Handler.Info.Value.Interface(),
		); err != nil {
			return err
		}
	}
	return nil
}

// addEnabledSourceRoutes projects documentable source-plugin routes for the
// plugins that are currently enabled.
func (s *serviceImpl) addEnabledSourceRoutes(
	ctx context.Context,
	document *goai.OpenApiV3,
	bindings []pluginhost.SourceRouteBinding,
) {
	if document == nil || len(bindings) == 0 {
		return
	}

	projectedRouteSet := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		if !binding.Documentable {
			continue
		}
		if s.pluginSvc != nil && !s.pluginSvc.IsEnabled(ctx, binding.PluginID) {
			continue
		}

		key := binding.Key()
		if _, ok := projectedRouteSet[key]; ok {
			continue
		}
		if err := addHandlerRouteToOpenAPI(document, binding.Path, binding.Method, binding.Handler); err != nil {
			logger.Warningf(
				ctx,
				"project source plugin route to OpenAPI failed plugin=%s method=%s path=%s err=%v",
				binding.PluginID,
				binding.Method,
				binding.Path,
				err,
			)
			continue
		}
		projectedRouteSet[key] = struct{}{}
	}
}

// listSourceRouteBindings reads the current source-plugin route binding snapshot
// from the plugin service when available.
func (s *serviceImpl) listSourceRouteBindings() []pluginhost.SourceRouteBinding {
	if s == nil || s.pluginSvc == nil {
		return nil
	}
	return s.pluginSvc.ListSourceRouteBindings()
}

// shouldIncludeHostStaticRoute reports whether the host route should stay in
// the document after removing plugin-owned strict-route duplicates.
func shouldIncludeHostStaticRoute(route ghttp.RouterItem, sourceRouteKeySet map[string]struct{}) bool {
	if route.Handler == nil || !route.Handler.Info.IsStrictRoute {
		return false
	}
	if _, ok := sourceRouteKeySet[buildRouteKey(route.Method, route.Route)]; ok {
		return false
	}
	return true
}

// addHandlerRouteToOpenAPI expands the handler's method set and registers it
// into the target OpenAPI document.
func addHandlerRouteToOpenAPI(
	document *goai.OpenApiV3,
	path string,
	method string,
	handler interface{},
) error {
	if document == nil {
		return nil
	}
	methods := expandOpenAPIMethods(method)
	publicAccess := handlerDeclaresPublicAccess(handler)
	for _, item := range methods {
		if err := document.Add(goai.AddInput{
			Path:   path,
			Method: item,
			Object: handler,
		}); err != nil {
			return err
		}
		if publicAccess {
			clearOpenAPIOperationSecurity(document, path, item)
		}
	}
	return nil
}

// handlerDeclaresPublicAccess reports whether one standard GoFrame handler
// marks its request DTO as `access:"public"`.
func handlerDeclaresPublicAccess(handler interface{}) bool {
	reqObject, ok := newOpenAPIHandlerReqObject(handler)
	if !ok {
		return false
	}
	accessMode := strings.TrimSpace(gmeta.Get(reqObject, openAPIAccessMetaKey).String())
	return strings.EqualFold(accessMode, pluginbridge.AccessPublic)
}

// newOpenAPIHandlerReqObject allocates the request DTO used for route metadata
// inspection when the handler follows GoFrame's standard API shape.
func newOpenAPIHandlerReqObject(handler interface{}) (interface{}, bool) {
	reflectType := reflect.TypeOf(handler)
	if !isOpenAPIHandlerTypeDocumentable(reflectType) {
		return nil, false
	}
	return reflect.New(reflectType.In(1).Elem()).Interface(), true
}

// isOpenAPIHandlerTypeDocumentable verifies the `(context.Context, *Req) (*Res, error)`
// shape required by GoFrame OpenAPI generation.
func isOpenAPIHandlerTypeDocumentable(reflectType reflect.Type) bool {
	if reflectType == nil || reflectType.Kind() != reflect.Func {
		return false
	}
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		return false
	}
	if !reflectType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return false
	}
	if !reflectType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return false
	}
	return reflectType.In(1).Kind() == reflect.Pointer && reflectType.In(1).Elem().Kind() == reflect.Struct
}

// clearOpenAPIOperationSecurity sets an explicit empty security requirement on
// one operation so it does not inherit the document-level JWT Bearer default.
func clearOpenAPIOperationSecurity(document *goai.OpenApiV3, path string, method string) {
	if document == nil || document.Paths == nil {
		return
	}
	pathKey, pathItem, ok := findOpenAPIPathItem(document, path)
	if !ok {
		return
	}
	operation := openAPIOperationForMethod(&pathItem, method)
	if operation == nil {
		return
	}
	emptySecurity := goai.SecurityRequirements{}
	operation.Security = &emptySecurity
	document.Paths[pathKey] = pathItem
}

// findOpenAPIPathItem retrieves a path item using the exact key first and the
// normalized key as a fallback for captured source-plugin routes.
func findOpenAPIPathItem(document *goai.OpenApiV3, path string) (string, goai.Path, bool) {
	if document == nil || document.Paths == nil {
		return "", goai.Path{}, false
	}
	if pathItem, ok := document.Paths[path]; ok {
		return path, pathItem, true
	}
	normalizedPath := normalizeOpenAPIPath(path)
	pathItem, ok := document.Paths[normalizedPath]
	return normalizedPath, pathItem, ok
}

// openAPIOperationForMethod returns the operation pointer that belongs to the
// normalized HTTP method on a path item.
func openAPIOperationForMethod(pathItem *goai.Path, method string) *goai.Operation {
	if pathItem == nil {
		return nil
	}
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodConnect:
		return pathItem.Connect
	case http.MethodDelete:
		return pathItem.Delete
	case http.MethodGet:
		return pathItem.Get
	case http.MethodHead:
		return pathItem.Head
	case http.MethodOptions:
		return pathItem.Options
	case http.MethodPatch:
		return pathItem.Patch
	case http.MethodPost:
		return pathItem.Post
	case http.MethodPut:
		return pathItem.Put
	case http.MethodTrace:
		return pathItem.Trace
	default:
		return nil
	}
}

// buildSourceRouteKeySet builds one lookup set for source-plugin route keys.
func buildSourceRouteKeySet(bindings []pluginhost.SourceRouteBinding) map[string]struct{} {
	items := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		items[binding.Key()] = struct{}{}
	}
	return items
}

// buildRouteKey combines one method and path into the normalized route key used
// by host and plugin route de-duplication.
func buildRouteKey(method string, path string) string {
	return strings.ToUpper(strings.TrimSpace(method)) + " " + normalizeOpenAPIPath(path)
}

// normalizeOpenAPIPath canonicalizes an OpenAPI path for stable key comparison.
func normalizeOpenAPIPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	return strings.TrimRight(trimmed, "/")
}

// expandOpenAPIMethods expands ALL or empty methods to the full supported HTTP
// method set used by GoFrame OpenAPI generation.
func expandOpenAPIMethods(method string) []string {
	normalized := strings.ToUpper(strings.TrimSpace(method))
	if normalized == "" || normalized == "ALL" {
		methods := ghttp.SupportedMethods()
		items := make([]string, 0, len(methods))
		for _, item := range methods {
			items = append(items, strings.ToUpper(strings.TrimSpace(item)))
		}
		return items
	}
	return []string{normalized}
}
