import { useState } from 'react'
import { Box, Typography, Paper, Chip, Button, Collapse, ToggleButton, ToggleButtonGroup } from '@mui/material'
import AccountTreeIcon from '@mui/icons-material/AccountTree'
import ViewStreamIcon from '@mui/icons-material/ViewStream'
import GrainIcon from '@mui/icons-material/Grain'
import TimerIcon from '@mui/icons-material/Timer'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import ExpandLessIcon from '@mui/icons-material/ExpandLess'
import { DagNetworkGraph } from './DagNetworkGraph'

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

export function JobTemplateCard({
  template,
  machineColorMap
}: {
  template: JobTemplate;
  machineColorMap: Record<number, string>
}) {
  const [expanded, setExpanded] = useState(false)
  const [viewMode, setViewMode] = useState<'graph' | 'list'>('graph') // За замовчуванням показуємо граф

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
        <Box sx={{ p: '20px', display: 'flex', flexDirection: 'column', gap: '20px' }}>

          {/* Панель перемикання режимів (Graph Flow / List View) */}
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
            <Typography sx={{ fontSize: '12px', fontWeight: 600, color: '#5f6368', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
              Operations Pipeline Topology
            </Typography>

            <ToggleButtonGroup
              value={viewMode}
              exclusive
              onChange={(_, newMode) => newMode && setViewMode(newMode)}
              size="small"
              sx={{
                backgroundColor: '#f1f3f4',
                '& .MuiToggleButton-root': { textTransform: 'none', fontSize: '12px', px: '12px', py: '4px' }
              }}
            >
              <ToggleButton value="graph">
                <GrainIcon sx={{ fontSize: '16px', mr: '6px' }} /> Graph Flow
              </ToggleButton>
              <ToggleButton value="list">
                <ViewStreamIcon sx={{ fontSize: '16px', mr: '6px' }} /> List View
              </ToggleButton>
            </ToggleButtonGroup>
          </Box>

          {/* Відображення: або інтерактивний граф React Flow, або класичний список */}
          {viewMode === 'graph' ? (
            <DagNetworkGraph
              operations={template.operations}
              edges={template.edges}
              machineColorMap={machineColorMap}
            />
          ) : (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {template.operations.map((op) => {
                const machineColor = machineColorMap[op.machine_type] || '#5f6368'

                return (
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
                      borderLeft: `4px solid ${machineColor}`,
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
                        sx={{
                          fontSize: '11px',
                          backgroundColor: `${machineColor}15`,
                          color: machineColor,
                          fontWeight: 600,
                          border: `1px solid ${machineColor}40`
                        }}
                      />
                      <Chip
                        icon={<TimerIcon style={{ fontSize: '14px' }} />}
                        label={op.processing_time}
                        size="small"
                        sx={{ backgroundColor: '#e8f5e9', color: '#2e7d32', fontWeight: 500, fontSize: '12px' }}
                      />
                    </Box>
                  </Box>
                )
              })}
            </Box>
          )}

          {/* Ребра / Залежності у вигляді плашок (якщо потрібні длявідладки) */}
          {template.edges.length > 0 && viewMode === 'list' && (
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
