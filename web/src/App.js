import './App.css';
import './index.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import { Fragment, useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Header from './components/Header/Header';
import Footer from './components/Footer/Footer';
import LoginForm from './components/Auth/LoginForm/LoginForm';
import RegisterForm from './components/Auth/RegisterForm/RegisterForm';
import Profile from './components/Profile/Profile';

export default function App() {
    const [activeTab, setActiveTab] = useState('login');
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [userRole, setUserRole] = useState(null);

    useEffect(() => {
        // Проверяем наличие токена в localStorage
        const token = localStorage.getItem('token');
        if (token) {
            const claims = JSON.parse(atob(token.split('.')[1])); // Декодируем токен
            setIsAuthenticated(true);
            setUserRole(claims.role);
        }
    }, []);

    const handleLoginSuccess = (token) => {
        setIsAuthenticated(true);
        const claims = JSON.parse(atob(token.split('.')[1])); // Декодируем токен
        setUserRole(claims.role);
    };

    const handleLogout = () => {
        localStorage.removeItem('token'); // Удаляем токен
        setIsAuthenticated(false); // Обновляем состояние аутентификации
        setUserRole(null); // Сбрасываем роль пользователя
    };

    return (
        <Router>
            <Fragment>
                <Header 
                    title={"RedHub"} 
                    userRole={userRole}
                    onLogout={handleLogout}
                />
                <main>
                    {/* Условный рендеринг заголовка */}
                    {!isAuthenticated && <h1>Авторизация и Регистрация</h1>}
                    <Routes>
                        <Route path="/" element={
                            isAuthenticated ? (
                                <Navigate to="/profile" />
                            ) : activeTab === 'login' ? (
                                <LoginForm onLoginSuccess={handleLoginSuccess} />
                            ) : (
                                <RegisterForm />
                            )
                        } />
                        <Route path="/register" element={
                            isAuthenticated ? <Navigate to="/profile" /> : <RegisterForm />
                        } />
                        <Route path="/profile" element={isAuthenticated ? <Profile /> : <Navigate to="/" />} />
                        <Route path="*" element={<Navigate to="/" />} />
                    </Routes>
                </main>
                <Footer />
            </Fragment>
        </Router>
    );
}