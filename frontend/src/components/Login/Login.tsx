import React, { useState } from "react";
import axios from "axios";

interface LoginProps {
  setToken: (token: string) => void;
  onRegisterClick: () => void;
}

const Login: React.FC<LoginProps> = ({ setToken, onRegisterClick }) => {
  const apiUrl = process.env.REACT_APP_API_URL;
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  const handleSubmit  = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post(apiUrl + "/login", {
        email,
        password,
      });
      const { token } = response.data;
      setToken(token);
      alert("Login successful");
    } catch (error) {
      console.error(error);
      alert("Login failed");
    }
  };

  return (
    <div className="login-container">
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit">Login</button>
        <p className="register-toggle" onClick={onRegisterClick}>
          Register
        </p>
      </form>
    </div>
  );
};

export default Login;
