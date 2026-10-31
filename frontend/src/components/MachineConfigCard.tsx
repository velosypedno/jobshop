import { Box, Typography, Paper, Chip } from '@mui/material'
import PrecisionManufacturingIcon from '@mui/icons-material/PrecisionManufacturing'
import { getMachineBackgroundColor } from '../utils/colorUtils'

interface MachineItem {
  type_id: number
  type_name: string
  count: number
}

export function MachineConfigCard({ machine, color }: { machine: MachineItem; color: string }) {
  const bgColor = getMachineBackgroundColor(machine.type_id);

  return (
    <Paper
      elevation={0}
      sx={{
        p: '16px',
        border: '1px solid #dadce0',
        borderRadius: '8px',
        background: '#ffffff',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        transition: 'all 0.2s',
        '&:hover': { borderColor: color, boxShadow: `0 2px 8px rgba(0, 0, 0, 0.06)` },
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
        <Box
          sx={{
            width: 40,
            height: 40,
            borderRadius: '8px',
            backgroundColor: bgColor,
            color: color,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <PrecisionManufacturingIcon fontSize="small" />
        </Box>
        <Box>
          <Typography sx={{ fontSize: '15px', fontWeight: 600, color: '#202124' }}>
            {machine.type_name}
          </Typography>
          <Typography sx={{ fontSize: '12px', color: '#5f6368', fontFamily: 'monospace' }}>
            ID: {machine.type_id}
          </Typography>
        </Box>
      </Box>
      <Chip
        label={`${machine.count} units`}
        size="small"
        sx={{
          backgroundColor: '#f1f3f4',
          color: '#3c4043',
          fontWeight: 500,
          borderRadius: '4px',
        }}
      />
    </Paper>
  )
}
