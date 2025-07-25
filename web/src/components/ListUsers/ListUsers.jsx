import React, { useEffect, useState } from 'react';
import './ListUsers.css';

const ListUsers = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const response = await fetch('http://localhost:80/api/v1/users');
                if (!response.ok) {
                    throw new Error('Ошибка при получении пользователей');
                }
                const data = await response.json();
                setUsers(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchUsers();
    }, []);

    const handleDeleteUser = async (userId) => {
        try {
            const token = localStorage.getItem('token');

            const response = await fetch(`http://localhost:80/api/v1/users/${userId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`,
                },
            });
            if (!response.ok) {
                throw new Error('Ошибка при удалении пользователя');
            }
            setUsers(users.filter(user => user.id !== userId));
        } catch (err) {
            console.error(err.message);
        }
    };

    if (loading) {
        return <div>Загрузка пользователей...</div>;
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    return (
        <div>
            <h1>Пользователи</h1>
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Имя</th>
                        <th>Email</th>
                        <th>Role</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {users.map(user => (
                        <tr key={user.id}>
                            <td>{user.id}</td>
                            <td>{user.nick}</td>
                            <td>{user.email}</td>
                            <td>{user.role}</td>
                            <td>
                                <button 
                                    className="btn btn-danger" 
                                    onClick={() => handleDeleteUser(user.id)}
                                >
                                    Удалить
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default ListUsers;