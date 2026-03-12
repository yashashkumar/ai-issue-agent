import { useEffect, useState, useRef } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getAgent, getAgentLogs, getProject, cancelAgent } from '../api';
import { Agent, AgentLog, Project } from '../types';
import { Terminal, StopCircle, ArrowLeft, FileTerminal, RefreshCw } from 'lucide-react';

export function AgentDetail() {
  const { id } = useParams<{ id: string }>();
  const [agent, setAgent] = useState<Agent | null>(null);
  const [project, setProject] = useState<Project | null>(null);
  const [logs, setLogs] = useState<AgentLog[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);

  const loadData = async (agentId: string) => {
    try {
      const a = await getAgent(agentId);
      setAgent(a);
      if (a.project_id && !project) {
        getProject(a.project_id).then(setProject).catch(() => {});
      }
      const l = await getAgentLogs(agentId);
      setLogs(l || []);
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    if (id) {
      loadData(id);
      const interval = setInterval(() => {
        if (agent?.status === 'running' || agent?.status === 'pending' || !agent) {
          loadData(id);
        }
      }, 2000);
      return () => clearInterval(interval);
    }
  }, [id, agent?.status]);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  const handleCancel = async () => {
    if (!agent) return;
    if (confirm("Are you sure you want to cancel this agent?")) {
      await cancelAgent(agent.id);
      loadData(agent.id);
    }
  };

  if (!agent) return <div className="p-8 text-zinc-400">Loading agent...</div>;

  const isRunning = agent.status === 'running' || agent.status === 'pending';

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center gap-4">
        {project && (
          <Link to={`/projects/${project.id}`} className="text-zinc-500 hover:text-white transition">
            <ArrowLeft className="w-5 h-5" />
          </Link>
        )}
        <h1 className="text-2xl font-bold text-white flex items-center gap-2">
          {agent.name}
          <span className={`text-xs px-2 py-0.5 rounded uppercase font-bold tracking-wider ${
            agent.status === 'running' ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20 animate-pulse' :
            agent.status === 'error' ? 'bg-red-500/10 text-red-400 border border-red-500/20' :
            agent.status === 'finished' ? 'bg-green-500/10 text-green-400 border border-green-500/20' :
            'bg-yellow-500/10 text-yellow-500 border border-yellow-500/20'
          }`}>
            {agent.status}
          </span>
        </h1>
        <div className="flex-1" />
        {isRunning && (
          <button 
            onClick={handleCancel}
            className="flex items-center gap-2 px-4 py-2 bg-red-600/10 hover:bg-red-600/20 text-red-400 rounded-lg text-sm font-medium transition"
          >
            <StopCircle className="w-4 h-4" /> Cancel Process
          </button>
        )}
      </div>

      <div className="grid md:grid-cols-3 gap-6">
        <div className="md:col-span-1 space-y-6">
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6">
            <h3 className="text-sm font-semibold text-white mb-4 uppercase tracking-wider">Metadata</h3>
            <dl className="space-y-4 text-sm">
              <div>
                <dt className="text-zinc-500">Source</dt>
                <dd className="text-zinc-200 mt-1 capitalize">{agent.source.replace('_', ' ')}</dd>
              </div>
              <div>
                <dt className="text-zinc-500">Work Directory</dt>
                <dd className="text-zinc-200 mt-1 font-mono text-xs break-all bg-zinc-950 p-2 rounded">{agent.work_dir}</dd>
              </div>
              <div>
                <dt className="text-zinc-500">Process ID</dt>
                <dd className="text-zinc-200 mt-1">{agent.pid || 'N/A'}</dd>
              </div>
              {agent.github_issue_number && (
                <div>
                  <dt className="text-zinc-500">GitHub Issue</dt>
                  <dd className="text-zinc-200 mt-1">#{agent.github_issue_number}</dd>
                </div>
              )}
            </dl>
          </div>

          <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6">
            <h3 className="text-sm font-semibold text-white mb-4 uppercase tracking-wider flex items-center gap-2">
              <FileTerminal className="w-4 h-4" /> System Prompt
            </h3>
            <div className="bg-zinc-950 border border-zinc-800 rounded-lg p-4 max-h-64 overflow-y-auto font-mono text-xs text-zinc-300 whitespace-pre-wrap">
              {agent.system_prompt}
            </div>
          </div>
        </div>

        <div className="md:col-span-2">
          <div className="bg-zinc-900 border border-zinc-800 rounded-xl overflow-hidden flex flex-col h-[600px]">
            <div className="bg-zinc-800 border-b border-zinc-700 p-3 flex items-center justify-between">
              <div className="flex items-center gap-2 text-sm text-zinc-300 font-medium">
                <Terminal className="w-4 h-4 text-zinc-400" /> CLI Output
              </div>
              {isRunning && <RefreshCw className="w-4 h-4 text-zinc-500 animate-spin" />}
            </div>
            <div className="flex-1 bg-black p-4 overflow-y-auto font-mono text-sm leading-relaxed">
              {logs.length === 0 ? (
                <div className="text-zinc-600 italic">Waiting for standard output...</div>
              ) : (
                logs.map(log => (
                  <div key={log.id} className={`break-words ${log.stream === 'stderr' ? 'text-red-400' : 'text-zinc-300'}`}>
                    <span className="text-zinc-600 select-none mr-3 text-xs">
                      {new Date(log.created_at).toISOString().substring(11, 23)}
                    </span>
                    {log.content}
                  </div>
                ))
              )}
              <div ref={logsEndRef} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
