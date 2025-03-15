import'./App.css';
import { Fragment } from 'react';
import Header from './components/Header/Header';
import Footer from './components/Footer/Footer';
import TabChoser from './components/TabChoser/Tab';
import Article from './components/Article/Article'
import { useState } from 'react';
import Profile from './components/Profile/Profile';

export default function App() {
  const [tab, setTab]  = useState('articles')

  const article = {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2023-03-05T12:00:00Z",
    "title": "Пример статьи",
    "content": "Это содержимое статьи, которое может быть довольно длинным и подробным.",
    "owner_id": "b771cd98-a256-4c61-9af3-89b57c2d8c61"
  };

  return (
    <Fragment>
        <Header title={"Redhub"} />
        <main>
          <TabChoser active={tab} onChange={(current) => setTab(current)}/>
          
          {tab === "articles" &&
          <>
            <Article article={article} />
            <Article article={article} />
            <Article article={article} />
            <Article article={article} />
          </>
          }
          {tab === "profile" &&
            <Profile />
          }
        </main>
               

        <Footer />
    </Fragment>
  );
}
