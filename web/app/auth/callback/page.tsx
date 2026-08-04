'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export default function AuthCallback() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get('token');
    if (token) {
      localStorage.setItem('automatomic_jwt', token);
      router.push('/dashboard/settings');
    } else {
      router.push('/?error=missing_token');
    }
  }, [router, searchParams]);

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center text-slate-200">
      <div className="flex items-center space-x-3">
        <div className="w-5 h-5 border-2 border-cyan-500 border-t-transparent rounded-full animate-spin" />
        <span className="font-mono text-sm">Finalizing authentication with Automatomic Control Plane...</span>
      </div>
    </div>
  );
}