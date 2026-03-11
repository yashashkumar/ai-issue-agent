import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createProject } from '../api';
import { Project } from '../types';

export function CreateProject() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState<Partial<Project>>({
    name: '',
    description: '',
    root_folder: '',
    allowed_emails: [],
    github_owner: '',
    github_repo: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const p = await createProject(formData);
      navigate(`/projects/${p.id}`);
    } catch (err) {
      console.error(err);
      alert('Failed to create project');
    }
  };

  return (
    <div className="max-w-3xl mx-auto space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-white mb-2">Create New Project</h1>
        <p className="text-zinc-400">Configure a workspace for your agents to operate in.</p>
      </div>

      <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 shadow-sm">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-zinc-300 mb-1">Project Name *</label>
              <input 
                required
                className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition"
                value={formData.name}
                onChange={e => setFormData({ ...formData, name: e.target.value })}
                placeholder="e.g. Acme Web App"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-zinc-300 mb-1">Description</label>
              <textarea 
                className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition h-24"
                value={formData.description}
                onChange={e => setFormData({ ...formData, description: e.target.value })}
                placeholder="What is this project?"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-zinc-300 mb-1">Root Folder *</label>
              <input 
                required
                className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition"
                value={formData.root_folder}
                onChange={e => setFormData({ ...formData, root_folder: e.target.value })}
                placeholder="/absolute/path/to/repo"
              />
              <p className="text-xs text-zinc-500 mt-1">Agents will execute operations within this underlying OS directory.</p>
            </div>

            <div className="pt-4 border-t border-zinc-800">
              <h3 className="text-lg font-medium text-white mb-4">GitHub Integration (Optional)</h3>
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-zinc-300 mb-1">GitHub Owner</label>
                  <input 
                    className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition"
                    value={formData.github_owner || ''}
                    onChange={e => setFormData({ ...formData, github_owner: e.target.value })}
                    placeholder="e.g. yashashkumar"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-zinc-300 mb-1">GitHub Repository</label>
                  <input 
                    className="w-full bg-zinc-950 border border-zinc-700 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition"
                    value={formData.github_repo || ''}
                    onChange={e => setFormData({ ...formData, github_repo: e.target.value })}
                    placeholder="e.g. ai-issue-agent"
                  />
                </div>
              </div>
            </div>
          </div>

          <div className="flex justify-end pt-4">
            <button 
              type="submit"
              className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition shadow-sm text-sm"
            >
              Create Project
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
