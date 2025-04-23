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
        const newCommentData = {
            article_id: article_id,
            owner_id: "вставьте_id_владельца_из_token",
            created_at: new Date().toISOString(),
            content: newComment,
        };

        try {
            const response = await fetch('http://localhost:80/api/v1/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(newCommentData),
            });

            if (!response.ok) {
                throw new Error('Ошибка при отправке комментария');
            }

            const result = await response.json();
            setComments([...comments, result]);
            setNewComment('');
        } catch (error) {
            console.error('Ошибка при отправке комментария:', error);
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

            <h3>Комментарии</h3>
            {comments.length === 0 ? (
                <p>Нет комментариев к этой статье.</p>
            ) : (
                <ul>
                    {comments.map(comment => (
                        <li key={comment.id}>
                            <strong>{comment.owner_id}</strong>: {comment.content}
                        </li>
                    ))}
                </ul>
            )}

            <form onSubmit={handleCommentSubmit}>
                <textarea 
                    value={newComment}
                    onChange={handleCommentChange}
                    placeholder="Введите ваш комментарий..."
                    required
                />
                <br />
                <button type="submit">Отправить</button>
            </form>
        </div>
    );
}
