export interface FactoryConfigItem {
  id: string
  created_at: string
  version: string
  name: string | null
  config: unknown
}

export type Page = 'home' | 'list' | 'detail'
