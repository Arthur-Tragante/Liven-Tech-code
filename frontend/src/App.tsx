import React, { useState } from 'react';
import Register from './components/Register';
import Login from './components/Login';
import User from './components/User';
import './styles.css';

const App: React.FC = () => {
  const [token, setToken] = useState<string | null>(null);
  const [showRegister, setShowRegister] = useState<boolean>(false);

  const handleShowRegister = () => {
    setShowRegister(!showRegister);
  };

  return (
    <div className="App">
      <h1>User Address Control</h1>
      {!token ? (
        <>
          <Login setToken={setToken} onRegisterClick={handleShowRegister} />
          {showRegister && <Register />}
        </>
      ) : (
        <User token={token} />
      )}
    </div>
  );
};

export default App;
