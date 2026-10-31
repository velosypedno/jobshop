import React from 'react'
import { Box, Typography, Button } from '@mui/material'

interface ConfirmModalProps {
  title: string
  message: string
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmModal({ title, message, onConfirm, onCancel }: ConfirmModalProps) {
  return (
    <Box
      className="modal-overlay"
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: 'rgba(0, 0, 0, 0.5)',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 1000,
      }}
    >
      <Box
        className="modal-card"
        sx={{
          background: '#fff',
          padding: '32px',
          borderRadius: '12px',
          width: '100%',
          maxWidth: '450px',
          boxShadow: '0 8px 24px rgba(0,0,0,0.2)',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <Typography variant="h6" className="modal-title" sx={{ margin: '0 0 12px', fontSize: '20px', color: '#202124' }}>
          {title}
        </Typography>
        <Typography variant="body2" className="modal-message" sx={{ margin: '0 0 24px', fontSize: '14px', color: '#5f6368', lineHeight: 1.5 }}>
          {message}
        </Typography>
        <Box className="modal-actions" sx={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
          <Button
            variant="outlined"
            onClick={onCancel}
            sx={{
              textTransform: 'none',
              borderColor: '#dadce0',
              color: '#3c4043',
              borderRadius: '6px',
              '&:hover': { borderColor: '#b0b3b8', background: '#f8f9fa' },
            }}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            onClick={onConfirm}
            sx={{
              textTransform: 'none',
              backgroundColor: '#d93025',
              boxShadow: 'none',
              borderRadius: '6px',
              '&:hover': { backgroundColor: '#c5221f', boxShadow: 'none' },
            }}
          >
            Confirm
          </Button>
        </Box>
      </Box>
    </Box>
  )
}

// --- НОВІ СПІЛЬНІ КОМПОНЕНТИ ---

interface PageCardProps {
  children: React.ReactNode
  borderColor?: string // Наприклад, зелений для головної/деталей, чи червоний для помилок
}

export function PageCard({ children, borderColor = '#2e7d32' }: PageCardProps) {
  return (
    <Box
      sx={{
        width: '100%',
        background: '#ffffff',
        border: '1px solid #dadce0',
        borderRadius: '12px',
        borderTop: `8px solid ${borderColor}`,
        padding: { xs: '24px', md: '32px 40px' },
        boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
        display: 'flex',
        flexDirection: 'column',
        gap: '28px',
      }}
    >
      {children}
    </Box>
  )
}

interface PageHeaderProps {
  title: string
  onBack: () => void
  backLabel?: string
  action?: React.ReactNode
}

export function PageHeader({ title, onBack, backLabel = '← Back', action }: PageHeaderProps) {
  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '4px',
        borderBottom: '1px solid #e0e0e0',
        paddingBottom: '16px',
      }}
    >
      <Button
        onClick={onBack}
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
        {backLabel}
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
        {title}
      </Typography>

      {action || <Box />}
    </Box>
  )
}


interface MetaItem {
  label: string
  value: React.ReactNode
}

interface MetaInfoBoxProps {
  items: MetaItem[]
}

export function MetaInfoBox({ items }: MetaInfoBoxProps) {
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        gap: '8px',
        fontSize: '14px',
        color: '#5f6368',
        background: '#f8f9fa',
        padding: '16px 20px',
        borderRadius: '8px',
        border: '1px solid #e0e0e0',
      }}
    >
      {items.map((item, index) => (
        <Typography key={index} component="div" sx={{ fontSize: '14px', color: '#5f6368' }}>
          <strong>{item.label}:</strong> {item.value}
        </Typography>
      ))}
    </Box>
  )
}
