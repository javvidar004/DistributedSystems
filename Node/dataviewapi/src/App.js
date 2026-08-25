import React, { useEffect, useMemo, useState } from 'react';
import './App.css';

const API_BASE_URL = 'http://localhost:8080';

const initialUser = {
  name: '',
  last_name: '',
  username: '',
  email: '',
};

function App() {
  const [authMode, setAuthMode] = useState('signin');
  const [view, setView] = useState('dashboard');
  const [currentUser, setCurrentUser] = useState(null);
  const [userData, setUserData] = useState(initialUser);
  const [loginData, setLoginData] = useState({ username: '', password: '' });
  const [signupData, setSignupData] = useState({
    name: '',
    last_name: '',
    username: '',
    email: '',
    password: '',
  });
  const [passwordData, setPasswordData] = useState({
    username: '',
    old_password: '',
    new_password: '',
  });
  const [message, setMessage] = useState({ type: '', text: '' });
  const [loading, setLoading] = useState(false);
  const [offices, setOffices] = useState([]);
  const [selectedOfficeId, setSelectedOfficeId] = useState(1);
  const [workspaces, setWorkspaces] = useState([]);
  const [selectedWorkspace, setSelectedWorkspace] = useState(null);
  const [bookingDate, setBookingDate] = useState('');
  const [futureBookings, setFutureBookings] = useState([]);

  const setSuccess = (text) => setMessage({ type: 'success', text });
  const setError = (text) => setMessage({ type: 'error', text });

  const request = async (path, method = 'GET', payload = null) => {
    const options = {
      method,
      headers: { 'Content-Type': 'application/json' },
    };

    if (payload) {
      options.body = JSON.stringify(payload);
    }

    const response = await fetch(`${API_BASE_URL}${path}`, options);
    const contentType = response.headers.get('content-type') || '';
    const body = contentType.includes('application/json') ? await response.json() : await response.text();

    if (!response.ok) {
      const message = typeof body === 'string' ? body : body.message || 'Request failed';
      throw new Error(message);
    }

    return body;
  };

  const fetchOffices = async () => {
    try {
      const data = await request('/offices');
      setOffices(Array.isArray(data) ? data : []);
      if (data[0]) {
        setSelectedOfficeId(data[0].id);
      }
    } catch (error) {
      console.error(error);
    }
  };

  const fetchWorkspaces = async (officeId = selectedOfficeId) => {
    try {
      const data = await request(`/offices/${officeId}/workspaces`);
      setWorkspaces(Array.isArray(data) ? data : []);
      setSelectedWorkspace(null);
    } catch (error) {
      console.error(error);
    }
  };

  const fetchFutureBookings = async (username = currentUser?.username) => {
    if (!username) {
      setFutureBookings([]);
      return;
    }

    try {
      const data = await request(`/bookings?username=${encodeURIComponent(username)}`);
      setFutureBookings(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchOffices();
  }, []);

  useEffect(() => {
    if (selectedOfficeId) {
      fetchWorkspaces(selectedOfficeId);
    }
  }, [selectedOfficeId]);

  useEffect(() => {
    if (currentUser?.username) {
      fetchFutureBookings(currentUser.username);
    }
  }, [currentUser]);

  const today = useMemo(() => new Date().toISOString().slice(0, 10), []);

  const handleLogin = async (event) => {
    event.preventDefault();
    setLoading(true);
    setMessage({ type: '', text: '' });

    try {
      const response = await request('/login', 'POST', loginData);
      const user = response.user || response;
      setCurrentUser(user);
      setView('dashboard');
      setSuccess('Welcome back.');
      setUserData({
        name: user.name || '',
        last_name: user.last_name || '',
        username: user.username || '',
        email: user.email || '',
      });
      setPasswordData({ username: user.username || '', old_password: '', new_password: '' });
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (event) => {
    event.preventDefault();
    setLoading(true);
    setMessage({ type: '', text: '' });

    try {
      await request('/register', 'POST', signupData);
      setSuccess('Account created successfully.');
      setAuthMode('signin');
      setSignupData({
        name: '',
        last_name: '',
        username: '',
        email: '',
        password: '',
      });
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordUpdate = async (event) => {
    event.preventDefault();
    setLoading(true);
    setMessage({ type: '', text: '' });

    try {
      await request('/update', 'PUT', passwordData);
      setSuccess('Password updated successfully.');
      setPasswordData((prev) => ({ ...prev, old_password: '', new_password: '' }));
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleBooking = async () => {
    if (!selectedWorkspace || !currentUser || !bookingDate) {
      setError('Select a workspace and choose a date first.');
      return;
    }

    setLoading(true);
    setMessage({ type: '', text: '' });

    try {
      await request('/bookings', 'POST', {
        user_id: currentUser.id,
        workspace_id: selectedWorkspace.id,
        booking_date: bookingDate,
      });
      setSuccess(`Workspace ${selectedWorkspace.number} booked for ${bookingDate}.`);
      setSelectedWorkspace(null);
      setBookingDate('');
      fetchFutureBookings(currentUser.username);
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    setCurrentUser(null);
    setView('dashboard');
    setAuthMode('signin');
    setUserData(initialUser);
    setMessage({ type: 'success', text: 'You have been signed out.' });
  };

  const renderAuth = () => (
    <div className="auth-shell">
      <div className="auth-card">
        <div className="auth-switch">
          <button className={authMode === 'signin' ? 'active' : ''} onClick={() => setAuthMode('signin')}>
            Sign in
          </button>
          <button className={authMode === 'signup' ? 'active' : ''} onClick={() => setAuthMode('signup')}>
            Sign up
          </button>
        </div>

        {authMode === 'signin' ? (
          <form className="stack-form" onSubmit={handleLogin}>
            <h2>Welcome back</h2>
            <label>
              Username
              <input value={loginData.username} onChange={(e) => setLoginData({ ...loginData, username: e.target.value })} required />
            </label>
            <label>
              Password
              <input type="password" value={loginData.password} onChange={(e) => setLoginData({ ...loginData, password: e.target.value })} required />
            </label>
            <button type="submit" disabled={loading}>{loading ? 'Signing in...' : 'Sign in'}</button>
          </form>
        ) : (
          <form className="stack-form" onSubmit={handleRegister}>
            <h2>Create account</h2>
            <div className="two-col">
              <label>
                Name
                <input value={signupData.name} onChange={(e) => setSignupData({ ...signupData, name: e.target.value })} required />
              </label>
              <label>
                Last name
                <input value={signupData.last_name} onChange={(e) => setSignupData({ ...signupData, last_name: e.target.value })} required />
              </label>
            </div>
            <label>
              Username
              <input value={signupData.username} onChange={(e) => setSignupData({ ...signupData, username: e.target.value })} required />
            </label>
            <label>
              Email
              <input type="email" value={signupData.email} onChange={(e) => setSignupData({ ...signupData, email: e.target.value })} required />
            </label>
            <label>
              Password
              <input type="password" value={signupData.password} onChange={(e) => setSignupData({ ...signupData, password: e.target.value })} required />
            </label>
            <button type="submit" disabled={loading}>{loading ? 'Creating...' : 'Sign up'}</button>
          </form>
        )}
      </div>
    </div>
  );

  const renderDashboard = () => (
    <div className="content-panel">
      <div className="section-header">
        <h2>Upcoming bookings</h2>
      </div>
      <div className="booking-list">
        {futureBookings.length === 0 ? (
          <p className="empty-state">No future bookings yet.</p>
        ) : (
          futureBookings.map((booking) => (
            <div className="booking-item" key={booking.id}>
              <div>
                <strong>Workspace {booking.workspace?.number || booking.workspace_id}</strong>
              </div>
              <span>{new Date(booking.booking_date).toLocaleDateString()}</span>
            </div>
          ))
        )}
      </div>
    </div>
  );

  const renderUserInfo = () => (
    <div className="content-panel stacked">
      <div className="section-header">
        <h2>User information</h2>
      </div>
      <div className="profile-grid">
        <div className="profile-card">
          <p><strong>Name:</strong> {userData.name}</p>
          <p><strong>Last name:</strong> {userData.last_name}</p>
          <p><strong>Username:</strong> {userData.username}</p>
          <p><strong>Email:</strong> {userData.email}</p>
        </div>

        <form className="stack-form compact" onSubmit={handlePasswordUpdate}>
          <h3>Update password</h3>
          <label>
            Current password
            <input type="password" value={passwordData.old_password} onChange={(e) => setPasswordData({ ...passwordData, old_password: e.target.value })} required />
          </label>
          <label>
            New password
            <input type="password" value={passwordData.new_password} onChange={(e) => setPasswordData({ ...passwordData, new_password: e.target.value })} required />
          </label>
          <button type="submit" disabled={loading}>{loading ? 'Updating...' : 'Update password'}</button>
        </form>
      </div>
    </div>
  );

  const renderOffice = () => (
    <div className="content-panel stacked">
      <div className="section-header">
        <h2>Book a workspace</h2>
      </div>

      <div className="office-toolbar">
        <label>
          Office
          <select value={selectedOfficeId} onChange={(e) => setSelectedOfficeId(Number(e.target.value))}>
            {offices.map((office) => (
              <option key={office.id} value={office.id}>{office.name}</option>
            ))}
          </select>
        </label>
      </div>

      <div className="office-grid">
        {workspaces.map((workspace) => (
          <button
            key={workspace.id}
            className={`workspace ${selectedWorkspace?.id === workspace.id ? 'selected' : ''}`}
            onClick={() => setSelectedWorkspace(workspace)}
            type="button"
            title={`Workspace ${workspace.number}`}
          >
            {workspace.number}
          </button>
        ))}
      </div>

      {selectedWorkspace && (
        <div className="booking-form">
          <h3>Selected workspace: {selectedWorkspace.number}</h3>
          <label>
            Booking date
            <input type="date" min={today} value={bookingDate} onChange={(e) => setBookingDate(e.target.value)} required />
          </label>
          <button type="button" onClick={handleBooking} disabled={loading}>Confirm booking</button>
        </div>
      )}
    </div>
  );

  if (!currentUser) {
    return (
      <div className="app-shell">
        {message.text && <div className={`flash ${message.type}`}>{message.text}</div>}
        {renderAuth()}
      </div>
    );
  }

  return (
    <div className="app-shell">
      <header className="topbar">
        <div>
          <h1>Workspace Hub</h1>
        </div>
        <nav className="nav-bar">
          <button className={view === 'dashboard' ? 'active' : ''} onClick={() => setView('dashboard')}>Main</button>
          <button className={view === 'user' ? 'active' : ''} onClick={() => setView('user')}>User info</button>
          <button className={view === 'office' ? 'active' : ''} onClick={() => setView('office')}>Office</button>
          <button onClick={handleLogout}>Sign out</button>
        </nav>
      </header>

      {message.text && <div className={`flash ${message.type}`}>{message.text}</div>}

      {view === 'dashboard' && renderDashboard()}
      {view === 'user' && renderUserInfo()}
      {view === 'office' && renderOffice()}
    </div>
  );
}

export default App;
