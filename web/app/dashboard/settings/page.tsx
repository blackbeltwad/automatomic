'use client';

import { useEffect, useState } from 'react';
import { fetchWithAuth } from '@lib/api';

export default function SettingsPage() {
  const [statusMessage, setStatusMessage] = useState<string>('Testing API connection...');
  const [authenticated, setAuthenticated] = useState<boolean | null>(null);

  useEffect(() => {
    async function verifyBackendConnection() {
      try {
        const res = await fetchWithAuth('/api/v1/pipelines');
        if (res.ok) {
          const data = await res.json();
          setStatusMessage(data.message || 'Successfully authenticated with Go Control Plane');
          setAuthenticated(true);
        } else {
          setStatusMessage(`Authorization Error: Received status ${res.status}`);
          setAuthenticated(false);
        }
      } catch (err) {
        setStatusMessage('Failed to reach Go Control Plane at http://localhost:8080');
        setAuthenticated(false);
      }
    }

    verifyBackendConnection();
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-100">Project Settings</h1>
        <p className="text-sm text-slate-400">Manage pipeline configurations and control plane credentials.</p>
      </div>

      <div className="bg-slate-900 border border-slate-800 rounded-lg p-6 space-y-4">
        <h2 className="text-lg font-semibold text-slate-200">Control Plane Connection Status</h2>
        
        <div className="flex items-center space-x-3 p-4 bg-slate-950 rounded-md border border-slate-800">
          <div className={`w-3 h-3 rounded-full ${authenticated ? 'bg-emerald-500' : 'bg-amber-500 animate-pulse'}`} />
          <span className="font-mono text-sm text-slate-300">{statusMessage}</span>
        </div>
      </div>

      <div className="bg-slate-900 border border-slate-800 rounded-lg p-6 space-y-4">
        <h2 className="text-lg font-semibold text-slate-200">Environment Details</h2>
        <div className="grid grid-cols-2 gap-4 text-sm font-mono">
          <div className="p-3 bg-slate-950 rounded border border-slate-800">
            <span className="text-slate-500 block text-xs">Target Environment</span>
            <span className="text-slate-200">Local (Week 1 Monorepo)</span>
          </div>
          <div className="p-3 bg-slate-950 rounded border border-slate-800">
            <span className="text-slate-500 block text-xs">Control Plane URL</span>
            <span className="text-slate-200">http://localhost:8080</span>
          </div>
        </div>
      </div>
    </div>
  );
}