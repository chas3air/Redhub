import React, { useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

const AddArticle = () => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();

        const token = localStorage.getItem('token'); // Получаем JWT из localStorage
        const newArticle = {
            id: crypto.randomUUID(), // Генерируем случайный UUID
            title,
            content,
            createdAt: new Date().toISOString(),
            ownerId: crypto.randomUUID(), // Генерируем случайный ownerId (можно заменить)
        };

        try {
            const response = await fetch('http://localhost:80/api/v1/moderation/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`, // Добавляем токен
                },
                body: JSON.stringify(newArticle),
            });

            if (!response.ok) {
                throw new Error('Ошибка при добавлении статьи');
            }

            alert('Статья успешно добавлена!');
            setTitle('');
            setContent('');
        } catch (err) {
            console.error(err);
            alert('Произошла ошибка при добавлении статьи');
        }
    };

    return (
        <div className="container mt-5">
            <h1 className="mb-4">Добавить статью</h1>
            <form onSubmit={handleSubmit}>
                <div className="mb-3">
                    <label className="form-label">Название:</label>
                    <input 
                        type="text" 
                        className="form-control" 
                        value={title} 
                        onChange={(e) => setTitle(e.target.value)} 
                        required 
                    />
                </div>
                <div className="mb-3">
                    <label className="form-label">Содержимое:</label>
                    <textarea 
                        className="form-control" 
                        value={content} 
                        onChange={(e) => setContent(e.target.value)} 
                        required 
                    />
                </div>
                <button type="submit" className="btn btn-primary">Сохранить статью</button>
            </form>
        </div>
    );
};

export default AddArticle;
