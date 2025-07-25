import { useEffect, useState } from 'react';
import Article from '../Article/Article';
import './Favorites.css';

export default function Favorites() {
    const [favorites, setFavorites] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            setError('Ошибка: Токен отсутствует');
            setLoading(false);
            return;
        }

        let uid;
        try {
            const claims = JSON.parse(atob(token.split('.')[1]));
            uid = claims.uid;
            if (!uid) throw new Error("UID не найден в токене");
        } catch (err) {
            setError(`Ошибка при обработке токена: ${err.message}`);
            setLoading(false);
            return;
        }

        const fetchFavorites = async () => {
            try {
                const response = await fetch(`http://localhost/api/v1/favorites/get?id=${uid}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) {
                    throw new Error('Ошибка при получении избранных статей');
                }

                const data = await response.json();
                setFavorites(data);
            } catch (err) {
                setError(`Ошибка: ${err.message}`);
            } finally {
                setLoading(false);
            }
        };

        fetchFavorites();
    }, []);

    const removeFromFavorites = async (articleId) => {
        const token = localStorage.getItem('token');
        if (!token) {
            alert('Ошибка: Токен отсутствует');
            return;
        }

        let uid;
        try {
            const claims = JSON.parse(atob(token.split('.')[1]));
            uid = claims.uid;
            if (!uid) throw new Error("UID не найден в токене");
        } catch (err) {
            alert(`Ошибка при обработке токена: ${err.message}`);
            return;
        }

        try {
            const response = await fetch('http://localhost/api/v1/favorites/delete', {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ user_id: uid, article_id: articleId }),
            });

            if (!response.ok) {
                throw new Error('Ошибка при удалении из избранного');
            }

            setFavorites(favorites.filter(article => article.id !== articleId));
            alert('Статья удалена из избранного!');
        } catch (err) {
            alert(`Ошибка: ${err.message}`);
        }
    };

    if (loading) return <div>Загрузка избранных статей...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div>
            <h1>Избранные статьи</h1>
            {favorites.length === 0 ? (
                <p>Нет избранных статей</p>
            ) : (
                favorites.map(article => (
                    <div key={article.id} className="favorite-article-container">
                        <Article article={article} />
                        <button className="remove-favorite-button" onClick={() => removeFromFavorites(article.id)}>
                            ✖ Удалить
                        </button>
                    </div>
                ))
            )}
        </div>
    );
}
