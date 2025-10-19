import React from 'react'
import MainLayout from '../components/layout/MainLayout'
import Dashboard from '../components/dashboard/Dashboard'

export default function Home() {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
  
  return (
    <MainLayout>
      <Dashboard apiBaseUrl={apiBaseUrl} />
    </MainLayout>
  )
}

