import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getProject, getProjectAgents, spawnAgent, deleteProject } from '../api';
import { Project, Agent } from '../types';
import { Play, FileText, Trash2, Folder, Clock, CheckCircle2, XCircle, Loader2 } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export function ProjectDetail() {
  const { id } = useParams<{ id: string }>();
  const [project, setProject] = useState<Project | null>(null);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [prompt, setPrompt] = useState('');
  const [spawning, setSpawning] = useState(false);

  useEffect(() => {
    if (id) {
      loadData(id);
      const interval = setInterval(() => loadData(id), 3000);
      return () => clearInterval(interval);
    }
  }, [id]);

  const loadData = async (projectId: string) => {
    try {
      const [p, a] = await Promise.all([getProject(projectId), getProjectAgents(projectId)]);
      setProject(p);
      setAgents(a || []);
    } catch (e) {
      console.error(e);
    }
  };

  const handleSpawn = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id || !prompt) return;
    setSpawning(true);
    try {
      await spawnAgent({ project_id: id, system_prompt: prompt });
      setPrompt('');
      loadData(id);
    } catch (err) {
      alert("Failed to spawn agent");
    } finally {
      setSpawning(false);
    }
  };

  const handleDelete = async () => {
    if (!id) return;
    if (confirm("Are you sure you want to delete this project?")) {
      await deleteProject(id);
      window.location.href = '/projects';
    }
  };

  if (!project) return <div className="p-8 text-zinc-400">Loading project...</div>;

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <Folder className="w-8 h-8 text-indigo-400" />
            {project.name}
          </h1>
          <p className="mt-2 text-zinc-400 max-w-2xl">{project.description}</p>
          <div className="flex gap-4 mt-4 text-xs font-mono text-zinc-500 bg-zinc-900 border border-zinc-800 rounded-md p-2 w-max">
            <span>Root: {project.root_folder}</span>
            {project.github_owner && <span>GitHub: {project.github_owner}/{project.github_repo}</span>}
          </div>
        </div>
        <button 
          onClick={handleDelete}
          className="text-red-400 hover:text-red-300 hover:bg-red-400/10 p-2 rounded transition"
          title="Delete Project"
        >
          <Trash2 className="w-5 h-5" />
        </button>
      </div>

      <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-white mb-4">Spawn Custom Agent</h2>
        <form onSubmit={handleSpawn} className="flex gap-4 items-start">
          <textarea 
            required
            className="flex-1 bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-3 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition resize-none h-24 text-sm"
            placeholder="Describe what the agent should do..."
            value={prompt}
            onChange={e => setPrompt(e.target.value)}
          />
          <button 
            type="submit"
            disabled={spawning || !prompt}
            className="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition shadow-sm disabled:opacity-50 flex items-center gap-2"
          >
            {spawning ? <Loader2 className="w-4 h-4 animate-spin" /> : <Play className="w-4 h-4" />}
            Spawn
          </button>
        </form>
      </div>

      <div>
        <h2 className="text-xl font-semibold text-white mb-4">Project Agents</h2>
        <div className="grid gap-4">
          {agents.length === 0 ? (
            <div className="text-zinc-500 border border-dashed border-zinc-800 p-8 rounded-xl text-center">
              No agents have been spawned for this project yet.
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
                    <h3 className="font-medium text-zinc-200 group-hover:text-indigo-400 transition">{agent.name}</h3>
                    <p className="text-xs text-zinc-500 mt-1 flex items-center gap-2">
                       <FileText className="w-3 h-3" />
                       <span className="truncate max-w-sm">{agent.system_prompt}</span>
                    </p>
                  </div>
                </div>
                <div className="text-right text-xs font-mono text-zinc-500">
                  <div className="flex items-center justify-end gap-1 mb-1">
                    <Clock className="w-3 h-3" />
                    {formatDistanceToNow(new Date(agent.created_at), { addSuffix: true })}
                  </div>
                  <span className="capitalize">{agent.source.replace('_', ' ')}</span>
                </div>
              </Link>
            ))
          )}
        </div>
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
