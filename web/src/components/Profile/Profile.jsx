import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Profile.module.css';

export default function Profile() {
    const token = localStorage.getItem('token');
    const claims = token ? JSON.parse(atob(token.split('.')[1])) : {};
    const [user, setUser] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await fetch(`http://localhost:80/api/v1/users/${claims.uid}`);
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
    }, [claims.uid]);

    const handleFavoritesRedirect = () => {
        navigate('/favorites');
    };

    const handleArticlesEditedRedirect = () => {
        navigate('/articles-edited');
    };

    const handleUsersRedirect = () => {
        navigate('/users');
    };

    const handleModerationRedirect = () => {
        navigate('/moderation');
    };

    const handleStatsArticlesRedirect = () => {
        navigate('/stats/articles');
    };

    const handleStatsUsersRedirect = () => {
        navigate('/stats/users');
    };

    if (!user) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className={styles.profile}>
            <h1 className={styles.title}>Профиль</h1>
            <div className={styles.photoContainer}>
                <img 
                    src={user.photoUrl || "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcROnah4MbL4lfJNAy5lgXMpFAnHnaQuiI0QuA&s"} 
                    alt="User Profile"
                    className={styles.photo}
                />
                <div className={styles.loginBlock}>
                    <p className={styles.email}><strong>Email:</strong> {user.email}</p>
                    <p className={styles.nick}><strong>Ник:</strong> {user.nick}</p>
                </div>
            </div>
            <button onClick={handleFavoritesRedirect} className={`${styles.button} btn btn-primary`}>
                Перейти в Избранное
            </button>

            {user.role === 'article_admin' && (
                <button onClick={handleArticlesEditedRedirect} className={`${styles.button} btn btn-secondary`}>
                    Перейти к редактированию статьей
                </button>
            )}
            {user.role === 'user_admin' && (
                <button onClick={handleUsersRedirect} className={`${styles.button} btn btn-secondary`}>
                    Перейти к редактированию пользователей
                </button>
            )}
            {user.role === 'moderator' && (
                <button onClick={handleModerationRedirect} className={`${styles.button} btn btn-secondary`}>
                    Перейти к модерации
                </button>
            )}
            {user.role === 'analyst' && (
                <>
                    <button onClick={handleStatsArticlesRedirect} className={`${styles.button} btn btn-secondary`}>
                        Перейти к статистике статей
                    </button>
                    <button onClick={handleStatsUsersRedirect} className={`${styles.button} btn btn-secondary`}>
                        Перейти к статистике пользователей
                    </button>
                </>
            )}
        </div>
    );
}