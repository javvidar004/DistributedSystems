import React, { useState } from 'react';
import './App.css';

const API_BASE_URL = 'http://localhost:8080';

const VIEWS = {
  LOGIN: 'login',
  REGISTER: 'register',
  CHANGE_PASSWORD: 'change-password',
  DELETE_USER: 'delete-user',
  USERS_TABLE: 'users-table'
};

function App() {
  const [currentView, setCurrentView] = useState(VIEWS.LOGIN);
  const [feedback, setFeedback] = useState({ type: '', message: '' });
  const [isLoading, setIsLoading] = useState(false);

  const [loginData, setLoginData] = useState({ username: '', password: '' });
  const [registerData, setRegisterData] = useState({ name: '', username: '', password: '' });
  const [passwordData, setPasswordData] = useState({
    username: '',
    old_password: '',
    new_password: ''
  });
  const [deleteData, setDeleteData] = useState({ username: '' });
  const [users, setUsers] = useState([]);

  const setSuccess = (message) => setFeedback({ type: 'success', message });
  const setError = (message) => setFeedback({ type: 'error', message });

  const clearFeedback = () => setFeedback({ type: '', message: '' });

  const parseResponse = async (response) => {
    const raw = await response.text();
    if (!raw) {
      return {};
    }

    try {
      return JSON.parse(raw);
    } catch (error) {
      return { message: raw };
    }
  };

  const request = async (path, method, payload) => {
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json'
      }
    };

    if (payload) {
      options.body = JSON.stringify(payload);
    }

    const response = await fetch(`${API_BASE_URL}${path}`, options);
    const data = await parseResponse(response);

    if (!response.ok) {
      const message = data.message || data.error || `HTTP ${response.status}`;
      throw new Error(message);
    }

    return data;
  };

  const onLogin = async (event) => {
    event.preventDefault();
    clearFeedback();
    setIsLoading(true);
    try {
      await request('/login', 'POST', loginData);
      setSuccess('Login request sent successfully.');
    } catch (error) {
      setError(`Login failed: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const onRegister = async (event) => {
    event.preventDefault();
    clearFeedback();
    setIsLoading(true);
    try {
      await request('/register', 'POST', registerData);
      setSuccess('User registered successfully.');
      setRegisterData({ name: '', username: '', password: '' });
    } catch (error) {
      setError(`Register failed: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const onChangePassword = async (event) => {
    event.preventDefault();
    clearFeedback();
    setIsLoading(true);
    try {
      await request('/update', 'PUT', passwordData);
      setSuccess('Password updated successfully.');
      setPasswordData({ username: '', old_password: '', new_password: '' });
    } catch (error) {
      setError(`Password update failed: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const onDeleteUser = async (event) => {
    event.preventDefault();
    clearFeedback();
    setIsLoading(true);
    try {
      await request('/delete', 'DELETE', deleteData);
      setSuccess('User deleted successfully.');
      setDeleteData({ username: '' });
    } catch (error) {
      setError(`Delete failed: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const onFetchUsers = async () => {
    clearFeedback();
    setIsLoading(true);
    try {
      const data = await request('/users', 'GET');
      if (Array.isArray(data)) {
        setUsers(data);
      } else if (Array.isArray(data.users)) {
        setUsers(data.users);
      } else {
        setUsers([]);
      }
      setSuccess('Users loaded successfully.');
    } catch (error) {
      setUsers([]);
      setError(`Could not load users: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const renderLogin = () => (
    <form className="card" onSubmit={onLogin}>
      <h2>Login</h2>
      <label htmlFor="login-username">Username</label>
      <input
        id="login-username"
        type="text"
        value={loginData.username}
        onChange={(event) => setLoginData({ ...loginData, username: event.target.value })}
        required
      />

      <label htmlFor="login-password">Password</label>
      <input
        id="login-password"
        type="password"
        value={loginData.password}
        onChange={(event) => setLoginData({ ...loginData, password: event.target.value })}
        required
      />

      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Sending...' : 'Login'}
      </button>
    </form>
  );

  const renderRegister = () => (
    <form className="card" onSubmit={onRegister}>
      <h2>Register</h2>
      <label htmlFor="register-name">Name</label>
      <input
        id="register-name"
        type="text"
        value={registerData.name}
        onChange={(event) => setRegisterData({ ...registerData, name: event.target.value })}
        required
      />

      <label htmlFor="register-username">Username</label>
      <input
        id="register-username"
        type="text"
        value={registerData.username}
        onChange={(event) => setRegisterData({ ...registerData, username: event.target.value })}
        required
      />

      <label htmlFor="register-password">Password</label>
      <input
        id="register-password"
        type="password"
        value={registerData.password}
        onChange={(event) => setRegisterData({ ...registerData, password: event.target.value })}
        required
      />

      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Sending...' : 'Register'}
      </button>
    </form>
  );

  const renderPasswordUpdate = () => (
    <form className="card" onSubmit={onChangePassword}>
      <h2>Change Password</h2>
      <label htmlFor="change-username">Username</label>
      <input
        id="change-username"
        type="text"
        value={passwordData.username}
        onChange={(event) => setPasswordData({ ...passwordData, username: event.target.value })}
        required
      />

      <label htmlFor="change-old-password">Old Password</label>
      <input
        id="change-old-password"
        type="password"
        value={passwordData.old_password}
        onChange={(event) => setPasswordData({ ...passwordData, old_password: event.target.value })}
        required
      />

      <label htmlFor="change-new-password">New Password</label>
      <input
        id="change-new-password"
        type="password"
        value={passwordData.new_password}
        onChange={(event) => setPasswordData({ ...passwordData, new_password: event.target.value })}
        required
      />

      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Sending...' : 'Update Password'}
      </button>
    </form>
  );

  const renderDeleteUser = () => (
    <form className="card" onSubmit={onDeleteUser}>
      <h2>Delete User</h2>
      <label htmlFor="delete-username">Username</label>
      <input
        id="delete-username"
        type="text"
        value={deleteData.username}
        onChange={(event) => setDeleteData({ username: event.target.value })}
        required
      />

      <button type="submit" className="danger" disabled={isLoading}>
        {isLoading ? 'Sending...' : 'Delete User'}
      </button>
    </form>
  );

  const renderUsersTable = () => (
    <section className="card table-card">
      <h2>Registered Users</h2>
      <button type="button" onClick={onFetchUsers} disabled={isLoading}>
        {isLoading ? 'Loading...' : 'Load Users'}
      </button>

      {users.length === 0 ? (
        <p className="empty">No users loaded yet.</p>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Username</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user, index) => (
                <tr key={user.id || `${user.username || 'user'}-${index}`}>
                  <td>{user.id}</td>
                  <td>{user.name}</td>
                  <td>{user.username}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );

  const renderCurrentView = () => {
    if (currentView === VIEWS.LOGIN) {
      return renderLogin();
    }

    if (currentView === VIEWS.REGISTER) {
      return renderRegister();
    }

    if (currentView === VIEWS.CHANGE_PASSWORD) {
      return renderPasswordUpdate();
    }

    if (currentView === VIEWS.DELETE_USER) {
      return renderDeleteUser();
    }

    return renderUsersTable();
  };

  return (
    <div className="App">
      <div className="bg-shape bg-shape-left" aria-hidden="true" />
      <div className="bg-shape bg-shape-right" aria-hidden="true" />

      <main className="layout">
        <header className="hero">
          <h1>User Management</h1>
          <p>Frontend connected to localhost:8080</p>
        </header>

        <nav className="nav-tabs" aria-label="Main navigation">
          <button type="button" onClick={() => setCurrentView(VIEWS.LOGIN)} className={currentView === VIEWS.LOGIN ? 'active' : ''}>
            Login
          </button>
          <button type="button" onClick={() => setCurrentView(VIEWS.REGISTER)} className={currentView === VIEWS.REGISTER ? 'active' : ''}>
            Register
          </button>
          <button type="button" onClick={() => setCurrentView(VIEWS.CHANGE_PASSWORD)} className={currentView === VIEWS.CHANGE_PASSWORD ? 'active' : ''}>
            Change Password
          </button>
          <button type="button" onClick={() => setCurrentView(VIEWS.DELETE_USER)} className={currentView === VIEWS.DELETE_USER ? 'active' : ''}>
            Delete User
          </button>
          <button type="button" onClick={() => setCurrentView(VIEWS.USERS_TABLE)} className={currentView === VIEWS.USERS_TABLE ? 'active' : ''}>
            Users Table
          </button>
        </nav>

        {feedback.message && (
          <p className={`feedback ${feedback.type === 'error' ? 'error' : 'success'}`}>
            {feedback.message}
          </p>
        )}

        {renderCurrentView()}
      </main>
    </div>
  );
}

export default App;
