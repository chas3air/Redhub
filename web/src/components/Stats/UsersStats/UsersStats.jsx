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

const UsersStats = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const token = localStorage.getItem('token'); // Получаем токен из локального хранилища
                const response = await fetch('http://localhost:80/api/v1/stats/users', {
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
                const ageGroups = result.array_of_ages; // Предполагается, что это [0-18, 19-25, 26-45, 46+]
                const formattedData = [
                    { name: '0-18', count: ageGroups[0] },
                    { name: '19-25', count: ageGroups[1] },
                    { name: '26-45', count: ageGroups[2] },
                    { name: '46+', count: ageGroups[3] },
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
        return <div>Загрузка данных...</div>;
    }

    if (error) {
        return <div>Ошибка: {error}</div>;
    }

    return (
        <div style={{ width: '100%', height: 400 }}>
            <h2>Статистика пользователей по возрастным группам</h2>
            <ResponsiveContainer>
                <BarChart data={data}>
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Bar dataKey="count" fill="#82ca9d" />
                </BarChart>
            </ResponsiveContainer>
        </div>
    );
};

export default UsersStats;