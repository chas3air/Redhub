import React, { useEffect, useState } from 'react';
import User from '../User/User';

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
            setUsers(users.filter(user => user.id !== userId)); // Обновляем состояние
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
            {users.map(user => (
                <User key={user.id} user={user} onDelete={handleDeleteUser} /> // Передаем функцию удаления
            ))}
        </div>
    );
};

export default ListUsers;