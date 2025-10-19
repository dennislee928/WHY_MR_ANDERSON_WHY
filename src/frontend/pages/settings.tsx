import React from 'react'
import MainLayout from '../components/layout/MainLayout'
import SettingsDashboard from '../components/settings/SettingsDashboard'

export default function SettingsPage() {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3001'
  
  return (
    <MainLayout>
      <SettingsDashboard apiBaseUrl={apiBaseUrl} />
    </MainLayout>
  )
}
