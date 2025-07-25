import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import logo from './../../../assets/redhub_logo.png';

const LoginForm = ({ onLoginSuccess }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        const data = {
            email: email,
            password: password,
        };

        try {
            const response = await fetch('http://localhost/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });

            if (!response.ok) {
                throw new Error('Ошибка при входе');
            }

            const token = await response.text();
            localStorage.setItem('token', token);
            
            onLoginSuccess(token);
            navigate('/profile');
        } catch (error) {
            setErrorMessage('Что-то пошло не так'); // Устанавливаем сообщение об ошибке
            console.error('Ошибка:', error);
        }
    };

    const handleRegisterRedirect = () => {
        navigate('/register'); // Переход на страницу регистрации
    };

    return (
        <div className="d-flex justify-content-center align-items-center" style={{ height: '100vh' }}>
            <div className="w-50 p-4">
                <h2 className="mb-4">Вход</h2>
                {errorMessage && (
                    <div className="alert alert-danger" role="alert">
                        {errorMessage}
                    </div>
                )}
                <form onSubmit={handleSubmit} className="border p-4 rounded bg-light">
                    <div className="form-group">
                        <label htmlFor="email">Email:</label>
                        <input
                            type="email"
                            id="email"
                            className="form-control"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="password">Пароль:</label>
                        <input
                            type="password"
                            id="password"
                            className="form-control"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </div>
                    <button type="submit" className="btn btn-primary">Войти</button>
                </form>
                <button 
                    onClick={handleRegisterRedirect} 
                    className="btn btn-link mt-3"
                >
                    Нет аккаунта? Зарегистрироваться
                </button>
            </div>
            <div className="w-50">
                <img src={logo} alt="Logo" className="img-fluid" style={{ height: '100%', objectFit: 'cover' }} />
            </div>
        </div>
    );
};

export default LoginForm;