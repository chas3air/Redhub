import React, { useEffect, useState } from 'react';
import Article from '../Article/Article';
import './ListArticles.css';

const ListArticles = () => {
    const [articles, setArticles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [currentArticle, setCurrentArticle] = useState(null);
    const [isEditing, setIsEditing] = useState(false);
    const [updatedTitle, setUpdatedTitle] = useState('');
    const [updatedContent, setUpdatedContent] = useState('');

    useEffect(() => {
        const fetchArticles = async () => {
            try {
                const response = await fetch('http://localhost:80/api/v1/articles');
                if (!response.ok) {
                    throw new Error('Ошибка при получении статей');
                }
                const data = await response.json();
                setArticles(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchArticles();
    }, []);

    const handleDeleteArticle = async (articleId) => {
        try {
            const token = localStorage.getItem('token'); // Получаем токен из локального хранилища

            const response = await fetch(`http://localhost:80/api/v1/articles/${articleId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`, // Добавляем токен в заголовки
                },
            });
            if (!response.ok) {
                throw new Error('Ошибка при удалении статьи');
            }
            setArticles(articles.filter(article => article.id !== articleId));
        } catch (err) {
            console.error(err.message);
        }
    };

    const openEditPopup = (article) => {
        setCurrentArticle(article);
        setUpdatedTitle(article.title);
        setUpdatedContent(article.content);
        setIsEditing(true);
        console.log("Редактирование статьи:", article);
    };

    const handleUpdateArticle = async () => {
        if (!currentArticle) return;

        try {
            const token = localStorage.getItem('token'); // Получаем токен из локального хранилища
            console.log(token);
            const response = await fetch(`http://localhost:80/api/v1/articles/${currentArticle.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`, // Добавляем токен в заголовки
                },
                body: JSON.stringify({
                    title: updatedTitle,
                    content: updatedContent,
                }),
            });

            if (!response.ok) {
                throw new Error('Ошибка при обновлении статьи');
            }

            const updatedData = await response.json();
            setArticles(articles.map(article => (article.id === updatedData.id ? updatedData : article)));
            setIsEditing(false);
        } catch (err) {
            console.error(err.message);
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
            <h1>Редактируемые статьи</h1>
            {articles.map(article => (
                <div key={article.id} className="article-item" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <Article article={article} />
                    <div style={{ display: 'flex', flexDirection: 'column', marginLeft: '10px' }}>
                        <button 
                            className="btn btn-danger mb-2" 
                            onClick={() => handleDeleteArticle(article.id)}
                        >
                            Удалить
                        </button>
                        <button 
                            className="btn btn-warning" 
                            onClick={() => openEditPopup(article)}
                        >
                            Редактировать
                        </button>
                    </div>
                </div>
            ))}

            {isEditing && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Редактировать статью</h2>
                        <div>
                            <label>Название:</label>
                            <input 
                                type="text" 
                                value={updatedTitle} 
                                onChange={(e) => setUpdatedTitle(e.target.value)} 
                            />
                        </div>
                        <div>
                            <label>Содержимое:</label>
                            <textarea 
                                value={updatedContent} 
                                onChange={(e) => setUpdatedContent(e.target.value)} 
                            />
                        </div>
                        <button className="btn btn-primary" onClick={handleUpdateArticle}>
                            Сохранить
                        </button>
                        <button className="btn btn-secondary" onClick={() => setIsEditing(false)}>
                            Отмена
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ListArticles;