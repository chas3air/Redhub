import React, { useEffect, useState } from 'react';
import {
    PieChart,
    Pie,
    Cell,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from 'recharts';
import './ArticlesStats.css';

const ArticlesStats = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const COLORS = ["#FF0000", "#000000", "#FFFFFF"]; // Красный, черный, белый

    useEffect(() => {
        const fetchData = async () => {
            try {
                const token = localStorage.getItem('token');
                const response = await fetch('http://localhost:80/api/v1/stats/articles', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) {
                    throw new Error('Ошибка при получении данных');
                }

                const result = await response.json();
                const formattedData = result.owner_articles.map((owner, index) => ({
                    name: `Пользователь ${owner.owner_id.substring(0, 6)}`,
                    count: owner.count_of_articles,
                    color: COLORS[index % COLORS.length],
                }));

                const sortedData = [...formattedData].sort((a, b) => b.count - a.count);
                setData(sortedData);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    if (loading) {
        return <div className="loading">Загрузка данных...</div>;
    }

    if (error) {
        return <div className="error">Ошибка: {error}</div>;
    }

    return (
        <div className="stats-container">
            <h2>📊 Статистика статей по владельцам</h2>
            <ResponsiveContainer width="100%" height={400}>
                <PieChart>
                    <Pie
                        data={data}
                        dataKey="count"
                        nameKey="name"
                        cx="50%"
                        cy="50%"
                        outerRadius={150}
                        label
                    >
                        {data.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={entry.color} />
                        ))}
                    </Pie>
                    <Tooltip />
                    <Legend />
                </PieChart>
            </ResponsiveContainer>

            <h3>🏆 Рейтинг владельцев статей</h3>
            <table className="ranking-table">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Владелец</th>
                        <th>Количество статей</th>
                    </tr>
                </thead>
                <tbody>
                    {data.map((owner, index) => (
                        <tr key={index}>
                            <td>{index + 1}</td>
                            <td>{owner.name}</td>
                            <td>{owner.count}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default ArticlesStats;
