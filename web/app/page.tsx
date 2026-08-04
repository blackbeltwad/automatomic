'use client';

export default function Home() {
  const handleGitHubLogin = () => {
    const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    window.location.href = `${backendUrl}/api/v1/auth/github/login`;
  };

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-center items-center px-4">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(56,189,248,0.05)_0,transparent_100%)] pointer-events-none" />
      
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-xl p-8 shadow-2xl space-y-6 text-center z-10">
        <div className="inline-flex items-center space-x-2">
          <div className="w-8 h-8 bg-cyan-500 rounded-lg flex items-center justify-center font-mono font-bold text-slate-950">
          </div>
          <span className="text-2xl font-bold font-mono tracking-wider">AUTOMATOMIC</span>
        </div>

        <div className="space-y-2">
          <h1 className="text-xl font-semibold tracking-tight text-white">Enterprise CI/CD Engine</h1>
          <p className="text-sm text-slate-400">
            Internal Developer Platform control plane with isolated container execution and real-time streaming.
          </p>
        </div>

        <button
          onClick={handleGitHubLogin}
          className="w-full flex items-center justify-center space-x-3 bg-slate-800 hover:bg-slate-700 text-white font-medium py-3 px-4 rounded-lg border border-slate-700 transition-all duration-150 shadow-md hover:border-slate-600 active:scale-[0.98]"
        >
          <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
          </svg>
          <span>Sign in with GitHub</span>
        </button>

        <p className="text-xs text-slate-500">
          Local-first execution. Authentication secured via JWT.
        </p>
      </div>
    </main>
  );
}