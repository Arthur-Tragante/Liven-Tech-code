import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './User.css';

interface UserProps {
  token: string;
}

interface Address {
  address_id: number;
  country: string;
  city: string;
  street: string;
  number: string;
  complement: string;
  state: string;
  zipcode: string;
}

interface User {
  ID: number;
  DeletedAt: string | null;
  name: string;
  email: string;
  addresses: Address[];
}

const User: React.FC<UserProps> = ({ token }) => {
  const [user, setUser] = useState<User | null>(null);
  const [selectedAddress, setSelectedAddress] = useState<Address | null>(null);
  const [name, setName] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [country, setCountry] = useState<string>('');
  const [city, setCity] = useState<string>('');
  const [street, setStreet] = useState<string>('');
  const [number, setNumber] = useState<string>('');
  const [complement, setComplement] = useState<string>('');
  const [state, setState] = useState<string>('');
  const [zipcode, setZipcode] = useState<string>('');
  const [isFormVisible, setIsFormVisible] = useState<boolean>(false);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await axios.get('http://localhost:8080/user', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        setUser(response.data);
        setName(response.data.name);
        setEmail(response.data.email);
      } catch (error) {
        console.error(error);
        alert('Failed to fetch user data');
      }
    };

    fetchUser();
  }, [token]);

  const handleAddAddress = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post(
        'http://localhost:8080/user/address',
        { country, city, street, number, complement, state, zipcode },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setUser((prevUser) => {
        if (prevUser) {
          return { ...prevUser, addresses: [...prevUser.addresses, response.data] };
        }
        return prevUser;
      });
      setCountry('');
      setCity('');
      setStreet('');
      setNumber('');
      setComplement('');
      setState('');
      setZipcode('');
      alert('Address added successfully');
    } catch (error) {
      console.error(error);
      alert('Failed to add address');
    }
  };

  const handleDeleteAddress = async (addressID: number) => {
    try {
      await axios.delete(`http://localhost:8080/user/address/${addressID}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setUser((prevUser) => {
        if (prevUser) {
          return { ...prevUser, addresses: prevUser.addresses.filter((address) => address.address_id !== addressID) };
        }
        return prevUser;
      });
      alert('Address deleted successfully');
    } catch (error) {
      console.error(error);
      alert('Failed to delete address');
    }
  };

  const handleEditAddress = (address: Address) => {
    setSelectedAddress(address);
    setCountry(address.country);
    setCity(address.city);
    setStreet(address.street);
    setNumber(address.number);
    setComplement(address.complement);
    setState(address.state);
    setZipcode(address.zipcode);
    setIsFormVisible(true);
  };

  const handleUpdateAddress = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedAddress) return;

    try {
      const response = await axios.put(
        `http://localhost:8080/user/address/${selectedAddress.address_id}`,
        { country, city, street, number, complement, state, zipcode },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setUser((prevUser) => {
        if (prevUser) {
          const updatedAddresses = prevUser.addresses.map((address) =>
            address.address_id === selectedAddress.address_id ? response.data : address
          );
          return { ...prevUser, addresses: updatedAddresses };
        }
        return prevUser;
      });
      setSelectedAddress(null);
      setCountry('');
      setCity('');
      setStreet('');
      setNumber('');
      setComplement('');
      setState('');
      setZipcode('');
      setIsFormVisible(false);
      alert('Address updated successfully');
    } catch (error) {
      console.error(error);
      alert('Failed to update address');
    }
  };

  const handleUpdateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.put(
        'http://localhost:8080/user',
        { name, email },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setUser((prevUser) => {
        if (prevUser) {
          return { ...prevUser, name: response.data.name, email: response.data.email };
        }
        return prevUser;
      });
      alert('User updated successfully');
    } catch (error) {
      console.error(error);
      alert('Failed to update user');
    }
  };

  const handleDeleteUser = async () => {
    try {
      await axios.delete('http://localhost:8080/user', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      alert('User deleted successfully');
    } catch (error) {
      console.error(error);
      alert('Failed to delete user');
    }
  };

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <div className="user-container">
      <h2>User Details</h2>
      <form onSubmit={handleUpdateUser}>
        <div>
          <label>Name:</label>
          <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div>
          <label>Email:</label>
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <button className="update-user-button" type="submit">Update User</button>
        <button className="delete-user-button" type="button" onClick={handleDeleteUser}>Delete User</button>
      </form>
      <h3>Addresses</h3>
      <ul>
        {user.addresses && user.addresses.map((address) => (
          <li key={address.address_id}>
            {address.street}, {address.number}, {address.complement}, {address.city}, {address.state}, {address.zipcode}, {address.country}
            <div className="address-buttons">
              <button className="edit-address-button" onClick={() => handleEditAddress(address)}>Edit</button>
              <button className="delete-address-button" onClick={() => handleDeleteAddress(address.address_id)}>Delete</button>
            </div>
          </li>
        ))}
      </ul>
      <button className="toggle-form-button" onClick={() => setIsFormVisible(!isFormVisible)}>
        {isFormVisible ? 'Cancel' : 'Add New Address'}
      </button>
      {isFormVisible && (
        <div>
          <h3>{selectedAddress ? 'Update Address' : 'Add Address'}</h3>
          <form onSubmit={selectedAddress ? handleUpdateAddress : handleAddAddress}>
            <div>
              <label>Street:</label>
              <input type="text" value={street} onChange={(e) => setStreet(e.target.value)} />
            </div>
            <div>
              <label>Number:</label>
              <input type="text" value={number} onChange={(e) => setNumber(e.target.value)} />
            </div>
            <div>
              <label>Complement:</label>
              <input type="text" value={complement} onChange={(e) => setComplement(e.target.value)} />
            </div>
            <div>
              <label>City:</label>
              <input type="text" value={city} onChange={(e) => setCity(e.target.value)} />
            </div>
            <div>
              <label>State:</label>
              <input type="text" value={state} onChange={(e) => setState(e.target.value)} />
            </div>
            <div>
              <label>Zipcode:</label>
              <input type="text" value={zipcode} onChange={(e) => setZipcode(e.target.value)} />
            </div>
            <div>
              <label>Country:</label>
              <input type="text" value={country} onChange={(e) => setCountry(e.target.value)} />
            </div>
            <button className="submit-address-button" type="submit">{selectedAddress ? 'Update Address' : 'Add Address'}</button>
            {selectedAddress && <button className="cancel-edit-button" type="button" onClick={() => setSelectedAddress(null)}>Cancel</button>}
          </form>
        </div>
      )}
    </div>
  );
};

export default User;
