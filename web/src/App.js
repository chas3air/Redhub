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
import ListUsers from './components/ListUsers/ListUsers';
import ShowingArticles from './components/ShowingArticles/ShowingArticles';
import SingleArticle from './components/SingleArticle/SingleArticle';
import AddArticle from './components/AddArticle/AddArticle';
import PrivateRoute, { RedirectIfAuthenticated } from './components/PrivateRoute';
import ListArticles from './components/ListArticles/ListArticles';
import ArticlesStats from './components/Stats/ArticleStats/ArticlesStats';
import UsersStats from './components/Stats/UsersStats/UsersStats';
import AgreementArticles from './components/ArticlesForAgreement/ArticleForAgreement';
import Favorites from './components/Favorites/Favorites';

export default function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [userRole, setUserRole] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            const claims = JSON.parse(atob(token.split('.')[1]));
            setIsAuthenticated(true);
            setUserRole(claims.role);
        }
        setLoading(false);
    }, []);

    const handleLoginSuccess = (token) => {
        setIsAuthenticated(true);
        const claims = JSON.parse(atob(token.split('.')[1]));
        setUserRole(claims.role);
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        setIsAuthenticated(false);
        setUserRole(null);
    };

    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <Router>
            <Fragment>
                <Header title={"RedHub"} userRole={userRole} onLogout={handleLogout} />
                <main>
                    <Routes>
                        <Route 
                            path="/register" 
                            element={<RedirectIfAuthenticated isAuthenticated={isAuthenticated} element={<RegisterForm />} />} 
                        />
                        <Route 
                            path="/login" 
                            element={<RedirectIfAuthenticated isAuthenticated={isAuthenticated} element={<LoginForm onLoginSuccess={handleLoginSuccess} />} />} 
                        />
                        <Route path="/profile" element={<Profile />} />
                        <Route path="/articles" element={<ShowingArticles />} />
                        <Route 
                            path="/add-article" 
                            element={<PrivateRoute element={<AddArticle />} isAuthenticated={isAuthenticated} />} 
                        />
                        <Route
                            path="/users" 
                            element={<PrivateRoute element={<ListUsers />} isAuthenticated={isAuthenticated} userRole={userRole} requiredRole="user_admin" />} 
                        />
                        <Route
                            path="/articles-edited" 
                            element={<PrivateRoute element={<ListArticles />} isAuthenticated={isAuthenticated} userRole={userRole} requiredRole="article_admin" />} 
                        />
                        <Route 
                            path="/moderation" 
                            element={<AgreementArticles />} 
                        />
                        <Route path="/favorites" element={<PrivateRoute element={<Favorites />} isAuthenticated={isAuthenticated} />} />
                        <Route path="/articles/:article_id" element={<SingleArticle />} /> 
                        <Route path="/stats/articles" element={<ArticlesStats />} />
                        <Route path="/stats/users" element={<UsersStats />} />
                        <Route path="*" element={<Navigate to="/login" />} />
                    </Routes>
                </main>
                <Footer />
            </Fragment>
        </Router>
    );
}