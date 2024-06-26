import React from 'react';
import { render, fireEvent, screen, waitFor, act } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import Cookies from 'js-cookie';
import User from './User';

describe('User component', () => {
  let mock: MockAdapter;
  const apiUrl = process.env.REACT_APP_API_URL;
  const mockUserData = {
    ID: 1,
    DeletedAt: null,
    name: 'John Doe',
    email: 'john@example.com',
    addresses: [
      {
        address_id: 1,
        country: 'Country',
        city: 'City',
        street: 'Street',
        number: 'Number',
        complement: 'Complement',
        state: 'State',
        zipcode: 'Zipcode',
      },
    ],
  };

  beforeEach(() => {
    mock = new MockAdapter(axios);
    Cookies.set = jest.fn();
    Cookies.get = jest.fn();
    Cookies.remove = jest.fn();
    global.alert = jest.fn();
    mock.onGet(apiUrl + '/user').reply(200, mockUserData);
  });

  afterEach(() => {
    mock.restore();
    jest.resetAllMocks();
  });

  it('fetches and displays user data on mount', async () => {
    await act(async () => {
      render(<User token="test-token" />);
    });

    await waitFor(() => {
      expect(screen.getByDisplayValue('John Doe')).toBeInTheDocument();
      expect(screen.getByDisplayValue('john@example.com')).toBeInTheDocument();
      expect(screen.getByText('Street, Number, Complement, City, State, Zipcode, Country')).toBeInTheDocument();
    });
  });

  it('allows the user to add a new address', async () => {
    await act(async () => {
      render(<User token="test-token" />);
    });

    fireEvent.click(screen.getByText(/add new address/i));

    fireEvent.change(screen.getByLabelText(/street/i), { target: { value: 'New Street' } });
    fireEvent.change(screen.getByLabelText(/number/i), { target: { value: 'New Number' } });
    fireEvent.change(screen.getByLabelText(/complement/i), { target: { value: 'New Complement' } });
    fireEvent.change(screen.getByLabelText(/city/i), { target: { value: 'New City' } });
    fireEvent.change(screen.getByLabelText(/state/i), { target: { value: 'New State' } });
    fireEvent.change(screen.getByLabelText(/zipcode/i), { target: { value: 'New Zipcode' } });
    fireEvent.change(screen.getByLabelText(/country/i), { target: { value: 'New Country' } });

    mock.onPost(apiUrl + '/user/address').reply(200, {
      address_id: 2,
      country: 'New Country',
      city: 'New City',
      street: 'New Street',
      number: 'New Number',
      complement: 'New Complement',
      state: 'New State',
      zipcode: 'New Zipcode',
    });

    fireEvent.click(screen.getByRole('button', { name: /add address/i }));

    await waitFor(() => {
      expect(global.alert).toHaveBeenCalledWith('Address added successfully');
      expect(screen.getByText('New Street, New Number, New Complement, New City, New State, New Zipcode, New Country')).toBeInTheDocument();
    });
  });

  it('allows the user to delete an address', async () => {
    await act(async () => {
      render(<User token="test-token" />);
    });

    mock.onDelete(apiUrl + '/user/address/1').reply(200);

    fireEvent.click(screen.getAllByText(/delete/i)[1]);

    await waitFor(() => {
      expect(global.alert).toHaveBeenCalledWith('Address deleted successfully');
      expect(screen.queryByText('Street, Number, Complement, City, State, Zipcode, Country')).not.toBeInTheDocument();
    });
  });

  it('allows the user to update their information', async () => {
    await act(async () => {
      render(<User token="test-token" />);
    });

    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'Jane Doe' } });
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'jane@example.com' } });

    mock.onPut(apiUrl + '/user').reply(200, {
      ID: 1,
      DeletedAt: null,
      name: 'Jane Doe',
      email: 'jane@example.com',
      addresses: mockUserData.addresses,
    });

    fireEvent.click(screen.getByText(/update user/i));

    await waitFor(() => {
      expect(global.alert).toHaveBeenCalledWith('User updated successfully');
      expect(screen.getByDisplayValue('Jane Doe')).toBeInTheDocument();
      expect(screen.getByDisplayValue('jane@example.com')).toBeInTheDocument();
    });
  });

});
