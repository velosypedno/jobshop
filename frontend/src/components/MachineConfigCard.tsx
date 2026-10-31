import { Box, Typography, Paper, Chip } from '@mui/material'
import PrecisionManufacturingIcon from '@mui/icons-material/PrecisionManufacturing'

interface MachineItem {
  type_id: number
  type_name: string
  count: number
}

export function MachineConfigCard({ machine }: { machine: MachineItem }) {
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
        '&:hover': { borderColor: '#2e7d32', boxShadow: '0 2px 8px rgba(46, 125, 50, 0.08)' },
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
        <Box
          sx={{
            width: 40,
            height: 40,
            borderRadius: '8px',
            backgroundColor: '#e8f5e9',
            color: '#2e7d32',
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
