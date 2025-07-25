import React, { useEffect, useState } from 'react';
import {
    PieChart,
    Pie,
    Cell,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from 'recharts';
import './UsersStats.css';

const UsersStats = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const COLORS = ["#FF0000", "#000000", "#FFFFFF"]; // Красный, черный, белый

    useEffect(() => {
        const fetchData = async () => {
            try {
                const token = localStorage.getItem('token');
                const response = await fetch('http://localhost:80/api/v1/stats/users', {
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
                const ageGroups = result.array_of_ages;
                const formattedData = [
                    { name: '0-18 лет', count: ageGroups[0], color: COLORS[0] },
                    { name: '19-25 лет', count: ageGroups[1], color: COLORS[1] },
                    { name: '26-45 лет', count: ageGroups[2], color: COLORS[2] },
                    { name: '46+ лет', count: ageGroups[0] },
                ];

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
        return <div className="loading">Загрузка данных...</div>;
    }

    if (error) {
        return <div className="error">Ошибка: {error}</div>;
    }

    return (
        <div className="stats-container">
            <h2>📊 Статистика пользователей по возрастным группам</h2>
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
                            <Cell key={`cell-${index}`} fill={entry.color || COLORS[index % COLORS.length]} />
                        ))}
                    </Pie>
                    <Tooltip />
                    <Legend />
                </PieChart>
            </ResponsiveContainer>

            <h3>🏆 Рейтинг возрастных групп</h3>
            <table className="ranking-table">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Возрастная группа</th>
                        <th>Количество пользователей</th>
                    </tr>
                </thead>
                <tbody>
                    {data.map((group, index) => (
                        <tr key={index}>
                            <td>{index + 1}</td>
                            <td>{group.name}</td>
                            <td>{group.count}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default UsersStats;
