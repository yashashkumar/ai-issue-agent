import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getProjects } from '../api';
import { Project } from '../types';
import { FolderGit2, Plus } from 'lucide-react';

export function Projects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getProjects().then(data => {
      setProjects(data || []);
      setLoading(false);
    });
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Projects</h1>
        <Link 
          to="/projects/new" 
          className="inline-flex items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-indigo-700 transition"
        >
          <Plus className="mr-2 h-4 w-4" /> New Project
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {loading ? (
          <p className="text-zinc-400">Loading projects...</p>
        ) : projects.length === 0 ? (
          <div className="col-span-full py-12 text-center text-zinc-500 border border-dashed border-zinc-700 rounded-xl">
            No projects configured. Start by creating one.
          </div>
        ) : (
          projects.map(project => (
            <Link 
              key={project.id} 
              to={`/projects/${project.id}`}
              className="group relative flex flex-col justify-between overflow-hidden rounded-xl border border-zinc-800 bg-zinc-900 p-6 shadow-sm transition-all hover:bg-zinc-800/80 hover:shadow-md hover:border-zinc-700"
            >
              <div>
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-indigo-500/10 rounded-lg text-indigo-400">
                    <FolderGit2 className="h-5 w-5" />
                  </div>
                  <h3 className="font-semibold text-lg text-white group-hover:text-indigo-400 transition-colors">
                    {project.name}
                  </h3>
                </div>
                <p className="mt-4 text-sm text-zinc-400 line-clamp-2">
                  {project.description || "No description provided."}
                </p>
              </div>
              
              <div className="mt-6 pt-4 border-t border-zinc-800/50 flex align-center text-xs text-zinc-500 justify-between">
                <span>{project.github_owner && project.github_repo ? `${project.github_owner}/${project.github_repo}` : 'Local'}</span>
                <span>{new Date(project.created_at).toLocaleDateString()}</span>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
