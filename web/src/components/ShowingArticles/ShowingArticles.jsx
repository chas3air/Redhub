import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import Article from '../Article/Article';
import './ShowingArticles.css';

const ShowingArticles = () => {
    const [articles, setArticles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedTag, setSelectedTag] = useState('');

    const tagOptions = ['Все', 'Новости', 'Технология', 'Язык программирования', 'Обзор', 'Реклама'];

    useEffect(() => {
        const fetchArticles = async () => {
            try {
                const response = await fetch('http://localhost:80/api/v1/articles');
                if (!response.ok) {
                    throw new Error('Ошибка при получении статей');
                }
                const data = await response.json();
                console.log("Полученные статьи:", data); // Проверяем данные
                if (Array.isArray(data)) {
                    setArticles(data);
                } else {
                    throw new Error("Неверный формат данных: ожидался массив");
                }
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchArticles();
    }, []);

    if (loading) {
        return <div>Загрузка статей...</div>;
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    const filteredArticles = Array.isArray(articles) ? articles.filter(article =>
        article.title.toLowerCase().includes(searchTerm.toLowerCase()) &&
        (selectedTag === 'Все' || !selectedTag || article.tag?.toLowerCase() === selectedTag.toLowerCase())
    ) : [];

    return (
        <div>
            <h1>Статьи</h1>
            <Link to="/add-article" className="btn btn-primary mb-3">Добавить статью</Link>

            <input
                type="text"
                placeholder="Поиск статей по заголовку..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="form-control mb-3"
            />

            <select
                value={selectedTag}
                onChange={(e) => setSelectedTag(e.target.value)}
                className="form-control mb-3"
            >
                {tagOptions.map(tag => (
                    <option key={tag} value={tag}>{tag}</option>
                ))}
            </select>

            {filteredArticles.map(article => (
                <div key={article.id} className="article-container">
                    <div className="article-content">
                        <Article article={article} />
                    </div>
                </div>
            ))}
        </div>
    );
};

export default ShowingArticles;
