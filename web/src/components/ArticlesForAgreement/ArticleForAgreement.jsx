import React, { useEffect, useState } from 'react';
import Article from '../Article/Article';

const AgreementArticles = () => {
    const [articles, setArticles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchArticles = async () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    throw new Error('Ошибка: отсутствует токен авторизации');
                }
        
                console.log("Отправляем запрос с токеном:", token);
        
                const response = await fetch('http://localhost:80/api/v1/moderation/get', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });

                console.log("Полный ответ сервера:", response); // Выводим Response в консоль
        
                if (response.status === 401) {
                    throw new Error('Ошибка 401: Недействительный токен или недостаточно прав.');
                }
        
                if (!response.ok) {
                    throw new Error(`Ошибка сервера: ${response.status}`);
                }
        
                const data = await response.json();
                console.log("Article:", data); // Логируем JSON-ответ сервера
                setArticles(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchArticles();
    }, []);

    const handleAddArticle = async (article) => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                alert('Ошибка: отсутствует токен авторизации!');
                return;
            }

            const response = await fetch('http://localhost:80/api/v1/articles', {
                method: 'POST',
                mode: "no-cors",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify(article),
            });

            console.log("Полный ответ сервера на добавление статьи:", response); // Логируем ответ

            if (!response.ok) {
                throw new Error('Ошибка при добавлении статьи');
            }

            alert('Статья успешно добавлена!');
        } catch (err) {
            console.error(err);
            alert('Произошла ошибка при добавлении статьи');
        }
    };

    const handleRemoveArticle = async (articleId) => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                alert('Ошибка: отсутствует токен авторизации!');
                return;
            }

            const response = await fetch(`http://localhost:80/api/v1/moderation/remove?id=${articleId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`,
                },
            });

            console.log("Полный ответ сервера на удаление статьи:", response); // Логируем ответ

            if (!response.ok) {
                throw new Error('Ошибка при удалении статьи');
            }

            setArticles(articles.filter(article => article.id !== articleId));
            alert('Статья успешно удалена!');
        } catch (err) {
            console.error(err);
            alert('Произошла ошибка при удалении статьи');
        }
    };

    if (loading) {
        return <div>Загрузка статей...</div>;
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    return (
        <div>
            <h1>Статьи на модерации</h1>
            {articles.map(article => (
                <div key={article.id} className="article-item" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <Article article={article} />
                    <div style={{ display: 'flex', flexDirection: 'column', marginLeft: '10px' }}>
                        <button 
                            className="btn btn-success mb-2" 
                            onClick={() => handleAddArticle(article)}
                        >
                            Добавить
                        </button>
                        <button
                            className="btn btn-danger" 
                            onClick={() => handleRemoveArticle(article.id)}
                        >
                            Удалить
                        </button>
                    </div>
                </div>
            ))}
        </div>
    );
};

export default AgreementArticles;
