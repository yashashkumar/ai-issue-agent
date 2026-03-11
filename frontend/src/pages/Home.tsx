import { Link } from 'react-router-dom';
import { Play, FolderGit2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { getProjects, getAgents } from '../api';
import { Project, Agent } from '../types';

export function Home() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([getProjects(), getAgents()])
      .then(([p, a]) => {
        setProjects(p || []);
        setAgents(a || []);
        setLoading(false);
      });
  }, []);

  const runningAgents = agents.filter(a => a.status === 'running');
  const finishedAgents = agents.filter(a => a.status === 'finished');
  const errorAgents = agents.filter(a => a.status === 'error');

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="px-6 py-12 bg-gradient-to-br from-zinc-800 to-zinc-900 shadow-2xl rounded-2xl border border-zinc-700 select-none text-center relative overflow-hidden">
        <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 -translate-y-1/2 translate-x-1/2"></div>
        <div className="absolute bottom-0 left-0 w-64 h-64 bg-purple-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 translate-y-1/2 -translate-x-1/2"></div>
        <div className="relative z-10">
          <h1 className="text-4xl font-extrabold tracking-tight text-white sm:text-5xl lg:text-6xl">
            Orchestrate Your <span className="text-indigo-400">AI Agents</span>
          </h1>
          <p className="mt-6 text-lg text-zinc-400 max-w-2xl mx-auto">
            Gateway and process manager for Gemini CLI powered agents. Monitor progress, handle GitHub webhooks, and automate your workflow safely.
          </p>
          <div className="mt-8 flex justify-center gap-4">
            <Link to="/projects/new" className="inline-flex items-center justify-center rounded-lg text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-11 px-6 bg-indigo-600 text-white hover:bg-indigo-700 shadow shadow-indigo-500/20 gap-2">
              <FolderGit2 className="w-4 h-4" />
              New Project
            </Link>
            <Link to="/agents" className="inline-flex items-center justify-center rounded-lg text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-11 px-6 bg-zinc-800 text-zinc-100 hover:bg-zinc-700 border border-zinc-700 gap-2">
              <Play className="w-4 h-4" />
              Quick Spawn
            </Link>
          </div>
        </div>
      </div>

      {!loading && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 shadow-sm flex flex-col justify-between">
            <h3 className="text-sm font-semibold text-zinc-400">Total Projects</h3>
            <p className="text-4xl font-bold text-white mt-4">{projects.length}</p>
          </div>
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 shadow-sm flex flex-col justify-between">
            <h3 className="text-sm font-semibold text-zinc-400">Running Agents</h3>
            <div className="flex items-center gap-3 mt-4">
              <span className="flex h-3 w-3 relative">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-blue-500"></span>
              </span>
              <p className="text-4xl font-bold text-white">{runningAgents.length}</p>
            </div>
          </div>
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 shadow-sm flex flex-col justify-between">
            <h3 className="text-sm font-semibold text-zinc-400">Completed</h3>
            <p className="text-4xl font-bold text-green-400 mt-4">{finishedAgents.length}</p>
          </div>
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 shadow-sm flex flex-col justify-between">
            <h3 className="text-sm font-semibold text-zinc-400">Errors</h3>
            <p className="text-4xl font-bold text-red-400 mt-4">{errorAgents.length}</p>
          </div>
        </div>
      )}
    </div>
  );
}
