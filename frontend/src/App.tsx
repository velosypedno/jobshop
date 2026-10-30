import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Box } from '@mui/material'
import { HomePage, ListPage, DetailPage } from './pages'

export default function App() {
  return (
    <BrowserRouter>
      <Box
        sx={{
          width: '100%',
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          padding: { xs: '16px', md: '32px 48px' },
          backgroundColor: '#f8f9fa', // Необов'язковий м'який фон для всієї сторінки
          boxSizing: 'border-box',
        }}
      >
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/configs" element={<ListPage />} />
          <Route path="/configs/:id" element={<DetailPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Box>
    </BrowserRouter>
  )
}
