import React from 'react';
import { render, fireEvent, screen, waitFor, act } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import Login from './Login';

describe('Login component', () => {
  let mock: MockAdapter;
  const apiUrl = process.env.REACT_APP_API_URL;
  beforeEach(() => {
    mock = new MockAdapter(axios);
    global.alert = jest.fn();
  });

  afterEach(() => {
    mock.restore();
    jest.resetAllMocks();
  });

  it('allows the user to log in successfully', async () => {
    const setToken = jest.fn();
    const onRegisterClick = jest.fn();

    mock.onPost(apiUrl + '/login').reply(200, { token: 'test-token' });

    await act(async () => {
      render(<Login setToken={setToken} onRegisterClick={onRegisterClick} />);
    });

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(setToken).toHaveBeenCalledWith('test-token');
    });
  });

  it('shows an error message on login failure', async () => {
    const setToken = jest.fn();
    const onRegisterClick = jest.fn();

    mock.onPost(apiUrl + '/login').reply(401, { message: 'Invalid credentials' });

    await act(async () => {
      render(<Login setToken={setToken} onRegisterClick={onRegisterClick} />);
    });

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'wrongpassword' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(global.alert).toHaveBeenCalledWith('Login failed');
    });
  });

  it('calls onRegisterClick when register is clicked', () => {
    const onRegisterClick = jest.fn();
    const setToken = jest.fn();

    render(<Login setToken={setToken} onRegisterClick={onRegisterClick} />);

    fireEvent.click(screen.getByText(/register/i));

    expect(onRegisterClick).toHaveBeenCalled();
  });
});
