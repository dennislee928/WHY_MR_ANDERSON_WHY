import React from 'react'
import MainLayout from '../components/layout/MainLayout'
import SecurityDashboard from '../components/security/SecurityDashboard'

export default function SecurityPage() {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3001'
  
  return (
    <MainLayout>
      <SecurityDashboard apiBaseUrl={apiBaseUrl} />
    </MainLayout>
  )
}
