import { useState } from 'react';
import { Activity, Flame, Heart, Timer, ChevronRight, Settings, Home, User, BarChart2 } from 'lucide-react';
import './App.css';

function App() {
  const [activeTab, setActiveTab] = useState('home');

  return (
    <div className="app-container">
      <nav className="sidebar">
        <div className="logo-container">
          <Activity className="logo-icon" size={32} />
          <h1 className="logo-text">FitDark</h1>
        </div>
        <ul className="nav-links">
          <li className={activeTab === 'home' ? 'active' : ''} onClick={() => setActiveTab('home')}>
            <Home size={20} />
            <span>Dashboard</span>
          </li>
          <li className={activeTab === 'analytics' ? 'active' : ''} onClick={() => setActiveTab('analytics')}>
            <BarChart2 size={20} />
            <span>Analytics</span>
          </li>
          <li className={activeTab === 'profile' ? 'active' : ''} onClick={() => setActiveTab('profile')}>
            <User size={20} />
            <span>Profile</span>
          </li>
          <li className={activeTab === 'settings' ? 'active' : ''} onClick={() => setActiveTab('settings')}>
            <Settings size={20} />
            <span>Settings</span>
          </li>
        </ul>
      </nav>

      <main className="main-content">
        <header className="header">
          <h2>Welcome back, Alex!</h2>
          <p className="subtitle">Here's your activity for today.</p>
        </header>

        <section className="stats-grid">
          <div className="stat-card">
            <div className="stat-header">
              <span className="stat-title">Calories</span>
              <div className="icon-wrapper orange">
                <Flame size={20} />
              </div>
            </div>
            <div className="stat-value">1,240 <span className="stat-unit">kcal</span></div>
            <div className="stat-progress">
              <div className="progress-bar orange" style={{ width: '65%' }}></div>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-header">
              <span className="stat-title">Active Time</span>
              <div className="icon-wrapper blue">
                <Timer size={20} />
              </div>
            </div>
            <div className="stat-value">45 <span className="stat-unit">mins</span></div>
            <div className="stat-progress">
              <div className="progress-bar blue" style={{ width: '45%' }}></div>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-header">
              <span className="stat-title">Heart Rate</span>
              <div className="icon-wrapper red">
                <Heart size={20} />
              </div>
            </div>
            <div className="stat-value">112 <span className="stat-unit">bpm</span></div>
            <div className="stat-progress">
              <div className="progress-bar red" style={{ width: '80%' }}></div>
            </div>
          </div>
        </section>

        <section className="recent-activity">
          <div className="section-header">
            <h3>Recent Activity</h3>
            <button className="view-all">View All</button>
          </div>
          <div className="activity-list">
            <div className="activity-item">
              <div className="activity-icon run">
                <Activity size={20} />
              </div>
              <div className="activity-details">
                <h4>Morning Run</h4>
                <p>Today, 6:00 AM</p>
              </div>
              <div className="activity-metrics">
                <span className="metric">5.2 km</span>
                <span className="metric">320 kcal</span>
              </div>
              <ChevronRight className="chevron" size={20} />
            </div>

            <div className="activity-item">
              <div className="activity-icon cycle">
                <Timer size={20} />
              </div>
              <div className="activity-details">
                <h4>Cycling</h4>
                <p>Yesterday, 5:30 PM</p>
              </div>
              <div className="activity-metrics">
                <span className="metric">12.5 km</span>
                <span className="metric">450 kcal</span>
              </div>
              <ChevronRight className="chevron" size={20} />
            </div>
          </div>
        </section>

        <section className="weekly-goal">
          <h3>Weekly Goal Progress</h3>
          <div className="goal-container">
            <div className="goal-info">
              <span>Workouts</span>
              <span>4 / 5 days</span>
            </div>
            <div className="goal-progress-bar">
              <div className="goal-progress-fill" style={{ width: '80%' }}></div>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}

export default App;
