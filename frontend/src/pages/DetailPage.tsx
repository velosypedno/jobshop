import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Box, Typography, Button, CircularProgress, TextField } from '@mui/material'
import type { FactoryConfigItem } from '../types'
import { ConfirmModal, PageCard, PageHeader, MetaInfoBox } from '../components/CommonComponents'
import { fetchFactoryConfigById, deleteFactoryConfig, updateFactoryConfig } from '../api/factoryConfigs'

export function DetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [config, setConfig] = useState<FactoryConfigItem | null>(null)
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showDeleteModal, setShowDeleteModal] = useState(false)

  useEffect(() => {
    if (!id) return
    fetchFactoryConfigById(id)
      .then((data) => {
        setConfig(data)
        setName(data.name || '')
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [id])

  const handleSave = async () => {
    if (!id) return
    setSaving(true)
    try {
      await updateFactoryConfig(id, name.trim() || null)
      if (config) {
        setConfig({ ...config, name: name.trim() || null })
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Update failed')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!id) return
    try {
      await deleteFactoryConfig(id)
      navigate('/configs')
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Delete failed')
    }
  }

  if (loading) {
    return (
      <PageCard>
        <Box sx={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
          <CircularProgress sx={{ color: '#2e7d32' }} />
        </Box>
      </PageCard>
    )
  }

  if (error || !config) {
    return (
      <PageCard borderColor="#d93025">
        <Typography sx={{ fontSize: '14px', color: '#d93025', textAlign: 'center' }}>
          {error || 'Not found'}
        </Typography>
      </PageCard>
    )
  }

  const isNameChanged = name !== (config.name || '')

  return (
    <PageCard>
      {/* Header Actions */}
      <PageHeader
        title=""
        onBack={() => navigate('/configs')}
        backLabel="← Back to List"
        action={
          <Box sx={{ display: 'flex', gap: '12px' }}>
            {isNameChanged && (
              <Button
                variant="contained"
                onClick={handleSave}
                disabled={saving}
                sx={{
                  backgroundColor: '#2e7d32',
                  color: '#fff',
                  borderRadius: '6px',
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: 500,
                  textTransform: 'none',
                  boxShadow: 'none',
                  '&:hover': {
                    backgroundColor: '#1b5e20',
                  },
                }}
              >
                {saving ? 'Saving...' : 'Save Changes'}
              </Button>
            )}
            <Button
              variant="contained"
              onClick={() => setShowDeleteModal(true)}
              sx={{
                backgroundColor: '#d93025',
                color: '#fff',
                borderRadius: '6px',
                padding: '10px 20px',
                fontSize: '14px',
                fontWeight: 500,
                textTransform: 'none',
                boxShadow: 'none',
                '&:hover': {
                  backgroundColor: '#c5221f',
                },
              }}
            >
              Delete Config
            </Button>
          </Box>
        }
      />

      {/* Editable Name Field */}
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
        <Typography component="label" sx={{ fontSize: '13px', fontWeight: 500, color: '#5f6368' }}>
          Configuration Name:
        </Typography>
        <TextField
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter configuration name..."
          size="small"
          variant="outlined"
          sx={{
            maxWidth: '400px',
            background: '#fff',
            '& .MuiOutlinedInput-root': {
              borderRadius: '6px',
              fontSize: '18px',
              '& fieldset': { borderColor: '#dadce0' },
              '&:hover fieldset': { borderColor: '#dadce0' },
              '&.Mui-focused fieldset': { borderColor: '#2e7d32' },
            },
          }}
        />
      </Box>

      <MetaInfoBox
        items={[
          { label: 'ID', value: config.id },
          { label: 'Version', value: config.version },
          { label: 'Created At', value: new Date(config.created_at).toLocaleString() },
        ]}
      />

      {/* JSON Structure Subtitle */}
      <Typography
        variant="h6"
        sx={{
          fontSize: '18px',
          color: '#202124',
          fontWeight: 500,
        }}
      >
        Configuration JSON Structure:
      </Typography>

      {/* JSON Viewer */}
      <Box
        component="pre"
        sx={{
          background: '#1e1e1e',
          color: '#d4d4d4',
          border: '1px solid #333',
          padding: '24px',
          borderRadius: '8px',
          fontSize: '13px',
          fontFamily: 'monospace',
          maxHeight: '70vh',
          overflow: 'auto',
          margin: 0,
        }}
      >
        {JSON.stringify(config.config, null, 2)}
      </Box>

      {/* Delete Confirmation Modal */}
      {showDeleteModal && (
        <ConfirmModal
          title="Delete Configuration"
          message="Are you sure you want to delete this configuration?"
          onConfirm={handleDelete}
          onCancel={() => setShowDeleteModal(false)}
        />
      )}
    </PageCard>
  )
}
