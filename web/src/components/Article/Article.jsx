import styles from './Article.module.css';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Article({ article }) {
    const [ownerNick, setOwnerNick] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const [successMessage, setSuccessMessage] = useState('');
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
        const token = localStorage.getItem('token');
        if (!token) {
            setErrorMessage('Ошибка: сначала зарегистрируйтесь.');
            return;
        }

        let uid;
        try {
            const claims = JSON.parse(atob(token.split('.')[1]));
            uid = claims.uid;
            if (!uid) throw new Error("UID не найден в токене");
        } catch (err) {
            setErrorMessage(`Ошибка при обработке токена: ${err.message}`);
            return;
        }

        try {
            const response = await fetch(`http://localhost/api/v1/favorites/add?id=${uid}`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(article),
            });

            if (!response.ok) {
                throw new Error('Ошибка при добавлении в избранное');
            }

            setSuccessMessage('Статья добавлена в избранное!');
            setErrorMessage(''); // Сброс сообщения об ошибке
        } catch (err) {
            setErrorMessage(`Ошибка: ${err.message}`);
            setSuccessMessage(''); // Сброс сообщения об успехе
        }
    };

    return (
        <div className={styles.block}>
            {errorMessage && (
                <div className={styles.errorMessage}>
                    {errorMessage}
                </div>
            )}
            {successMessage && (
                <div className={styles.successMessage}>
                    {successMessage}
                </div>
            )}
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