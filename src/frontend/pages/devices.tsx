import React from 'react'
import MainLayout from '../components/layout/MainLayout'
import DevicesDashboard from '../components/devices/DevicesDashboard'

export default function DevicesPage() {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3001'
  
  return (
    <MainLayout>
      <DevicesDashboard apiBaseUrl={apiBaseUrl} />
    </MainLayout>
  )
}
