import React, { useEffect, useState } from 'react';
import styles from './Profile.module.css';

export default function Profile() {
    const token = localStorage.getItem('token');
    const claims = token ? JSON.parse(atob(token.split('.')[1])) : {}; // Проверяем существование токена

    const [user, setUser] = useState(null); // Используем состояние для хранения данных пользователя

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await fetch(`http://localhost:80/api/v1/users/${claims.uid}`); // Используем uid из claims
                if (!response.ok) {
                    throw new Error('Ошибка при загрузке данных пользователя');
                }
                const userData = await response.json();
                setUser(userData);
            } catch (error) {
                console.error('Ошибка:', error);
            }
        };

        fetchUser();
    }, [claims.uid]); // Зависимость — uid

    if (!user) {
        return <div>Загрузка...</div>; // Предоставляем обратную связь при загрузке
    }

    return (
        <div className={styles.profile}>
            <h1>Профиль</h1>
            <div className={styles.photo}>
                <img 
                    src={user.photoUrl || "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcROnah4MbL4lfJNAy5lgXMpFAnHnaQuiI0QuA&s"} 
                    alt="User Profile"
                />
                <div className={styles.loginBlock}>
                    <p><strong>Ник:</strong> {user.nick}</p>
                    <p><strong>Email:</strong> {user.email}</p>
                </div>
            </div>
        </div>
    );
}