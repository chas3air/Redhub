import React from 'react';

export default function Profile() {
    const token = localStorage.getItem('token');

    return (
        <div className="profile">
            <h1>Профиль</h1>
            {token ? <p>Ваш токен: {token}</p> : <p>Токен не найден.</p>}
        </div>
    );
}