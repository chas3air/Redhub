import React from 'react';
import { Navigate } from 'react-router-dom';

export const RedirectIfAuthenticated = ({ isAuthenticated, element }) => {
    if (isAuthenticated) {
        return <Navigate to="/articles" />;
    }
    return element;
};

const PrivateRoute = ({ element, isAuthenticated, userRole, requiredRole }) => {
    console.log(`Authenticated: ${isAuthenticated}, User Role: ${userRole}, Required Role: ${requiredRole}`);

    if (!isAuthenticated) {
        alert('Вы не авторизованы! Перенаправление на страницу входа.');
        return <Navigate to="/login" />;
    }

    if (requiredRole && userRole?.trim().toLowerCase() !== requiredRole?.trim().toLowerCase()) {
        alert(`Недостаточно прав! Ваша роль: ${userRole}. Требуется: ${requiredRole}`);
        return <Navigate to="/" />;
    }

    return element;
};

export default PrivateRoute;