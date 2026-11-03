import type { FactoryConfigItem } from '../types'

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/+$/, '')

export async function pingApi() {
  const res = await fetch(`${apiBaseUrl}/ping`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return await res.json()
}

export async function fetchFactoryConfigs(): Promise<FactoryConfigItem[]> {
  const res = await fetch(`${apiBaseUrl}/api/v1/factory-configs`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const data = await res.json()
  return data.items || []
}

export async function fetchFactoryConfigById(id: string): Promise<FactoryConfigItem> {
  const res = await fetch(`${apiBaseUrl}/api/v1/factory-configs/${id}`)
  if (!res.ok) throw new Error('Failed to fetch config details')
  return await res.json()
}

export async function uploadFactoryConfig(formData: FormData): Promise<void> {
  const res = await fetch(`${apiBaseUrl}/api/v1/factory-configs/upload`, {
    method: 'POST',
    body: formData,
  })
  if (!res.ok) throw new Error('Failed to upload configuration')
}

export async function updateFactoryConfig(id: string, name: string | null): Promise<void> {
  const res = await fetch(`${apiBaseUrl}/api/v1/factory-configs/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name }),
  })
  if (!res.ok) throw new Error('Failed to update configuration')
}

export async function deleteFactoryConfig(id: string): Promise<void> {
  const res = await fetch(`${apiBaseUrl}/api/v1/factory-configs/${id}`, {
    method: 'DELETE',
  })
  if (!res.ok) throw new Error('Failed to delete')
}
