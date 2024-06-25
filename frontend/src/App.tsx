import React, { useState } from 'react';
import Register from './components/Register';
import Login from './components/Login';
import User from './components/User';

const App: React.FC = () => {
  const [token, setToken] = useState<string | null>(null);

  return (
    <div className="App">
      <h1>My App</h1>
      {!token ? (
        <>
          <Login setToken={setToken} />
          <Register />
        </>
      ) : (
        <User token={token} />
      )}
    </div>
  );
};

export default App;
