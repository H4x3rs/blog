package service

import (
	"context"
	"fmt"
)

// FunctionCall 函数调用结构
type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// FunctionCallResult 函数调用结果
type FunctionCallResult struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

// 获取可用的函数定义（tools）
func (s *sChat) GetAvailableFunctions() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_articles",
				"description": "查询文章列表，支持按标题、分类、状态等条件筛选",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "文章标题关键词（模糊匹配）",
						},
						"categoryId": map[string]interface{}{
							"type":        "integer",
							"description": "分类ID",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "文章状态：published（已发布）、draft（草稿）",
							"enum":        []string{"published", "draft"},
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_article",
				"description": "根据ID查询单篇文章详情",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "文章ID",
						},
					},
					"required": []string{"id"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_users",
				"description": "查询用户列表，支持按用户名模糊搜索",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "用户名关键词（模糊匹配）",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_user",
				"description": "根据ID查询单个用户详情",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "用户ID",
						},
					},
					"required": []string{"id"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_categories",
				"description": "查询分类列表，支持按名称、slug搜索",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "分类名称关键词（模糊匹配）",
						},
						"slug": map[string]interface{}{
							"type":        "string",
							"description": "分类slug关键词（模糊匹配）",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_category",
				"description": "根据ID或slug查询单个分类详情",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "分类ID",
						},
						"slug": map[string]interface{}{
							"type":        "string",
							"description": "分类slug",
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_tags",
				"description": "查询标签列表，支持按名称搜索",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "标签名称关键词（模糊匹配）",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_tag",
				"description": "根据ID或slug查询单个标签详情",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "标签ID",
						},
						"slug": map[string]interface{}{
							"type":        "string",
							"description": "标签slug",
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_menus",
				"description": "查询菜单列表，返回所有菜单的树形结构",
				"parameters": map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_permissions",
				"description": "查询权限列表，返回所有权限信息",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"parentId": map[string]interface{}{
							"type":        "integer",
							"description": "父权限ID，-1表示查询所有",
							"default":     -1,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_roles",
				"description": "查询角色列表，支持按名称搜索",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "角色名称关键词（模糊匹配）",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_topics",
				"description": "查询专题列表，支持按名称搜索",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "专题名称关键词（模糊匹配）",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_topic",
				"description": "根据ID查询单个专题详情",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "专题ID",
						},
					},
					"required": []string{"id"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_topic_articles",
				"description": "查询专题下的文章列表",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"topicId": map[string]interface{}{
							"type":        "integer",
							"description": "专题ID",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "页码，从1开始",
							"default":     1,
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "每页数量",
							"default":     10,
						},
					},
					"required": []string{"topicId"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_settings",
				"description": "查询系统设置，可以查询所有设置或指定key的设置",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"key": map[string]interface{}{
							"type":        "string",
							"description": "设置key，如果为空则返回所有设置",
						},
					},
				},
			},
		},
	}
}

// 执行函数调用
func (s *sChat) ExecuteFunction(ctx context.Context, userId int, functionCall *FunctionCall) (*FunctionCallResult, error) {
	switch functionCall.Name {
	case "query_articles":
		return s.executeQueryArticles(ctx, userId, functionCall.Arguments)
	case "query_article":
		return s.executeQueryArticle(ctx, userId, functionCall.Arguments)
	case "query_users":
		return s.executeQueryUsers(ctx, functionCall.Arguments)
	case "query_user":
		return s.executeQueryUser(ctx, functionCall.Arguments)
	case "query_categories":
		return s.executeQueryCategories(ctx, functionCall.Arguments)
	case "query_category":
		return s.executeQueryCategory(ctx, functionCall.Arguments)
	case "query_tags":
		return s.executeQueryTags(ctx, functionCall.Arguments)
	case "query_tag":
		return s.executeQueryTag(ctx, functionCall.Arguments)
	case "query_menus":
		return s.executeQueryMenus(ctx, functionCall.Arguments)
	case "query_permissions":
		return s.executeQueryPermissions(ctx, functionCall.Arguments)
	case "query_roles":
		return s.executeQueryRoles(ctx, functionCall.Arguments)
	case "query_topics":
		return s.executeQueryTopics(ctx, functionCall.Arguments)
	case "query_topic":
		return s.executeQueryTopic(ctx, functionCall.Arguments)
	case "query_topic_articles":
		return s.executeQueryTopicArticles(ctx, functionCall.Arguments)
	case "query_settings":
		return s.executeQuerySettings(ctx, functionCall.Arguments)
	default:
		return nil, fmt.Errorf("未知的函数: %s", functionCall.Name)
	}
}

// 查询文章列表
func (s *sChat) executeQueryArticles(ctx context.Context, userId int, args map[string]interface{}) (*FunctionCallResult, error) {
	title := ""
	if v, ok := args["title"].(string); ok {
		title = v
	}
	categoryId := 0
	if v, ok := args["categoryId"].(float64); ok {
		categoryId = int(v)
	}
	status := ""
	if v, ok := args["status"].(string); ok {
		status = v
	}
	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := Article.GetList(ctx, title, categoryId, status, page, size, userId, false, 0)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_articles",
		Content: map[string]interface{}{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}, nil
}

// 查询单篇文章
func (s *sChat) executeQueryArticle(ctx context.Context, userId int, args map[string]interface{}) (*FunctionCallResult, error) {
	id := 0
	if v, ok := args["id"].(float64); ok {
		id = int(v)
	}
	if id == 0 {
		return nil, fmt.Errorf("文章ID不能为空")
	}

	article, err := Article.GetOne(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name:    "query_article",
		Content: article,
	}, nil
}

// 查询用户列表
func (s *sChat) executeQueryUsers(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	username := ""
	if v, ok := args["username"].(string); ok {
		username = v
	}
	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := User.GetList(ctx, username, page, size)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_users",
		Content: map[string]interface{}{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}, nil
}

// 查询单个用户
func (s *sChat) executeQueryUser(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	id := 0
	if v, ok := args["id"].(float64); ok {
		id = int(v)
	}
	if id == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	user, err := User.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name:    "query_user",
		Content: user,
	}, nil
}

// 查询分类列表
func (s *sChat) executeQueryCategories(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	name := ""
	if v, ok := args["name"].(string); ok {
		name = v
	}
	slug := ""
	if v, ok := args["slug"].(string); ok {
		slug = v
	}
	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := Category.GetList(ctx, name, slug, page, size)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_categories",
		Content: map[string]interface{}{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}, nil
}

// 查询单个分类
func (s *sChat) executeQueryCategory(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	if id, ok := args["id"].(float64); ok && id > 0 {
		category, err := Category.GetOne(ctx, int(id))
		if err != nil {
			return nil, err
		}
		return &FunctionCallResult{
			Name:    "query_category",
			Content: category,
		}, nil
	}

	if slug, ok := args["slug"].(string); ok && slug != "" {
		category, err := Category.GetBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
		return &FunctionCallResult{
			Name:    "query_category",
			Content: category,
		}, nil
	}

	return nil, fmt.Errorf("分类ID或slug不能为空")
}

// 查询标签列表
func (s *sChat) executeQueryTags(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	name := ""
	if v, ok := args["name"].(string); ok {
		name = v
	}
	slug := ""
	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := Tag.GetList(ctx, name, slug, page, size)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_tags",
		Content: map[string]interface{}{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}, nil
}

// 查询单个标签
func (s *sChat) executeQueryTag(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	if id, ok := args["id"].(float64); ok && id > 0 {
		tag, err := Tag.GetOne(ctx, int(id))
		if err != nil {
			return nil, err
		}
		return &FunctionCallResult{
			Name:    "query_tag",
			Content: tag,
		}, nil
	}

	if slug, ok := args["slug"].(string); ok && slug != "" {
		tag, err := Tag.GetBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
		return &FunctionCallResult{
			Name:    "query_tag",
			Content: tag,
		}, nil
	}

	return nil, fmt.Errorf("标签ID或slug不能为空")
}

// 查询菜单列表
func (s *sChat) executeQueryMenus(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	list, err := Menu.GetTree(ctx)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name:    "query_menus",
		Content: list,
	}, nil
}

// 查询权限列表
func (s *sChat) executeQueryPermissions(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	parentId := -1
	if v, ok := args["parentId"].(float64); ok {
		parentId = int(v)
	}

	list, err := Permission.GetList(ctx, parentId)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name:    "query_permissions",
		Content: list,
	}, nil
}

// 查询角色列表
func (s *sChat) executeQueryRoles(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	name := ""
	if v, ok := args["name"].(string); ok {
		name = v
	}
	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := Role.GetList(ctx, name, page, size)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_roles",
		Content: map[string]interface{}{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}, nil
}

// 查询专题列表
func (s *sChat) executeQueryTopics(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	name := ""
	if v, ok := args["name"].(string); ok {
		name = v
	}
	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := Topic.GetList(ctx, name, page, size)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_topics",
		Content: map[string]interface{}{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}, nil
}

// 查询单个专题
func (s *sChat) executeQueryTopic(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	id := 0
	if v, ok := args["id"].(float64); ok {
		id = int(v)
	}
	if id == 0 {
		return nil, fmt.Errorf("专题ID不能为空")
	}

	topic, err := Topic.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name:    "query_topic",
		Content: topic,
	}, nil
}

// 查询专题下的文章列表
func (s *sChat) executeQueryTopicArticles(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	topicId := 0
	if v, ok := args["topicId"].(float64); ok {
		topicId = int(v)
	}
	if topicId == 0 {
		return nil, fmt.Errorf("专题ID不能为空")
	}

	page := 1
	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	size := 10
	if v, ok := args["size"].(float64); ok {
		size = int(v)
	}

	list, total, err := Topic.GetTopicArticles(ctx, topicId, page, size)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name: "query_topic_articles",
		Content: map[string]interface{}{
			"list":   list,
			"total":  total,
			"page":   page,
			"size":   size,
			"topicId": topicId,
		},
	}, nil
}

// 查询系统设置
func (s *sChat) executeQuerySettings(ctx context.Context, args map[string]interface{}) (*FunctionCallResult, error) {
	key := ""
	if v, ok := args["key"].(string); ok {
		key = v
	}

	if key != "" {
		// 查询单个设置
		value, err := Settings.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return &FunctionCallResult{
			Name: "query_settings",
			Content: map[string]interface{}{
				"key":   key,
				"value": value,
			},
		}, nil
	}

	// 查询所有设置
	allSettings, err := Settings.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &FunctionCallResult{
		Name:    "query_settings",
		Content: allSettings,
	}, nil
}

