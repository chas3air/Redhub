import React, { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { useNavigate } from 'react-router-dom';
import logo from '../../../assets/redhub_logo.png';

const RegisterForm = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [nick, setNick] = useState('');
    const [birthday, setBirthday] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        const data = {
            id: uuidv4(),
            email: email,
            password: password,
            nick: nick,
            birthday: new Date(birthday).toISOString(),
        };

        try {
            const response = await fetch('http://localhost/api/v1/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });

            if (!response.ok) {
                const errorText = await response.text();
                setErrorMessage(`Ошибка при регистрации: ${errorText}`);
                throw new Error('Ошибка при регистрации');
            }

            const result = await response.json();

            if (result.success) {
                navigate('/profile');
            } else {
                setErrorMessage(result.message || 'Неизвестная ошибка');
            }
        } catch (error) {
            console.error('Ошибка:', error);
        }
    };

    const handleLoginRedirect = () => {
        navigate('/login');
    };

    return (
        <div className="d-flex justify-content-center align-items-center" style={{ height: '100vh' }}>
            <div className="w-50 p-4">
                <h2 className="mb-4">Регистрация</h2>
                
                {errorMessage && <div className="alert alert-danger">{errorMessage}</div>}

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
                    <div className="form-group">
                        <label htmlFor="nick">Ник:</label>
                        <input
                            type="text"
                            id="nick"
                            className="form-control"
                            value={nick}
                            onChange={(e) => setNick(e.target.value)}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="birthday">Дата рождения:</label>
                        <input
                            type="datetime-local"
                            id="birthday"
                            className="form-control"
                            value={birthday}
                            onChange={(e) => setBirthday(e.target.value)}
                            required
                        />
                    </div>
                    <button type="submit" className="btn btn-primary mt-3">Зарегистрироваться</button>
                </form>
                <button 
                    onClick={handleLoginRedirect} 
                    className="btn btn-link mt-3"
                >
                    У вас уже есть аккаунт? Войти
                </button>
            </div>
            <div className="w-50 d-flex justify-content-center align-items-center">
                <img src={logo} alt="Logo" style={{ width: '400px', height: 'auto' }} />
            </div>
        </div>
    );
};

export default RegisterForm;