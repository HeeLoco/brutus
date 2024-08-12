import { useEffect, useState } from 'react'
import { Toaster } from 'react-hot-toast';
import BrutusNavbar from './features/BrutusNavbar';
import DefaultPage from './pages/DefaultPage/DefaultPage';

function initializeApp() {
  Promise.all([

  ])
}

function App() {
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    initializeApp();
    setInitialized(true);
  }, []);

  return (
    <div className='flex flex-col gap-3 w-screen h-screen'>
      <BrutusNavbar />
      {initialized ? <DefaultPage /> : "Loading..."}
      <Toaster />
    </div>
  );
}

export default App;
