'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function Sidebar() {
  const pathname = usePathname();

  const navItems = [
    { label: 'Overview', href: '/dashboard' },
    { label: 'Pipelines', href: '/dashboard/pipelines' },
    { label: 'Project Settings', href: '/dashboard/settings' },
  ];

  return (
    <aside className="w-64 bg-slate-900 border-r border-slate-800 min-h-screen flex flex-col p-4">
      <div className="flex items-center space-x-2 px-2 py-4 border-b border-slate-800 mb-6">
        <div className="w-6 h-6 bg-cyan-500 rounded flex items-center justify-center font-mono text-xs font-bold text-slate-950">
          A
        </div>
        <span className="font-mono font-bold tracking-wider text-slate-100">AUTOMATOMIC</span>
      </div>

      <nav className="flex-1 space-y-1">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`block px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                isActive
                  ? 'bg-slate-800 text-cyan-400 border border-slate-700'
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
              }`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>

      <div className="border-t border-slate-800 pt-4 px-2 text-xs text-slate-500 font-mono">
        Status: <span className="text-emerald-400">Control Plane Connected</span>
      </div>
    </aside>
  );
}