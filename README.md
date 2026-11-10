# RememberMe — 数字纪念馆

基于家族树的数字纪念平台，支持墓园管理、逝者信息、家族关系图可视化。

## 技术栈

| 层 | 技术 |
|----|------|
| **后端** | Go 1.25+ · Gin · GORM · MySQL 8.0 |
| **前端** | React 18 · TypeScript · Vite · Ant Design 5 |
| **地图** | 高德地图 JS API 2.0 · @amap/amap-jsapi-loader |
| **认证** | JWT (access token + refresh token) |

## 目录结构

```
.
├── server/                    # Go 后端
│   ├── main.go                # 入口
│   ├── config.yaml            # 配置文件
│   ├── handler/               # HTTP 处理器
│   ├── service/               # 业务逻辑
│   ├── repository/            # 数据访问
│   ├── model/                 # 数据模型
│   ├── middleware/             # 中间件（JWT 等）
│   ├── routes/                # 路由注册
│   ├── config/                # 配置加载
│   ├── types/                 # 响应类型定义
│   ├── utils/                 # 工具（JWT/验证/响应）
│   ├── errors/                # 业务错误定义
│   ├── errcode/               # 错误码
│   └── uploads/               # 上传文件存储目录（运行时生成）
│       ├── cover/             # 墓园封面图
│       └── portrait/          # 逝者遗像
│
├── rememberme-web/            # React 前端
│   ├── src/
│   │   ├── api/               # API 调用封装
│   │   ├── components/        # 组件
│   │   │   ├── deceased/      #   - 逝者相关组件
│   │   │   ├── grave/         #   - 墓园相关组件
│   │   │   ├── relation/      #   - 关系管理组件
│   │   │   ├── relation-graph/#   - 关系图组件
│   │   │   ├── layout/        #   - 布局组件
│   │   │   └── member/        #   - 成员管理组件
│   │   ├── pages/             # 页面
│   │   │   ├── auth/          #   - 登录/注册
│   │   │   ├── dashboard/     #   - 墓园列表
│   │   │   ├── grave/         #   - 墓园详情
│   │   │   ├── members/       #   - 成员管理
│   │   │   └── profile/       #   - 个人中心
│   │   ├── types/             # TypeScript 类型定义
│   │   ├── contexts/          # React Context
│   │   ├── utils/             # 工具函数
│   │   ├── App.tsx            # 路由配置
│   │   └── main.tsx           # 入口
│   ├── vite.config.ts         # Vite 配置（含代理）
│   └── package.json
│
└── README.md                  # 本文件
```

## 快速开始

### 环境要求

- Go 1.25+
- Node.js 18+
- MySQL 8.0+

### 1. 启动数据库

确保 MySQL 运行中，创建数据库：

```sql
CREATE DATABASE RememberMe CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 配置后端

编辑 `server/config.yaml`：

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: ""
  name: RememberMe
```

敏感配置（数据库密码、JWT Secret）通过环境变量覆盖，避免硬编码在 yaml 中：

```bash
# 复制模板并填入真实值
cp server/.env.example server/.env
```

```env
# server/.env
DB_PASSWORD=your_password_here
JWT_SECRET=your_jwt_secret_here
```

### 3. 启动后端

```bash
cd server
go run main.go
# 监听 :8080，GORM 自动建表
```

> **开发调试提示**：`config.yaml` 中 `app.gin_mode: release` 会禁用控制台日志，开发时建议改为 `debug` 方便调试。JWT `access_token` 过期时间默认 **15 分钟**，`refresh_token` 默认 **7 天**。

### 4. 配置前端环境变量

编辑 `rememberme-web/.env`，填写高德地图 Key 和安全密钥：

```env
VITE_AMAP_KEY=你的高德Key
VITE_AMAP_SECRET=你的安全密钥
```

获取方式：[高德开放平台控制台](https://console.amap.com/dev/key/app) → 创建应用 → Web端(JS API)。

### 5. 安装前端依赖

```bash
cd rememberme-web
npm install
```

### 6. 启动前端

```bash
npm run dev
# 监听 :3000，自动代理 /api 和 /uploads 到后端 :8080
```

### 7. 访问

浏览器打开 `http://localhost:3000`，注册账号即可开始使用。

## 开发说明

### 前端代理配置

`vite.config.ts` 中配置了开发代理：

```
/api/*     → http://localhost:8080   (API 请求)
/uploads/* → http://localhost:8080   (静态文件访问)
```

### 图片上传

- 图片通过 `POST /api/v1/upload` 接口上传
- `subDir` 参数指定子目录：`cover`（封面）或 `portrait`（遗像）
- 服务端保存在 `server/uploads/` 目录下
- 返回格式：`{ "url": "/uploads/portrait/uuid.png" }`

### API 基础路径

所有 API 以 `/api/v1/` 为前缀，需要 JWT 认证（`Authorization: Bearer <token>`）。

### 编译检查

```bash
# 后端
cd server && go build ./...

# 前端
cd rememberme-web && npx tsc --noEmit
```

> 注意：`gin_mode: release` 下 `go build` 的产物启动后无控制台日志，生产部署前确认日志写入文件路径。

### 地图选点功能

墓园创建和编辑表单中集成了高德地图选点组件，支持以下操作：

| 功能 | 实现方式 |
|------|----------|
| **关键字搜索** | 调用 `AMap.PlaceSearch`，输入地名后搜索并展示结果列表 |
| **浏览器定位** | 调用 `HTML5 Geolocation`，点击定位按钮获取当前 GPS 坐标 |
| **地图点击选点** | 点击地图任意位置放置标记，自动逆地理编码获取地址 |
| **手动输入坐标** | 直接输入经纬度，点击定位到坐标按钮跳转 |
| **拖动标记** | 标记点可拖动调整，松开后自动更新地址 |
| **清除标记** | 一键清除已选位置 |

#### 相关文件

```
rememberme-web/src/
├── components/grave/
│   ├── MapPicker.tsx          # 地图选点组件（SDK 加载、搜索、定位、标记）
│   ├── GraveForm.tsx           # 创建/编辑表单（集成 MapPicker）
│   └── GraveDetail.tsx         # 墓园详情（展示位置信息）
└── types/
    └── grave.ts                # 类型定义（LatLng 接口）
```

#### 高德 SDK 加载

使用 `@amap/amap-jsapi-loader` 动态加载，启动时输出诊断日志：

```
[Amap] loading SDK, key=xxx..., secret configured=true
[Amap PlaceSearch] status: complete result: {info: 'OK', ...}
```

#### 免费配额

| 服务 | 个人配额/月 | 说明 |
|------|:----------:|------|
| JS 地图初始化 | 150万次 | 页面加载展示地图 |
| 浏览器定位 | 含在地图配额 | Geolocation 插件 |
| 逆地理编码 | 15万次 | 坐标 → 地址 |
| POI 关键字搜索 | 5,000次 | 地点搜索 |

超出后可购买流量包（30元/万次），详见 [高德定价页](https://lbs.amap.com/upgrade#price)。
