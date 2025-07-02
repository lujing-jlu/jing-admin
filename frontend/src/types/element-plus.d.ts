declare module 'element-plus/dist/locale/zh-cn.mjs' {
  const zhCn: {
    name: string
    el: Record<string, string>
  }
  export default zhCn
}

declare module 'element-plus/es/components/form/src/types' {
  export interface RuleItem {
    required?: boolean
    message?: string
    trigger?: string | string[]
    min?: number
    max?: number
    type?: string
    validator?: (rule: RuleItem, value: any, callback: (error?: Error) => void) => void
  }
} 