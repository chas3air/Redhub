import React, { useEffect, useState } from 'react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from 'recharts';

const ArticlesStats = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const token = localStorage.getItem('token'); // Получаем токен из локального хранилища
                const response = await fetch('http://localhost:80/api/v1/stats/articles', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`, // Добавляем токен в заголовок
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) {
                    throw new Error('Ошибка при получении данных');
                }

                const result = await response.json();
                const formattedData = result.owner_articles.map(owner => ({
                    name: owner.owner_id,
                    count: owner.count_of_articles,
                }));
                setData(formattedData);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    if (loading) {
        return <div>Загрузка данных...</div>;
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    return (
        <div style={{ width: '100%', height: 400 }}>
            <h2>Статистика статей по владельцам</h2>
            <ResponsiveContainer>
                <BarChart data={data}>
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Bar dataKey="count" fill="#8884d8" />
                </BarChart>
            </ResponsiveContainer>
        </div>
    );
};

export default ArticlesStats;