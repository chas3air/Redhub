import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import './SingleArticle.css';

export default function SingleArticle() {
    const { article_id } = useParams();
    const [article, setArticle] = useState(null);
    const [owner, setOwner] = useState('');
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState('');

    useEffect(() => {
        const fetchArticleData = async () => {
            try {
                const articleResponse = await fetch(`http://localhost:80/api/v1/articles/${article_id}/`);
                const articleData = await articleResponse.json();

                if (!articleResponse.ok) {
                    throw new Error('Ошибка при получении статьи');
                }

                setArticle(articleData);
                fetchOwnerData(articleData.owner_id);
                fetchCommentsData(article_id);
            } catch (error) {
                console.error('Ошибка при запросе статьи:', error);
            }
        };

        const fetchOwnerData = async (ownerId) => {
            try {
                const response = await fetch(`http://localhost:80/api/v1/users/${ownerId}`);
                const data = await response.json();

                if (!response.ok) {
                    throw new Error('Ошибка при получении владельца');
                }

                setOwner(data);
            } catch (error) {
                console.error('Ошибка при запросе владельца:', error);
            }
        };

        const fetchCommentsData = async (articleId) => {
            try {
                const response = await fetch(`http://localhost:80/api/v1/${articleId}/comments`);
                const data = await response.json();

                if (!response.ok) {
                    throw new Error('Ошибка при получении комментариев');
                }

                setComments(Array.isArray(data) ? data : []);
            } catch (error) {
                console.error('Ошибка при запросе комментариев:', error);
            }
        };

        fetchArticleData();
    }, [article_id]);

    const handleCommentChange = (event) => {
        setNewComment(event.target.value);
    };

    const handleCommentSubmit = async (event) => {
        event.preventDefault();
        const token = localStorage.getItem('token');
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
            const response = await fetch(`http://localhost:80/api/v1/users/${uid}`);
            const data = await response.json();

            if (!response.ok) {
                throw new Error('Ошибка при получении пользователя');
            }

            const currentDate = new Date().toLocaleString('ru-RU', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });

            const newCommentData = {
                owner_id: data.nick, // Используем ник вместо ID
                created_at: currentDate,
                content: newComment,
            };

            setComments(prevComments => [...prevComments, newCommentData]);
            setNewComment('');
        } catch (error) {
            console.error('Ошибка при получении пользователя:', error);
        }
    };

    if (!article) {
        return <div>Загрузка статьи...</div>;
    }

    return (
        <div className="article-container">
            <div className="article-header">
                <img src="../../assets/logo512.png" alt="article" className="article-image" style={{ width: '20px', height: '20px' }} />
                <span className="author-nick">{owner.nick}</span>
                <span className="publication-date">
                    {new Date(article.created_at).toLocaleString('ru-RU', { 
                        year: 'numeric',
                        month: '2-digit',
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit'
                    })}
                </span>
            </div>
            <h2>{article.title}</h2>
            {article.tag && (
                <span className="article-tag">{article.tag}</span>
            )}
            
            <div>{article.content}</div>

            <h3>Добавить комментарий</h3>
            <form onSubmit={handleCommentSubmit} className="comment-form">
                <div className="comment-input-container">
                    <textarea 
                        value={newComment}
                        onChange={handleCommentChange}
                        placeholder="Введите ваш комментарий..."
                        required
                        className="comment-input"
                    />
                    <button type="submit" className="submit-button">Отправить</button>
                </div>
            </form>

            <h3>Комментарии</h3>
            {comments.length === 0 ? (
                <p>Нет комментариев к этой статье.</p>
            ) : (
                <div className="comments-container">
                    {comments.map((comment, index) => (
                        <div key={index} className="comment-block" style={{ border: '1px solid #ccc', borderRadius: '5px', padding: '10px', marginBottom: '10px' }}>
                            <div className="comment-header">
                                <strong>{comment.owner_id || 'Аноним'}</strong>
                                <span className="comment-date" style={{ marginLeft: '10px' }}>{comment.created_at}</span>
                            </div>
                            <div className="comment-content">{comment.content}</div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}