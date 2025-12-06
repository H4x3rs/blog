# 阿里云OSS文件上传配置说明

本文档说明如何配置阿里云OSS文件上传功能。

## 功能说明

系统已集成阿里云OSS文件上传功能，支持：
- 图片上传（用于文章封面、编辑器图片等）
- 文件上传（支持PDF、Word、Excel、压缩包等）
- 自动返回OSS完整URL
- 支持自定义CDN域名

## 配置步骤

### 1. 创建阿里云OSS Bucket

1. 登录 [阿里云控制台](https://oss.console.aliyun.com/)
2. 进入"对象存储OSS"服务
3. 点击"创建Bucket"
4. 填写Bucket信息：
   - **Bucket名称**：例如 `my-blog-files`
   - **地域**：选择离你服务器最近的地域，例如 `华东1（杭州）`
   - **存储类型**：选择 `标准存储`
   - **读写权限**：选择 `公共读`（如果文件需要公开访问）或 `私有`（如果文件需要签名访问）
5. 点击"确定"创建Bucket

### 2. 获取AccessKey

1. 登录阿里云控制台
2. 鼠标悬停在右上角头像，点击"AccessKey管理"
3. 如果提示安全验证，完成验证
4. 创建AccessKey（如果还没有）
5. 记录以下信息：
   - **AccessKey ID**
   - **AccessKey Secret**

**安全提示**：AccessKey Secret只显示一次，请妥善保管。

### 3. 配置OSS域名（可选）

如果你有CDN域名，可以配置自定义域名：

1. 在OSS控制台，进入你的Bucket
2. 点击"传输管理" → "域名管理"
3. 绑定自定义域名（需要先完成域名备案和CNAME配置）
4. 记录自定义域名，例如 `cdn.example.com`

### 4. 配置后端

#### 方法一：使用环境变量（推荐）

在服务器上设置环境变量：

```bash
# OSS端点（地域节点）
export OSS_ENDPOINT="oss-cn-hangzhou.aliyuncs.com"

# OSS AccessKey ID
export OSS_ACCESS_KEY_ID="your_access_key_id"

# OSS AccessKey Secret
export OSS_ACCESS_KEY_SECRET="your_access_key_secret"

# OSS Bucket名称
export OSS_BUCKET_NAME="my-blog-files"

# OSS自定义域名（可选，如果配置了CDN）
export OSS_DOMAIN="https://cdn.example.com"
```

#### 方法二：修改配置文件

编辑 `backend/manifest/config/config.yaml`：

```yaml
oss:
  endpoint: "oss-cn-hangzhou.aliyuncs.com"  # OSS端点
  accessKeyId: "your_access_key_id"         # AccessKey ID
  accessKeySecret: "your_access_key_secret"  # AccessKey Secret
  bucketName: "my-blog-files"                # Bucket名称
  domain: "https://cdn.example.com"         # 自定义域名（可选）
```

**注意**：生产环境建议使用环境变量，避免敏感信息泄露。

### 5. 安装依赖

在 `backend` 目录下执行：

```bash
go mod tidy
```

这会自动下载阿里云OSS SDK依赖。

### 6. 启动服务

启动后端服务，系统会自动初始化OSS客户端：

```bash
cd backend
go run main.go
```

如果看到日志 `OSS服务初始化成功`，说明配置正确。

## OSS端点列表

根据你的Bucket所在地域，选择对应的端点：

| 地域 | 端点 |
|------|------|
| 华东1（杭州） | oss-cn-hangzhou.aliyuncs.com |
| 华东2（上海） | oss-cn-shanghai.aliyuncs.com |
| 华北1（青岛） | oss-cn-qingdao.aliyuncs.com |
| 华北2（北京） | oss-cn-beijing.aliyuncs.com |
| 华北3（张家口） | oss-cn-zhangjiakou.aliyuncs.com |
| 华北5（呼和浩特） | oss-cn-huhehaote.aliyuncs.com |
| 华南1（深圳） | oss-cn-shenzhen.aliyuncs.com |
| 华南2（河源） | oss-cn-heyuan.aliyuncs.com |
| 华南3（广州） | oss-cn-guangzhou.aliyuncs.com |
| 西南1（成都） | oss-cn-chengdu.aliyuncs.com |
| 中国（香港） | oss-cn-hongkong.aliyuncs.com |
| 美国（硅谷） | oss-us-west-1.aliyuncs.com |
| 美国（弗吉尼亚） | oss-us-east-1.aliyuncs.com |
| 新加坡 | oss-ap-southeast-1.aliyuncs.com |
| 澳大利亚（悉尼） | oss-ap-southeast-2.aliyuncs.com |
| 马来西亚（吉隆坡） | oss-ap-southeast-3.aliyuncs.com |
| 印度尼西亚（雅加达） | oss-ap-southeast-5.aliyuncs.com |
| 日本（东京） | oss-ap-northeast-1.aliyuncs.com |
| 韩国（首尔） | oss-ap-northeast-2.aliyuncs.com |
| 印度（孟买） | oss-ap-south-1.aliyuncs.com |
| 德国（法兰克福） | oss-eu-central-1.aliyuncs.com |
| 英国（伦敦） | oss-eu-west-1.aliyuncs.com |
| 阿联酋（迪拜） | oss-me-east-1.aliyuncs.com |

## 文件访问URL格式

### 使用OSS默认域名

```
https://bucket-name.endpoint/object-key
```

例如：
```
https://my-blog-files.oss-cn-hangzhou.aliyuncs.com/uploads/images/20240101_123456_photo.jpg
```

### 使用自定义域名（CDN）

```
https://cdn.example.com/uploads/images/20240101_123456_photo.jpg
```

## 文件存储结构

上传的文件在OSS中的存储结构：

```
bucket-name/
├── uploads/
│   ├── images/          # 图片文件
│   │   └── 20240101_123456_photo.jpg
│   └── files/           # 其他文件
│       └── 20240101_123456_document.pdf
```

## 前端使用

前端上传图片后，会收到OSS的完整URL，可以直接使用：

```javascript
// 上传图片
import { uploadImage } from '@/api/upload'

const handleUpload = async (file) => {
  try {
    const res = await uploadImage(file)
    // res.data.url 是OSS的完整URL
    console.log('图片URL:', res.data.url)
    // 例如: https://my-blog-files.oss-cn-hangzhou.aliyuncs.com/uploads/images/20240101_123456_photo.jpg
  } catch (error) {
    console.error('上传失败:', error)
  }
}
```

## 常见问题

### Q: 上传失败，提示"OSS配置不完整"

A: 检查以下配置：
1. 环境变量或配置文件中的OSS配置是否完整
2. AccessKey ID和Secret是否正确
3. Bucket名称是否正确
4. Endpoint是否正确

### Q: 上传成功但无法访问文件

A: 检查以下设置：
1. Bucket的读写权限是否设置为"公共读"（如果文件需要公开访问）
2. 如果使用自定义域名，检查CNAME配置是否正确
3. 检查OSS的防盗链设置

### Q: 如何切换到本地存储？

A: 如果不想使用OSS，可以：
1. 不配置OSS相关环境变量或配置
2. 修改上传控制器代码，使用本地存储逻辑
3. 系统会在启动时提示OSS初始化失败，但不影响其他功能

### Q: 如何配置CDN加速？

A: 
1. 在阿里云CDN控制台创建加速域名
2. 将CDN域名指向OSS Bucket
3. 在OSS配置中设置 `domain` 为CDN域名
4. 文件URL将自动使用CDN域名

### Q: 文件大小限制是多少？

A: 
- 图片文件：最大5MB
- 其他文件：最大10MB

可以在 `backend/internal/controller/upload/upload.go` 中修改限制。

## 安全建议

1. **不要将AccessKey提交到代码仓库**
   - 使用环境变量或配置文件，并将配置文件添加到 `.gitignore`
   
2. **使用RAM子账号**
   - 创建RAM子账号，只授予OSS相关权限
   - 使用子账号的AccessKey，而不是主账号的AccessKey

3. **定期轮换AccessKey**
   - 定期更换AccessKey，提高安全性

4. **设置Bucket权限**
   - 根据实际需求设置Bucket的读写权限
   - 如果文件不需要公开访问，设置为"私有"，并使用签名URL

5. **启用OSS日志**
   - 在OSS控制台启用访问日志，监控文件访问情况

## 技术支持

如有问题，请参考：
- [阿里云OSS官方文档](https://help.aliyun.com/product/31815.html)
- [Go SDK文档](https://help.aliyun.com/document_detail/32100.html)

