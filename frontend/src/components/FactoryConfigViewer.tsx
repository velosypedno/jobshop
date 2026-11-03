import { useState, useMemo } from 'react'
import { Box, Typography, Button } from '@mui/material'
import CodeIcon from '@mui/icons-material/Code'
import VisibilityIcon from '@mui/icons-material/Visibility'
import { MachineConfigCard } from './MachineConfigCard'
import { JobTemplateCard } from './JobTemplateCard'
import { buildMachineColorMap } from '../utils/colorUtils'

export interface FactoryConfigData {
  machines: Array<{ type_id: number; type_name: string; count: number }>
  job_templates: Array<{
    id: number
    name: string
    operations: Array<{ id: number; name: string; machine_type: number; processing_time: string }>
    edges: Array<{ from_op_id: number; to_op_id: number }>
  }>
}

export function FactoryConfigViewer({ configData }: { configData: FactoryConfigData }) {
  const [showRawJson, setShowRawJson] = useState(false)

  // Генеруємо мапу кольорів на основі всіх доступних машин
  const machineColorMap = useMemo(() => {
    return buildMachineColorMap(configData?.machines || [])
  }, [configData?.machines])

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>

      {/* Кнопка перемикання між Візуальним режимом та сирим JSON */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h6" sx={{ fontSize: '18px', color: '#202124', fontWeight: 600 }}>
          Factory Infrastructure & Templates
        </Typography>
        <Button
          size="small"
          variant="outlined"
          onClick={() => setShowRawJson(!showRawJson)}
          startIcon={showRawJson ? <VisibilityIcon /> : <CodeIcon />}
          sx={{
            borderColor: '#dadce0',
            color: '#3c4043',
            textTransform: 'none',
            '&:hover': { borderColor: '#2e7d32', backgroundColor: '#f8f9fa' },
          }}
        >
          {showRawJson ? 'Show Visual View' : 'Show Raw JSON'}
        </Button>
      </Box>

      {showRawJson ? (
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
          {JSON.stringify(configData, null, 2)}
        </Box>
      ) : (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: '32px' }}>

          {/* СЕКЦІЯ 1: Список машин (Machines) */}
          <Box>
            <Typography sx={{ fontSize: '13px', fontWeight: 600, color: '#5f6368', textTransform: 'uppercase', mb: '12px', letterSpacing: '0.8px' }}>
              1. Available Machines ({configData.machines?.length || 0})
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' }, gap: '16px' }}>
              {configData.machines?.map((machine) => (
                <MachineConfigCard
                  key={machine.type_id}
                  machine={machine}
                  color={machineColorMap[machine.type_id]}
                />
              ))}
            </Box>
          </Box>

          {/* СЕКЦІЯ 2: Шаблони завдань (Job Templates) */}
          <Box>
            <Typography sx={{ fontSize: '13px', fontWeight: 600, color: '#5f6368', textTransform: 'uppercase', mb: '12px', letterSpacing: '0.8px' }}>
              2. Job Templates & DAG Workflows ({configData.job_templates?.length || 0})
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              {configData.job_templates?.map((template) => (
                <JobTemplateCard
                  key={template.id}
                  template={template}
                  machineColorMap={machineColorMap}
                />
              ))}
            </Box>
          </Box>

        </Box>
      )}

    </Box>
  )
}
