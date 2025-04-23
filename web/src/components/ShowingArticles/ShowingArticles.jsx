import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import Article from '../Article/Article';

const ShowingArticles = () => {
    const [articles, setArticles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedTag, setSelectedTag] = useState(''); // Состояние для выбранного тега

    const tagOptions = ['Все', 'Новости', 'Технология', 'Язык программирования', 'Обзор', 'Реклама'];

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

    if (loading) {
        return <div>Загрузка статей...</div>;
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    const filteredArticles = articles.filter(article =>
        article.title.toLowerCase().includes(searchTerm.toLowerCase()) &&
        (selectedTag === 'Все' || !selectedTag || article.tag?.toLowerCase() === selectedTag.toLowerCase())
    );
      

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
                <Article key={article.id} article={article} />
            ))}
        </div>
    );
};

export default ShowingArticles;
