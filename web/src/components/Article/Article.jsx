import styles from './Article.module.css';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom'; // Импортируем useNavigate

export default function Article({ article }) {
    const [ownerNick, setOwnerNick] = useState('');
    const navigate = useNavigate(); // Получаем функцию для навигации

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

    const handleArticleClick = () => {
        navigate(`/articles/${article.id}`); // Перенаправление на страницу статьи
    };

    return (
        <div key={article.id} className={styles.block} onClick={handleArticleClick} style={{ cursor: 'pointer' }}>
            <div className={styles.info}>
                <p>{ownerNick}</p>
                <p>{new Date(article.created_at).toLocaleDateString()}</p>
            </div>
            <h2>{article.title}</h2>
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