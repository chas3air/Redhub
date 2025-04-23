import styles from './Article.module.css';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Article({ article }) {
    const [ownerNick, setOwnerNick] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchOwnerData = async () => {
            try {
                const response = await fetch(`http://localhost:80/api/v1/users/${article.owner_id}`);
                const text = await response.text();
                console.log('Ответ сервера:', text);

                if (!response.ok) {
                    const errorData = JSON.parse(text);
                    console.error('Ошибка:', errorData);
                    return;
                }

                const data = JSON.parse(text);
                setOwnerNick(data.nick);
            } catch (error) {
                console.error('Ошибка при запросе:', error);
            }
        };

        fetchOwnerData();
    }, [article.owner_id]);

    const handleTitleClick = () => {
        navigate(`/articles/${article.id}`);
    };

    const addToFavorites = async () => {
        const token = localStorage.getItem('token'); // Получаем JWT токен
        if (!token) {
            alert('Ошибка: JWT токен отсутствует');
            return;
        }

        let uid;
        try {
            const claims = JSON.parse(atob(token.split('.')[1])); // Декодируем токен
            uid = claims.uid; // Извлекаем `uid`
            if (!uid) throw new Error("UID не найден в токене");
        } catch (err) {
            alert(`Ошибка при обработке токена: ${err.message}`);
            return;
        }

        try {
            const response = await fetch(`http://localhost/api/v1/favorites/add?id=${uid}`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(article), // Отправляем весь объект статьи
            });

            if (!response.ok) {
                throw new Error('Ошибка при добавлении в избранное');
            }

            alert('Статья добавлена в избранное!');
        } catch (err) {
            alert(`Ошибка: ${err.message}`);
        }
    };

    return (
        <div className={styles.block}>
            <div className={styles.favoriteButtonContainer}>
                <button className={styles.addFavoriteButton} onClick={addToFavorites}>★</button>
            </div>
            <div className={styles.info}> 
                <p>{ownerNick}</p>
                <p>{new Date(article.created_at).toLocaleDateString()}</p>
            </div>
            <h2 onClick={handleTitleClick} className={styles.title} style={{ cursor: 'pointer' }}>
                {article.title}
            </h2>
            {article.tag && (
                <span className={styles.tag}>
                    {article.tag}
                </span>
            )}
            <p>
                {article.content.length > 100 ? article.content.slice(0, 100) + '...' : article.content}
            </p>
        </div>
    );
}
