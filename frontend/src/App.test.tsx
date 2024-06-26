import React from 'react';
import { render, fireEvent, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import Cookies from 'js-cookie';
import App from './App';
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';

jest.mock('./components/Register/Register', () => () => <div>Register Component</div>);
jest.mock('./components/Login/Login', () => ({ setToken, onRegisterClick }: any) => (
  <div>
    <button onClick={() => setToken('test-token')}>Mock Login</button>
    <button onClick={onRegisterClick}>Mock Register</button>
  </div>
));
jest.mock('./components/User/User', () => ({ token }: any) => <div>User Component with token: {token}</div>);

describe('App component', () => {
  let mock: MockAdapter;

  beforeEach(() => {
    mock = new MockAdapter(axios);
    Cookies.set = jest.fn();
    Cookies.get = jest.fn();
  });

  afterEach(() => {
    mock.restore();
    jest.resetAllMocks();
  });

  it('renders Login component when no token is present', () => {
    (Cookies.get as jest.Mock).mockReturnValue(null);
    render(<App />);

    expect(screen.getByText('Mock Login')).toBeInTheDocument();
    expect(screen.queryByText('User Component')).not.toBeInTheDocument();
  });

  it('renders User component when token is present', () => {
    (Cookies.get as jest.Mock).mockReturnValue('test-token');
    render(<App />);

    expect(screen.getByText('User Component with token: test-token')).toBeInTheDocument();
    expect(screen.queryByText('Mock Login')).not.toBeInTheDocument();
  });

  it('shows Register component when Register button is clicked', () => {
    (Cookies.get as jest.Mock).mockReturnValue(null);
    render(<App />);

    fireEvent.click(screen.getByText('Mock Register'));
    expect(screen.getByText('Register Component')).toBeInTheDocument();
  });

  it('handles token setting correctly', async () => {
    (Cookies.get as jest.Mock).mockReturnValue(null);
    render(<App />);

    fireEvent.click(screen.getByText('Mock Login'));

    await waitFor(() => {
      expect(Cookies.set).toHaveBeenCalledWith('token', 'test-token');
      expect(screen.getByText('User Component with token: test-token')).toBeInTheDocument();
    });
  });
});
