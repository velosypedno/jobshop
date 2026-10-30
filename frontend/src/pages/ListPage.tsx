import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Box,
  Typography,
  Button,
  TextField,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  CircularProgress,
} from '@mui/material'
import type { FactoryConfigItem } from '../types'
import { ConfirmModal } from '../components/CommonComponents'
import { fetchFactoryConfigs, uploadFactoryConfig, deleteFactoryConfig } from '../api/factoryConfigs'

export function ListPage() {
  const navigate = useNavigate()
  const [configs, setConfigs] = useState<FactoryConfigItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [uploading, setUploading] = useState(false)
  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null)

  const loadConfigs = async () => {
    setLoading(true)
    try {
      const items = await fetchFactoryConfigs()
      setConfigs(items)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load configs')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadConfigs()
  }, [])

  const handleUpload = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = e.currentTarget
    const fileInput = form.elements.namedItem('file') as HTMLInputElement
    const nameInput = form.elements.namedItem('name') as HTMLInputElement

    if (!fileInput.files || fileInput.files.length === 0) return

    const formData = new FormData()
    formData.append('file', fileInput.files[0])
    if (nameInput.value) {
      formData.append('name', nameInput.value)
    }

    setUploading(true)
    try {
      await uploadFactoryConfig(formData)
      form.reset()
      await loadConfigs()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setUploading(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTargetId) return
    try {
      await deleteFactoryConfig(deleteTargetId)
      setConfigs(configs.filter((c) => c.id !== deleteTargetId))
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Delete failed')
    } finally {
      setDeleteTargetId(null)
    }
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
      {/* Page Header */}
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
          onClick={() => navigate('/')}
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
          ← Back to Home
        </Button>
        <Typography
          variant="h4"
          sx={{
            margin: 0,
            fontSize: '28px',
            fontWeight: 400,
            color: '#202124',
          }}
        >
          Factory Configurations
        </Typography>
      </Box>

      {/* Upload Form */}
      <Box
        component="form"
        onSubmit={handleUpload}
        sx={{
          background: '#f8f9fa',
          padding: '24px',
          borderRadius: '8px',
          marginBottom: '32px',
          border: '1px solid #e0e0e0',
          display: 'flex',
          gap: '20px',
          alignItems: 'flex-end',
          flexWrap: 'wrap',
        }}
      >
        <Box sx={{ flex: 1, minWidth: '250px', display: 'flex', flexDirection: 'column', gap: '6px' }}>
          <Typography component="label" sx={{ fontSize: '13px', fontWeight: 500, color: '#5f6368' }}>
            Name (optional):
          </Typography>
          <TextField
            name="name"
            placeholder="Enter name..."
            size="small"
            variant="outlined"
            sx={{
              background: '#fff',
              '& .MuiOutlinedInput-root': {
                borderRadius: '6px',
                '& fieldset': { borderColor: '#dadce0' },
                '&:hover fieldset': { borderColor: '#dadce0' },
                '&.Mui-focused fieldset': { borderColor: '#2e7d32' },
              },
            }}
          />
        </Box>

        <Box sx={{ flex: 1, minWidth: '250px', display: 'flex', flexDirection: 'column', gap: '6px' }}>
          <Typography component="label" sx={{ fontSize: '13px', fontWeight: 500, color: '#5f6368' }}>
            JSON Config File:
          </Typography>
          <Button
            component="label"
            variant="outlined"
            sx={{
              justifyContent: 'flex-start',
              borderColor: '#dadce0',
              color: '#3c4043',
              background: '#fff',
              borderRadius: '6px',
              padding: '8px 12px',
              textTransform: 'none',
              fontSize: '14px',
              '&:hover': {
                borderColor: '#2e7d32',
                background: '#fff',
              },
            }}
          >
            Choose File...
            <input type="file" name="file" accept=".json" required hidden />
          </Button>
        </Box>

        <Button
          type="submit"
          variant="contained"
          disabled={uploading}
          sx={{
            backgroundColor: '#2e7d32',
            '&:hover': { backgroundColor: '#1b5e20' },
            textTransform: 'none',
            borderRadius: '6px',
            padding: '10px 20px',
            fontSize: '14px',
            fontWeight: 500,
            boxShadow: 'none',
            height: '40px',
          }}
        >
          {uploading ? 'Uploading...' : 'Upload Config'}
        </Button>
      </Box>

      {/* List / Table Container */}
      <Box sx={{ overflowX: 'auto' }}>
        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
            <CircularProgress sx={{ color: '#2e7d32' }} />
          </Box>
        )}

        {error && (
          <Typography sx={{ fontSize: '14px', color: '#d93025', textAlign: 'center', padding: '40px 0' }}>
            {error}
          </Typography>
        )}

        {!loading && !error && configs.length === 0 && (
          <Typography sx={{ fontSize: '14px', color: '#5f6368', textAlign: 'center', padding: '40px 0' }}>
            No configurations found.
          </Typography>
        )}

        {!loading && configs.length > 0 && (
          <TableContainer>
            <Table sx={{ width: '100%', borderCollapse: 'collapse' }}>
              <TableHead>
                <TableRow>
                  <TableCell sx={{ color: '#5f6368', fontWeight: 600, background: '#f8f9fa', borderBottom: '1px solid #e0e0e0', padding: '16px 20px' }}>
                    Name
                  </TableCell>
                  <TableCell sx={{ color: '#5f6368', fontWeight: 600, background: '#f8f9fa', borderBottom: '1px solid #e0e0e0', padding: '16px 20px' }}>
                    Version
                  </TableCell>
                  <TableCell sx={{ color: '#5f6368', fontWeight: 600, background: '#f8f9fa', borderBottom: '1px solid #e0e0e0', padding: '16px 20px' }}>
                    Created At
                  </TableCell>
                  <TableCell sx={{ color: '#5f6368', fontWeight: 600, background: '#f8f9fa', borderBottom: '1px solid #e0e0e0', padding: '16px 20px' }}>
                    Actions
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {configs.map((cfg) => (
                  <TableRow
                    key={cfg.id}
                    onClick={() => navigate(`/configs/${cfg.id}`)}
                    sx={{
                      cursor: 'pointer',
                      transition: 'background 0.15s',
                      '&:hover': { backgroundColor: '#f1f3f4' },
                    }}
                  >
                    <TableCell sx={{ padding: '16px 20px', borderBottom: '1px solid #e0e0e0', fontSize: '14px', color: '#202124' }}>
                      {cfg.name || 'Unnamed'}
                    </TableCell>
                    <TableCell sx={{ padding: '16px 20px', borderBottom: '1px solid #e0e0e0' }}>
                      <Chip
                        label={cfg.version}
                        size="small"
                        sx={{
                          background: '#e8f5e9',
                          color: '#2e7d32',
                          borderRadius: '4px',
                          fontSize: '12px',
                          fontWeight: 500,
                        }}
                      />
                    </TableCell>
                    <TableCell sx={{ padding: '16px 20px', borderBottom: '1px solid #e0e0e0', fontSize: '14px', color: '#5f6368' }}>
                      {new Date(cfg.created_at).toLocaleString()}
                    </TableCell>
                    <TableCell sx={{ padding: '16px 20px', borderBottom: '1px solid #e0e0e0' }}>
                      <Button
                        variant="contained"
                        size="small"
                        onClick={(e) => {
                          e.stopPropagation()
                          setDeleteTargetId(cfg.id)
                        }}
                        sx={{
                          backgroundColor: '#fce8e6',
                          color: '#c5221f',
                          padding: '6px 12px',
                          fontSize: '13px',
                          borderRadius: '4px',
                          boxShadow: 'none',
                          textTransform: 'none',
                          fontWeight: 500,
                          '&:hover': {
                            backgroundColor: '#fad2cf',
                            boxShadow: 'none',
                          },
                        }}
                      >
                        Delete
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Box>

      {/* Delete Confirmation Modal */}
      {deleteTargetId && (
        <ConfirmModal
          title="Confirm Deletion"
          message="Are you sure you want to delete this configuration? This action cannot be undone."
          onConfirm={handleDelete}
          onCancel={() => setDeleteTargetId(null)}
        />
      )}
    </Box>
  )
}
