import React, { useState, useEffect } from 'react';
import Cookies from 'js-cookie';
import Register from './components/Register';
import Login from './components/Login';
import User from './components/User';
import './styles.css';

const App: React.FC = () => {
  const [token, setToken] = useState<string | null>(null);
  const [showRegister, setShowRegister] = useState<boolean>(false);

  useEffect(() => {
    const savedToken = Cookies.get('token');
    if (savedToken) {
      setToken(savedToken);
    }
  }, []);

  const handleSetToken = (newToken: string) => {
    Cookies.set('token', newToken);
    setToken(newToken);
  };

  const handleShowRegister = () => {
    setShowRegister(!showRegister);
  };

  return (
    <div className="App">
      <h1>User Address Control</h1>
      {!token ? (
        <>
          <Login setToken={handleSetToken} onRegisterClick={handleShowRegister} />
          {showRegister && <Register />}
        </>
      ) : (
        <User token={token} />
      )}
    </div>
  );
};

export default App;
