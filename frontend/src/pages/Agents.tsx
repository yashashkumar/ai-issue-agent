import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getAgents, getProjects } from '../api';
import { Agent, Project } from '../types';
import { Clock, CheckCircle2, XCircle, Loader2 } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export function Agents() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [projects, setProjects] = useState<Record<string, Project>>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([getAgents(), getProjects()]).then(([a, p]) => {
      setAgents(a || []);
      setProjects((p || []).reduce((acc, curr) => ({ ...acc, [curr.id]: curr }), {}));
      setLoading(false);
    });
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">All Agents</h1>
      </div>

      <div className="grid gap-4">
        {loading ? (
          <p className="text-zinc-400">Loading agents...</p>
        ) : agents.length === 0 ? (
          <div className="col-span-full py-12 text-center text-zinc-500 border border-dashed border-zinc-700 rounded-xl">
            No agents found across any projects.
          </div>
        ) : (
          agents.map(agent => (
            <Link 
              key={agent.id} 
              to={`/agents/${agent.id}`}
              className="flex items-center justify-between p-4 bg-zinc-900 border border-zinc-800 rounded-lg hover:border-zinc-700 transition group"
            >
              <div className="flex items-center gap-4">
                <StatusIcon status={agent.status} />
                <div>
                  <h3 className="font-medium text-zinc-200 group-hover:text-indigo-400 transition flex items-center gap-2">
                    {agent.name}
                    {agent.project_id && projects[agent.project_id] && (
                      <span className="text-xs bg-zinc-800 px-2 py-0.5 rounded text-zinc-400">
                        {projects[agent.project_id].name}
                      </span>
                    )}
                  </h3>
                  <p className="text-xs text-zinc-500 mt-1 truncate max-w-xl font-mono">
                     {agent.system_prompt}
                  </p>
                </div>
              </div>
              <div className="text-right text-xs font-mono text-zinc-500 flex flex-col items-end">
                <div className="flex items-center gap-1 mb-1">
                  <Clock className="w-3 h-3" />
                  {formatDistanceToNow(new Date(agent.created_at), { addSuffix: true })}
                </div>
                <span className="capitalize px-2 py-0.5 rounded bg-zinc-800/50 border border-zinc-700 shrink-0">
                  {agent.source.replace('_', ' ')}
                </span>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}

function StatusIcon({ status }: { status: string }) {
  switch(status) {
    case 'running':
      return <Loader2 className="w-5 h-5 text-blue-400 animate-spin" />;
    case 'finished':
      return <CheckCircle2 className="w-5 h-5 text-green-400" />;
    case 'error':
      return <XCircle className="w-5 h-5 text-red-400" />;
    default:
      return <Clock className="w-5 h-5 text-yellow-400" />;
  }
}
