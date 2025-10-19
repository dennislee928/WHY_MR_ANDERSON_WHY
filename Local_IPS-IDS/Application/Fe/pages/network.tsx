import React from 'react'
import MainLayout from '../components/layout/MainLayout'
import NetworkDashboard from '../components/network/NetworkDashboard'

export default function NetworkPage() {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3001'
  
  return (
    <MainLayout>
      <NetworkDashboard apiBaseUrl={apiBaseUrl} />
    </MainLayout>
  )
}
