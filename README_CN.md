# MDC 规则文件拷贝工具

这个小工具用于根据你的需求拷贝指定的规则到你的项目里面

**这个项目能帮你快速让 Cursor 遵守你的指令！**

## 功能

- 根据映射关系拷贝文件
- 支持按分类和类型筛选拷贝（如仅拷贝后端语言中的 python 相关规则）
- 支持列出所有可用的规则分类和类型
- 使用配置文件定义映射关系
- 按字母顺序排序输出结果
- 内置默认配置文件（无需外部配置文件即可使用）

## 使用方法

### 前置条件

- 系统中已安装 Go 语言环境
  - 如果您尚未安装 Go，请访问 [https://go.dev/doc/install](https://go.dev/doc/install) 安装最新版本
  - 确保正确设置 GO_BIN 环境变量，以便使用 `go install` 命令

### 安装

您可以直接从 GitHub 安装：

```bash
go install github.com/zhaolion/cursor-mdc-work@latest
```

安装完成后，您可以直接使用 `cursor-mdc-work` 命令生成 Cursor MDC (Markdown Cursor) 规则文件。

### 手动编译

如果您更喜欢手动编译：

```bash
go build -o cursor-mdc-work
```

### 使用示例

1. 使用内置默认配置列出所有可用的规则分类和类型：

```bash
cursor-mdc-work list
```

2. 列出特定分类下的所有类型：

```bash
cursor-mdc-work list --category backend-languages
```

3. 拷贝特定分类下特定类型的规则文件：

```bash
cursor-mdc-work copy --category backend-languages --types python,go --target /path/to/destination
```

4. 拷贝特定分类下所有规则文件：

```bash
cursor-mdc-work copy --category backend-languages --target /path/to/destination
```

5. 拷贝所有规则文件：

```bash
cursor-mdc-work copy --target /path/to/destination
```

6. 使用自定义配置文件：

```bash
cursor-mdc-work copy --config my-mapping.json --target /path/to/destination
```

## 配置文件格式

配置文件格式为 JSON，示例如下：

```json
{
  "mappings": {
    "frontend-frameworks": {
      "react": [
        "cursor/rules-mdc/react.mdc",
        "cursor/rules-mdc/javascript.mdc",
        "cursor/rules-mdc/typescript.mdc",
        "cursor/rules-mdc/html.mdc",
        "cursor/rules-mdc/css.mdc"
      ]
    },
    "backend-languages": {
      "python": ["cursor/rules-mdc/python.mdc", "cursor/rules-mdc/fastapi.mdc"]
    }
  }
}
```

如果没有通过 `--config` 参数指定配置文件，工具将使用内置的默认配置。如果指定了配置文件但该文件不存在，程序会在指定位置创建一个示例配置文件。

## 当前支持的分类

该工具目前支持以下分类：

> 要获取最新完整的分类和类型列表，请运行 `cursor-mdc-work list`。

```
All rule categories and types:

ai-ml:
  - ai-agents
  - computer-vision
  - data-science
  - machine-learning
  - nlp

backend-languages:
  - go
  - java
  - node
  - other-backend
  - python
  - rust

build-tools:
  - bundlers
  - dependency-management

cloud-services:
  - aws
  - azure
  - gcp
  - other-cloud

databases:
  - nosql
  - sql

desktop-apps:
  - cross-platform
  - gui-frameworks

devops:
  - ci-cd
  - configuration-management
  - containerization
  - service-deployment

editors:
  - code-editors

frontend-frameworks:
  - angular
  - other-frontend
  - react
  - svelte
  - vue

frontend-ui:
  - ui-components
  - visualization

game-development:
  - game-engines
  - game-related

mobile:
  - flutter
  - other-mobile
  - react-native

testing:
  - api-testing
  - code-quality
  - testing-frameworks

utility-libraries:
  - api-integration
  - auth-security
  - documentation
  - monitoring
```
