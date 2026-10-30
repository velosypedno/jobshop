import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Box, Typography, Button, CircularProgress } from '@mui/material'
import type { FactoryConfigItem } from '../types'
import { ConfirmModal } from '../components/CommonComponents'
import { fetchFactoryConfigById, deleteFactoryConfig } from '../api/factoryConfigs'

export function DetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [config, setConfig] = useState<FactoryConfigItem | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showDeleteModal, setShowDeleteModal] = useState(false)

  useEffect(() => {
    if (!id) return
    fetchFactoryConfigById(id)
      .then((data) => setConfig(data))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [id])

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
      <Box
        sx={{
          width: '100%',
          background: '#ffffff',
          border: '1px solid #dadce0',
          borderRadius: '12px',
          borderTop: '8px solid #2e7d32',
          padding: '40px',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
        }}
      >
        <CircularProgress sx={{ color: '#2e7d32' }} />
      </Box>
    )
  }

  if (error || !config) {
    return (
      <Box
        sx={{
          width: '100%',
          background: '#ffffff',
          border: '1px solid #dadce0',
          borderRadius: '12px',
          borderTop: '8px solid #d93025',
          padding: '40px',
          textAlign: 'center',
          boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
        }}
      >
        <Typography sx={{ fontSize: '14px', color: '#d93025' }}>
          {error || 'Not found'}
        </Typography>
      </Box>
    )
  }

  return (
    <Box
      sx={{
        width: '100%',
        background: '#ffffff',
        border: '1px solid #dadce0',
        borderRadius: '12px',
        borderTop: '8px solid #2e7d32',
        padding: { xs: '24px', md: '32px 40px' },
        boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      {/* Header Actions */}
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '28px',
          borderBottom: '1px solid #e0e0e0',
          paddingBottom: '16px',
        }}
      >
        <Button
          onClick={() => navigate('/configs')}
          sx={{
            background: 'transparent',
            color: '#2e7d32',
            padding: 0,
            fontSize: '15px',
            fontWeight: 500,
            textTransform: 'none',
            '&:hover': {
              background: 'transparent',
              textDecoration: 'underline',
            },
          }}
        >
          ← Back to List
        </Button>
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

      {/* Title */}
      <Typography
        variant="h4"
        sx={{
          margin: '0 0 20px 0',
          fontSize: '28px',
          fontWeight: 400,
          color: '#202124',
        }}
      >
        {config.name || 'Unnamed'}
      </Typography>

      {/* Meta Info Box */}
      <Box
        sx={{
          display: 'flex',
          gap: { xs: '16px', md: '32px' },
          fontSize: '14px',
          color: '#5f6368',
          marginBottom: '24px',
          flexWrap: 'wrap',
          background: '#f8f9fa',
          padding: '16px 20px',
          borderRadius: '8px',
          border: '1px solid #e0e0e0',
        }}
      >
        <Typography component="span" sx={{ fontSize: '14px', color: '#5f6368' }}>
          <strong>ID:</strong> {config.id}
        </Typography>
        <Typography component="span" sx={{ fontSize: '14px', color: '#5f6368' }}>
          <strong>Version:</strong> {config.version}
        </Typography>
        <Typography component="span" sx={{ fontSize: '14px', color: '#5f6368' }}>
          <strong>Created At:</strong> {new Date(config.created_at).toLocaleString()}
        </Typography>
      </Box>

      {/* JSON Structure Subtitle */}
      <Typography
        variant="h6"
        sx={{
          margin: '24px 0 12px 0',
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
    </Box>
  )
}
