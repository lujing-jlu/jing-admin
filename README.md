# Jing Admin

现代化全栈管理后台系统 | Modern Full-Stack Admin System

---

## 项目简介 | Project Overview

Jing Admin 是一个基于 Vue3 + Element Plus + TypeScript + Go + Gin + SQLite 的现代化开源管理后台系统，支持用户、角色、权限、日志、系统配置、文件上传、批量导入导出等核心功能，适合中小型企业、团队或个人快速搭建后台管理平台。

Jing Admin is a modern open-source admin system based on Vue3, Element Plus, TypeScript, Go, Gin, and SQLite. It features user/role/permission management, logs, system config, file upload, batch import/export, and more. Ideal for SMEs, teams, or individuals to quickly build admin platforms.

---

## 功能特性 | Features

- 用户、角色、权限（RBAC）管理
- 操作日志、系统配置、系统监控
- 文件上传（头像/文档/图片/通用）
- 批量导入导出（CSV）
- 主题切换、响应式布局、统一设计系统
- JWT认证、权限控制、接口安全
- 支持多语言（中/英）
- 易于二次开发和扩展

---

## 技术栈 | Tech Stack

- **前端**：Vue3, TypeScript, Element Plus, Vite, Pinia, Vue Router
- **后端**：Go, Gin, GORM, SQLite
- **工具**：Prettier, ESLint, Docker（可选）

---

## 快速开始 | Quick Start

### 1. 克隆项目 | Clone

```bash
git clone https://github.com/yourname/jing-admin.git
cd jing-admin
```

### 2. 启动后端 | Backend

```bash
cd backend
go mod tidy
go run main.go
```

默认端口：8081，数据库文件：backend/jing_admin.db

### 3. 启动前端 | Frontend

```bash
cd frontend
npm install
npm run dev
```

访问：http://localhost:5173

---

## 目录结构 | Directory Structure

```
jing-admin/
├── backend/    # Go后端服务
├── frontend/   # Vue3前端项目
└── README.md   # 项目说明
```

---

## 贡献指南 | Contributing

欢迎PR和Issue！请遵循以下规范：
- 前端：Prettier + ESLint
- 后端：Go标准代码规范
- 提交信息：Conventional Commits
- 详细贡献流程见CONTRIBUTING.md（如有）

---

## 许可证 | License

MIT License

---

## 联系与支持 | Contact & Support

- GitHub Issues

---

> 本项目持续迭代中，欢迎Star、Fork、反馈建议！

---

English version below.

---

# Jing Admin

A modern full-stack admin system based on Vue3 + Go. See above for features, tech stack, and quick start. Contributions welcome! 