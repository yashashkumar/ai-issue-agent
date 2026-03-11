import axios from 'axios';
import { Project, Agent, AgentLog } from '../types';

const api = axios.create({
  baseURL: '/api/v1',
});

// Projects
export const getProjects = () => api.get<Project[]>('/projects').then(res => res.data);
export const getProject = (id: string) => api.get<Project>(`/projects/${id}`).then(res => res.data);
export const createProject = (data: Partial<Project>) => api.post<Project>('/projects', data).then(res => res.data);
export const updateProject = (id: string, data: Partial<Project>) => api.put<Project>(`/projects/${id}`, data).then(res => res.data);
export const deleteProject = (id: string) => api.delete(`/projects/${id}`);

// Agents
export const getAgents = () => api.get<Agent[]>('/agents').then(res => res.data);
export const getProjectAgents = (projectId: string) => api.get<Agent[]>(`/projects/${projectId}/agents`).then(res => res.data);
export const getAgent = (id: string) => api.get<Agent>(`/agents/${id}`).then(res => res.data);
export const getAgentLogs = (id: string, stream: 'stdout'|'stderr'|'all' = 'all') => api.get<AgentLog[]>(`/agents/${id}/logs?stream=${stream}`).then(res => res.data);
export const getAgentFiles = (id: string) => api.get<any[]>(`/agents/${id}/files`).then(res => res.data);
export const cancelAgent = (id: string) => api.post(`/agents/${id}/cancel`);

// Spawn
export const spawnAgent = (data: { system_prompt: string; project_id?: string; work_dir?: string; name?: string; description?: string }) => 
  api.post<{ agent_id: string; status: string }>('/spawn', data).then(res => res.data);
