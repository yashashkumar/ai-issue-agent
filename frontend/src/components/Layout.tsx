import { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Activity, Folder, LayoutDashboard } from 'lucide-react';
import { cn } from '../lib/utils';

export function Layout({ children }: { children: ReactNode }) {
  const location = useLocation();

  const navigation = [
    { name: 'Dashboard', href: '/', icon: LayoutDashboard },
    { name: 'Projects', href: '/projects', icon: Folder },
    { name: 'Agents', href: '/agents', icon: Activity },
  ];

  return (
    <div className="min-h-screen flex flex-col font-sans dark text-zinc-100 bg-zinc-950">
      <header className="border-b border-zinc-800 bg-zinc-900 shadow-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-xl shadow-lg">
              A
            </div>
            <Link to="/" className="text-xl font-bold tracking-tight text-white hover:text-indigo-400 transition-colors">
              AgentForge
            </Link>
          </div>
          <nav className="flex gap-1">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href || (item.href !== '/' && location.pathname.startsWith(item.href));
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={cn(
                    "flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                    isActive 
                      ? "bg-zinc-800 bg-opacity-70 text-indigo-400" 
                      : "text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800"
                  )}
                >
                  <item.icon className="w-4 h-4" />
                  {item.name}
                </Link>
              );
            })}
          </nav>
        </div>
      </header>

      <main className="flex-grow max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        {children}
      </main>

      <footer className="border-t border-zinc-800 py-6 text-center text-zinc-500 text-sm mt-auto">
        <p>AgentForge © {new Date().getFullYear()}. React + TailwindCSS + Go.</p>
      </footer>
    </div>
  );
}
