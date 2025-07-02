# Jing Admin 设计系统

这是 Jing Admin 项目的统一设计系统，确保所有页面和组件具有一致的视觉风格和用户体验。

## 🎨 设计原则

- **一致性**: 所有组件使用统一的颜色、字体、间距等
- **简洁性**: 简洁明了的界面设计，避免冗余元素
- **响应式**: 适配各种屏幕尺寸
- **可访问性**: 良好的色彩对比度和键盘导航支持

## 📁 文件结构

```
src/styles/
├── design-system.css  # 主设计系统文件
└── README.md         # 本文档
```

## 🎯 使用方法

### 1. 导入设计系统

设计系统已在 `main.ts` 中全局导入：

```typescript
import './styles/design-system.css'
```

### 2. 页面布局类

#### 页面容器
```html
<div class="page-container">
  <!-- 页面内容 -->
</div>
```

#### 页面头部
```html
<div class="page-header">
  <h2 class="page-title">页面标题</h2>
  <div class="page-actions">
    <el-button type="primary">操作按钮</el-button>
  </div>
</div>
```

#### 卡片容器
```html
<div class="card-container">
  <!-- 卡片内容 -->
</div>
```

#### 搜索区域
```html
<div class="search-section">
  <el-form class="search-form">
    <!-- 搜索表单 -->
  </el-form>
</div>
```

#### 表格区域
```html
<div class="table-section">
  <el-table>
    <!-- 表格内容 -->
  </el-table>
  
  <div class="pagination-section">
    <el-pagination />
  </div>
</div>
```

## 🎨 颜色系统

### 主色调
- `--primary-color`: #3b82f6 (主要按钮、链接)
- `--primary-light`: #60a5fa (悬停状态)
- `--primary-dark`: #1d4ed8 (按下状态)

### 状态颜色
- `--success-color`: #10b981 (成功状态)
- `--warning-color`: #f59e0b (警告状态)
- `--error-color`: #ef4444 (错误状态)
- `--info-color`: #6b7280 (信息状态)

### 背景色
- `--bg-page`: 页面背景渐变
- `--bg-primary`: #ffffff (主要背景)
- `--bg-secondary`: #f8fafc (次要背景)
- `--bg-card`: rgba(255, 255, 255, 0.9) (卡片背景)

### 文字色
- `--text-primary`: #1e293b (主要文字)
- `--text-secondary`: #64748b (次要文字)
- `--text-tertiary`: #94a3b8 (辅助文字)

## 📏 尺寸系统

### 间距
- `--spacing-xs`: 4px
- `--spacing-sm`: 8px
- `--spacing-md`: 12px
- `--spacing-lg`: 16px
- `--spacing-xl`: 20px
- `--spacing-2xl`: 24px

### 圆角
- `--radius-sm`: 4px
- `--radius-md`: 8px
- `--radius-lg`: 12px
- `--radius-xl`: 16px

### 字体大小
- `--font-xs`: 12px
- `--font-sm`: 14px
- `--font-md`: 16px
- `--font-lg`: 18px
- `--font-xl`: 20px
- `--font-2xl`: 24px

## 🌙 深色模式

设计系统支持自动深色模式切换，只需在根元素添加 `.dark-theme` 类：

```html
<html class="dark-theme">
```

## 📱 响应式断点

- `1200px` 以下：减少内边距
- `768px` 以下：垂直布局，调整字体大小
- `480px` 以下：移动优化布局

## 🧩 Element Plus 组件样式

设计系统自动覆盖 Element Plus 组件样式，确保：

- 表格背景透明，使用统一边框和颜色
- 按钮使用统一圆角和悬停效果
- 卡片使用毛玻璃效果和一致的阴影
- 表单元素使用统一的聚焦状态
- 分页组件使用统一的颜色方案

## 💡 最佳实践

### 页面结构模板

```html
<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">页面标题</h2>
      <div class="page-actions">
        <el-button type="primary">主要操作</el-button>
      </div>
    </div>

    <!-- 搜索区域（可选） -->
    <div class="search-section">
      <el-form class="search-form">
        <!-- 搜索表单项 -->
      </el-form>
    </div>

    <!-- 主要内容 -->
    <div class="table-section">
      <!-- 表格或其他内容 -->
      
      <!-- 分页（如果需要） -->
      <div class="pagination-section">
        <el-pagination />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 只添加页面特定的样式，使用设计系统变量 */
.custom-element {
  margin-top: var(--spacing-xl);
  border-radius: var(--radius-md);
  color: var(--text-primary);
}
</style>
```

### 样式编写原则

1. **优先使用类名**：使用预定义的类名而不是自定义样式
2. **使用CSS变量**：自定义样式时使用设计系统的CSS变量
3. **最小化自定义**：尽量减少页面特定的样式代码
4. **保持一致性**：新组件应该遵循现有的设计模式

### 示例对比

❌ **不推荐**：
```css
.my-header {
  padding: 20px;
  background: #ffffff;
  border-radius: 8px;
  margin-bottom: 24px;
}
```

✅ **推荐**：
```html
<div class="page-header">
  <!-- 直接使用设计系统类 -->
</div>
```

或者：
```css
.my-custom-element {
  padding: var(--spacing-xl);
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-2xl);
}
```

## 🔧 扩展设计系统

如需添加新的设计元素：

1. 在 `design-system.css` 中添加新的CSS变量
2. 添加相应的深色模式变量
3. 创建对应的类名
4. 更新本文档

## 📋 检查清单

新页面创建时请确保：

- [ ] 使用 `.page-container` 作为根容器
- [ ] 页面头部使用 `.page-header` 和 `.page-title`
- [ ] 搜索区域使用 `.search-section`
- [ ] 表格区域使用 `.table-section`
- [ ] 分页使用 `.pagination-section`
- [ ] 自定义样式使用设计系统变量
- [ ] 在深色模式下测试页面效果
- [ ] 在不同屏幕尺寸下测试响应式布局 