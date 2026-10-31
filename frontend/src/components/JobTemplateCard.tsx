import { useState } from 'react'
import { Box, Typography, Paper, Chip, Button, Collapse } from '@mui/material'
import AccountTreeIcon from '@mui/icons-material/AccountTree'
import TimerIcon from '@mui/icons-material/Timer'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import ExpandLessIcon from '@mui/icons-material/ExpandLess'

interface Operation {
  id: number
  name: string
  machine_type: number
  processing_time: string
}

interface Edge {
  from_op_id: number
  to_op_id: number
}

interface JobTemplate {
  id: number
  name: string
  operations: Operation[]
  edges: Edge[]
}

export function JobTemplateCard({ template }: { template: JobTemplate }) {
  const [expanded, setExpanded] = useState(false)

  return (
    <Paper
      elevation={0}
      sx={{
        border: '1px solid #dadce0',
        borderRadius: '8px',
        background: '#ffffff',
        overflow: 'hidden',
        transition: 'all 0.2s',
      }}
    >
      {/* Header Картки Темплейту */}
      <Box
        sx={{
          p: '20px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          backgroundColor: '#f8f9fa',
          borderBottom: expanded ? '1px solid #e0e0e0' : 'none',
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <Box
            sx={{
              width: 36,
              height: 36,
              borderRadius: '8px',
              backgroundColor: '#e8f0fe',
              color: '#1a73e8',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <AccountTreeIcon fontSize="small" />
          </Box>
          <Box>
            <Typography sx={{ fontSize: '16px', fontWeight: 600, color: '#202124' }}>
              {template.name}
            </Typography>
            <Typography sx={{ fontSize: '12px', color: '#5f6368' }}>
              ID: {template.id} • {template.operations.length} operations • {template.edges.length} dependencies
            </Typography>
          </Box>
        </Box>

        <Button
          size="small"
          onClick={() => setExpanded(!expanded)}
          endIcon={expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
          sx={{ textTransform: 'none', color: '#1a73e8', fontWeight: 500 }}
        >
          {expanded ? 'Hide Details' : 'View Operations'}
        </Button>
      </Box>

      {/* Список операцій та зв'язків (Розгортається) */}
      <Collapse in={expanded} timeout="auto" unmountOnExit>
        <Box sx={{ p: '20px', display: 'flex', flexDirection: 'column', gap: '16px' }}>

          {/* Операції */}
          <Box>
            <Typography sx={{ fontSize: '12px', fontWeight: 600, color: '#5f6368', textTransform: 'uppercase', mb: '8px', letterSpacing: '0.5px' }}>
              Operations Pipeline
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {template.operations.map((op) => (
                <Box
                  key={op.id}
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    p: '10px 14px',
                    borderRadius: '6px',
                    backgroundColor: '#fff',
                    border: '1px solid #e0e0e0',
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                    <Chip label={`#${op.id}`} size="small" sx={{ height: '22px', fontSize: '11px', fontWeight: 600 }} />
                    <Typography sx={{ fontSize: '14px', fontWeight: 500, color: '#202124' }}>
                      {op.name}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                    <Chip
                      label={`Machine Type: ${op.machine_type}`}
                      size="small"
                      variant="outlined"
                      sx={{ fontSize: '11px', borderColor: '#dadce0', color: '#5f6368' }}
                    />
                    <Chip
                      icon={<TimerIcon style={{ fontSize: '14px' }} />}
                      label={op.processing_time}
                      size="small"
                      sx={{ backgroundColor: '#e8f5e9', color: '#2e7d32', fontWeight: 500, fontSize: '12px' }}
                    />
                  </Box>
                </Box>
              ))}
            </Box>
          </Box>

          {/* Залежності (Edges / DAG) */}
          {template.edges.length > 0 && (
            <Box>
              <Typography sx={{ fontSize: '12px', fontWeight: 600, color: '#5f6368', textTransform: 'uppercase', mb: '6px', letterSpacing: '0.5px' }}>
                DAG Dependencies (Edges)
              </Typography>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
                {template.edges.map((edge, idx) => (
                  <Chip
                    key={idx}
                    label={`Op #${edge.from_op_id} → Op #${edge.to_op_id}`}
                    size="small"
                    sx={{ backgroundColor: '#f1f3f4', color: '#3c4043', fontSize: '11px', fontFamily: 'monospace' }}
                  />
                ))}
              </Box>
            </Box>
          )}

        </Box>
      </Collapse>
    </Paper>
  )
}
