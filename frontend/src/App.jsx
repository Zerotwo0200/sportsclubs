import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [field, setField] = useState('club_name');
  const [value, setValue] = useState('');
  const [results, setResults] = useState([]);
  const [error, setError] = useState(null);

  const handleSearch = async () => {
    setError(null);
    if (!value.trim()) {
      setError('Введите значение для поиска');
      return;
    }
    if ((field === 'titles' || field === 'avg_player_age') && isNaN(value)) {
      setError('Для титулов или среднего возраста введите число');
      return;
    }

    try {
      console.log(`Searching for field: ${field}, value: ${value}`);
      const response = await axios.get('http://localhost:8080/search', {
        params: { field, value },
      });
      setResults(Array.isArray(response.data) ? response.data : []);
    } catch (error) {
      console.error('Error fetching data:', error);
      setError(error.response?.data?.error || 'Ошибка при выполнении поиска');
      setResults([]);
    }
  };

  return (
    <div className="App">
      <h1>Спортивные клубы</h1>
      <div className="search-form">
        <select value={field} onChange={(e) => setField(e.target.value)}>
          <option value="club_name">Название клуба</option>
          <option value="city_name">Город</option>
          <option value="titles">Титулы</option>
          <option value="avg_player_age">Средний возраст</option>
        </select>
        <input
          type="text"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder="Введите значение для поиска..."
        />
        <button onClick={handleSearch}>Поиск</button>
      </div>
      {error && <div className="error-message">{error}</div>}
      {results.length > 0 ? (
        <table>
          <thead>
            <tr>
              <th>Название клуба</th>
              <th>Город</th>
              <th>Титулы</th>
              <th>Средний возраст</th>
            </tr>
          </thead>
          <tbody>
            {results.map((club, index) => (
              <tr key={index}>
                <td>{club.club_name}</td>
                <td>{club.city_name}</td>
                <td>{club.titles}</td>
                <td>{club.avg_player_age}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <p className="no-results">Результаты не найдены</p>
      )}
    </div>
  );
}

export default App;