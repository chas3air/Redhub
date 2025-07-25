import React, { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import 'bootstrap/dist/css/bootstrap.min.css';

const tagOptions = ['Новости', 'Технология', 'Язык программирования', 'Обзор', 'Реклама'];

const AddArticle = () => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [tag, setTag] = useState(tagOptions[0]);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const token = localStorage.getItem('token');
        if (!token) {
            alert('Ошибка: JWT токен отсутствует');
            return;
        }

        let uid;
        try {
            const claims = JSON.parse(atob(token.split('.')[1]));
            uid = claims.uid;
            if (!uid) throw new Error("UID не найден в токене");
        } catch (err) {
            alert(`Ошибка токена: ${err.message}`);
            return;
        }

        const newArticle = {
            id: uuidv4(),
            title,
            content,
            tag,
            created_at: new Date().toISOString(),
            owner_id: uid,
        };

        console.log(newArticle);

        try {
            const response = await fetch('http://localhost:80/api/v1/moderation/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify(newArticle),
            });

            if (!response.ok) {
                throw new Error('Ошибка при добавлении статьи');
            }

            alert('Статья успешно добавлена!');
            setTitle('');
            setContent('');
            setTag(tagOptions[0]);
        } catch (err) {
            console.error(err);
            alert('Ошибка при добавлении статьи');
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
                <div className="mb-3">
                    <label className="form-label">Категория:</label>
                    <select 
                        className="form-control"
                        value={tag}
                        onChange={(e) => setTag(e.target.value)}
                    >
                        {tagOptions.map(option => (
                            <option key={option} value={option}>{option}</option>
                        ))}
                    </select>
                </div>
                <button type="submit" className="btn btn-primary">Сохранить статью</button>
            </form>
        </div>
    );
};

export default AddArticle;
