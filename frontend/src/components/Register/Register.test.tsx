import { render, fireEvent, screen, waitFor, act } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import Register from './Register';

describe('Register component', () => {
  let mock: MockAdapter;

  beforeEach(() => {
    mock = new MockAdapter(axios);
    global.alert = jest.fn();
  });

  afterEach(() => {
    mock.restore();
    jest.resetAllMocks();
  });

  it('allows the user to register successfully', async () => {
    mock.onPost('http://localhost:8080/register').reply(200, { message: 'Registration successful' });

    await act(async () => {
      render(<Register />);
    });

    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'John Doe' } });
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'john@example.com' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /register/i }));

    await waitFor(() => {
      expect(global.alert).toHaveBeenCalledWith('Registration successful');
    });
  });

  it('shows an error message on registration failure', async () => {
    mock.onPost('http://localhost:8080/register').reply(500, { message: 'Registration failed' });

    await act(async () => {
      render(<Register />);
    });

    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'John Doe' } });
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'john@example.com' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /register/i }));

    await waitFor(() => {
      expect(global.alert).toHaveBeenCalledWith('Registration failed');
    });
  });
});
