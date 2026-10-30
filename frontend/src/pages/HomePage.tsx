import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Box, Typography, Button, Divider, Chip } from '@mui/material'
import { pingApi } from '../api/factoryConfigs'

export function HomePage() {
  const [status, setStatus] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/+$/, '')

  const handlePing = async () => {
    setLoading(true)
    setStatus('Pinging...')
    try {
      const data = await pingApi()
      setStatus(`OK: ${JSON.stringify(data)}`)
    } catch (err: unknown) {
      setStatus(`Error: ${err instanceof Error ? err.message : 'failed to connect'}`)
    } finally {
      setLoading(false)
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
        gap: '28px',
      }}
    >
      {/* Top Header & System Status */}
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 2,
        }}
      >
        <Box>
          <Typography variant="h4" sx={{ fontSize: '32px', fontWeight: '500', color: '#202124', margin: 0 }}>
            JobShop Control Center
          </Typography>
          <Typography variant="body2" sx={{ mt: '4px', color: '#5f6368', fontSize: '14px' }}>
            Advanced JSSP Optimization & Factory Configuration Management System
          </Typography>
        </Box>
        <Box sx={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
          <Typography component="span" sx={{ fontSize: '13px', color: '#5f6368', fontFamily: 'monospace' }}>
            {apiBaseUrl}/api/ping
          </Typography>
          <Button
            variant="contained"
            onClick={handlePing}
            disabled={loading}
            sx={{
              backgroundColor: '#2e7d32',
              '&:hover': { backgroundColor: '#1b5e20' },
              textTransform: 'none',
              padding: '8px 16px',
              borderRadius: '6px',
              fontSize: '14px',
              fontWeight: 500,
              boxShadow: 'none',
            }}
          >
            {loading ? 'Pinging...' : 'System Ping'}
          </Button>
        </Box>
      </Box>

      {status && (
        <Box
          sx={{
            padding: '16px',
            background: '#f8f9fa',
            borderRadius: '6px',
            fontSize: '13px',
            fontFamily: 'monospace',
            color: '#3c4043',
            wordBreak: 'break-all',
            border: '1px solid #e0e0e0',
          }}
        >
          {status}
        </Box>
      )}

      <Divider sx={{ borderTop: '1px solid #e0e0e0', m: 0 }} />

      {/* --- PIPELINE ARCHITECTURE LAYOUT (Flexbox Grid) --- */}
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>

        {/* SECTION 1: FOUNDATION (Independent Entities) */}
        <Box>
          <Typography
            variant="subtitle2"
            sx={{
              mb: '12px',
              fontSize: '13px',
              textTransform: 'uppercase',
              letterSpacing: '0.8px',
              color: '#5f6368',
              fontWeight: 600,
            }}
          >
            1. Core Foundations (Independent)
          </Typography>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' }, gap: '20px' }}>

            {/* Configurations (Active) */}
            <Box
              sx={{
                background: '#f8f9fa',
                border: '1px solid #2e7d32',
                borderRadius: '8px',
                padding: '24px',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between',
                boxShadow: '0 2px 8px rgba(46, 125, 50, 0.1)',
              }}
            >
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: '8px' }}>
                  <Typography variant="h6" sx={{ fontSize: '18px', color: '#202124', fontWeight: 600, margin: 0 }}>
                    Configurations
                  </Typography>
                  <Chip
                    label="Active"
                    size="small"
                    sx={{ bgcolor: '#e8f5e9', color: '#2e7d32', fontWeight: 500, borderRadius: '4px', fontSize: '12px' }}
                  />
                </Box>
                <Typography variant="body2" sx={{ mb: '20px', color: '#5f6368', lineHeight: '1.4', fontSize: '14px' }}>
                  Manage factory layouts, DAG templates, upload JSON configurations, and inspect machine parameters.
                </Typography>
              </Box>
              <Button
                component={Link}
                to="/configs"
                variant="contained"
                sx={{
                  backgroundColor: '#2e7d32',
                  '&:hover': { backgroundColor: '#1b5e20' },
                  textTransform: 'none',
                  textAlign: 'center',
                  display: 'block',
                  width: '100%',
                  borderRadius: '6px',
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: 500,
                  boxShadow: 'none',
                }}
              >
                Open Configurations →
              </Button>
            </Box>

            {/* Strategies (Fake) */}
            <Box
              sx={{
                background: '#ffffff',
                border: '1px solid #dadce0',
                borderRadius: '8px',
                padding: '24px',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between',
              }}
            >
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: '8px' }}>
                  <Typography variant="h6" sx={{ fontSize: '18px', color: '#202124', fontWeight: 600, margin: 0 }}>
                    Strategies
                  </Typography>
                  <Chip
                    label="Coming Soon"
                    size="small"
                    sx={{
                      background: '#f1f3f4',
                      color: '#5f6368',
                      borderRadius: '4px',
                      fontSize: '11px',
                      fontWeight: 500,
                      textTransform: 'uppercase',
                      letterSpacing: '0.5px',
                    }}
                  />
                </Box>
                <Typography variant="body2" sx={{ mb: '20px', color: '#5f6368', lineHeight: '1.4', fontSize: '14px' }}>
                  Configure metaheuristics, genetic algorithm parameters, mutation rates, and DAG optimization rules.
                </Typography>
              </Box>
              <Button
                variant="contained"
                disabled
                sx={{
                  width: '100%',
                  textTransform: 'none',
                  cursor: 'not-allowed',
                  backgroundColor: '#f1f3f4',
                  color: '#80868b',
                  boxShadow: 'none',
                  borderRadius: '6px',
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: 500,
                  '&.Mui-disabled': {
                    backgroundColor: '#f1f3f4',
                    color: '#80868b',
                  },
                }}
              >
                Configure Strategies
              </Button>
            </Box>

          </Box>
        </Box>

        {/* SECTION 2: EXECUTION LAYER (Dependent Entities) */}
        <Box>
          <Typography
            variant="subtitle2"
            sx={{
              mb: '12px',
              fontSize: '13px',
              textTransform: 'uppercase',
              letterSpacing: '0.8px',
              color: '#5f6368',
              fontWeight: 600,
            }}
          >
            2. Optimization Pipeline (Derived & Synthesis)
          </Typography>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' }, gap: '20px' }}>

            {/* Problems */}
            <Box
              sx={{
                background: '#ffffff',
                border: '1px solid #dadce0',
                borderRadius: '8px',
                padding: '24px',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between',
              }}
            >
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: '8px' }}>
                  <Typography variant="h6" sx={{ fontSize: '18px', color: '#202124', fontWeight: 600, margin: 0 }}>
                    Problems
                  </Typography>
                  <Chip
                    label="Depends on: Config"
                    size="small"
                    sx={{
                      background: '#e8f0fe',
                      color: '#1a73e8',
                      borderRadius: '4px',
                      fontSize: '11px',
                      fontWeight: 500,
                      letterSpacing: '0.3px',
                    }}
                  />
                </Box>
                <Typography variant="body2" sx={{ mb: '20px', color: '#5f6368', lineHeight: '1.4', fontSize: '14px' }}>
                  Define multi-objective JSSP problem instances, constraints, time windows, and resource limitations.
                </Typography>
              </Box>
              <Button
                variant="contained"
                disabled
                sx={{
                  width: '100%',
                  textTransform: 'none',
                  cursor: 'not-allowed',
                  backgroundColor: '#f1f3f4',
                  color: '#80868b',
                  boxShadow: 'none',
                  borderRadius: '6px',
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: 500,
                  '&.Mui-disabled': {
                    backgroundColor: '#f1f3f4',
                    color: '#80868b',
                  },
                }}
              >
                View Problems
              </Button>
            </Box>

            {/* Solutions */}
            <Box
              sx={{
                background: '#ffffff',
                border: '1px solid #dadce0',
                borderRadius: '8px',
                padding: '24px',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between',
              }}
            >
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: '8px' }}>
                  <Typography variant="h6" sx={{ fontSize: '18px', color: '#202124', fontWeight: 600, margin: 0 }}>
                    Solutions
                  </Typography>
                  <Chip
                    label="Problem + Strategy"
                    size="small"
                    sx={{
                      background: '#e8f0fe',
                      color: '#1a73e8',
                      borderRadius: '4px',
                      fontSize: '11px',
                      fontWeight: 500,
                      letterSpacing: '0.3px',
                    }}
                  />
                </Box>
                <Typography variant="body2" sx={{ mb: '20px', color: '#5f6368', lineHeight: '1.4', fontSize: '14px' }}>
                  Inspect computed manufacturing schedules, sequence graphs, and strongly connected operation flows.
                </Typography>
              </Box>
              <Button
                variant="contained"
                disabled
                sx={{
                  width: '100%',
                  textTransform: 'none',
                  cursor: 'not-allowed',
                  backgroundColor: '#f1f3f4',
                  color: '#80868b',
                  boxShadow: 'none',
                  borderRadius: '6px',
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: 500,
                  '&.Mui-disabled': {
                    backgroundColor: '#f1f3f4',
                    color: '#80868b',
                  },
                }}
              >
                View Solutions
              </Button>
            </Box>

          </Box>
        </Box>

        {/* SECTION 3: ANALYTICS LAYER */}
        <Box>
          <Typography
            variant="subtitle2"
            sx={{
              mb: '12px',
              fontSize: '13px',
              textTransform: 'uppercase',
              letterSpacing: '0.8px',
              color: '#5f6368',
              fontWeight: 600,
            }}
          >
            3. Monitoring & Analytics
          </Typography>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' }, gap: '20px' }}>

            {/* Dashboards */}
            <Box
              sx={{
                background: '#ffffff',
                border: '1px solid #dadce0',
                borderRadius: '8px',
                padding: '24px',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between',
              }}
            >
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: '8px' }}>
                  <Typography variant="h6" sx={{ fontSize: '18px', color: '#202124', fontWeight: 600, margin: 0 }}>
                    Dashboards
                  </Typography>
                  <Chip
                    label="Coming Soon"
                    size="small"
                    sx={{
                      background: '#f1f3f4',
                      color: '#5f6368',
                      borderRadius: '4px',
                      fontSize: '11px',
                      fontWeight: 500,
                      textTransform: 'uppercase',
                      letterSpacing: '0.5px',
                    }}
                  />
                </Box>
                <Typography variant="body2" sx={{ mb: '20px', color: '#5f6368', lineHeight: '1.4', fontSize: '14px' }}>
                  Real-time analytics, performance metrics, convergence charts, and manufactory throughput stats.
                </Typography>
              </Box>
              <Button
                variant="contained"
                disabled
                sx={{
                  width: '100%',
                  textTransform: 'none',
                  cursor: 'not-allowed',
                  backgroundColor: '#f1f3f4',
                  color: '#80868b',
                  boxShadow: 'none',
                  borderRadius: '6px',
                  padding: '10px 20px',
                  fontSize: '14px',
                  fontWeight: 500,
                  '&.Mui-disabled': {
                    backgroundColor: '#f1f3f4',
                    color: '#80868b',
                  },
                }}
              >
                Open Dashboards
              </Button>
            </Box>

          </Box>
        </Box>

      </Box>
    </Box>
  )
}
