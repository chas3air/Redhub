import styles from './Article.module.css';
import { useEffect, useState } from 'react';

export default function Article({ article }) {
    const [ownerNick, setOwnerNick] = useState('');

    useEffect(() => {
        const fetchOwnerData = async () => {
            try {
                const response = await fetch(`http://localhost:80/api/v1/users/${article.owner_id}`);
                const text = await response.text(); // Получаем текст ответа
                console.log('Ответ сервера:', text); // Выводим ответ в консоль
    
                if (!response.ok) {
                    const errorData = JSON.parse(text); // Пробуем парсить JSON
                    console.error('Ошибка:', errorData);
                    return;
                }
    
                const data = JSON.parse(text); // Парсим текст как JSON
                setOwnerNick(data.nick);
            } catch (error) {
                console.error('Ошибка при запросе:', error);
            }
        };
    
        fetchOwnerData();
    }, [article.owner_id]);

    return (
        <div key={article.id} className={styles.block}>
            <div className={styles.info}>
                <p>{ownerNick}</p>
                <p>{new Date(article.created_at).toLocaleDateString()}</p>
            </div>
            <h2>{article.title}</h2>
            <p>
                {article.content.length > 100 ? article.content.slice(0, 100) + '...' : article.content}
            </p>
        </div>
    );
}