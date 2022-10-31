import logo from './logo.svg';
import './App.css';
import { useEffect, useState } from 'react';

function Albums() {
  const [error, setError] = useState(null);
  const [albums, setAlbums] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/albums")
    .then(res => res.json())
    .then(
      (result) => {
        const titles = result.map(album => <li key={album.id}>{album.title}</li>);
        setAlbums(titles);
        setIsLoaded(true);
      },
      (error) => {
        setAlbums(error);
        setIsLoaded(true);
      }
    )
  }, [])

  // useEffect(() => {
  //   console.log(albums)
  // }, [albums])

  if (error) {
    return <div>Error: {error.message} </div>
  } else if (!isLoaded) {
    return <div>Loading ...</div>
  } else {
    return (
      <div>
        {albums}
      </div>
    )
  }
}

function App() {
  return (
    <div className="App">
      <Albums />
    </div>
  );
}

export default App;
