import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Login from './pages/Auth/Login'
import Signup from './pages/Auth/Signup'
import { Toaster } from 'react-hot-toast'
import { MdCheckBox, MdError } from 'react-icons/md'
import ProtectedLayout from './layouts/protectedLayout/ProtectedLayout'
import Chat from './pages/chat/Chat'
import NotFound from './pages/NotFound/NotFound'
import { useEffect } from 'react'
import { useWebSocketStore } from './store/useWebSocketStore'
import { useWebSocketConnection } from './hooks/useWebSocket'

function App() {
  const { connectWebSocket } = useWebSocketConnection();
  const { socket } = useWebSocketStore()
  useEffect(() => {
    // Only connect when the component mounts
    if(!socket) {
      connectWebSocket()
    }
    // Cleanup WebSocket on unmount
    return () => {
      if (socket) {
        socket.close();
      }
    };
  }, [connectWebSocket, socket]);
  
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<ProtectedLayout/>}>
          <Route path='/' element={<Chat/>}/>
        </Route>
          <Route path='/login' element={<Login/>}/>
          <Route path='/signup' element={<Signup/>}/>
          <Route path='*' element={<NotFound/>}/>
      </Routes>
      <Toaster
      position="bottom-left"
      toastOptions={{
        className: 'toaster',
        duration: 5000,
        success: {
          duration: 2500,
          className: 'toaster success-toast',
          icon: <MdCheckBox />
        },
        error: {
          duration: 5000,
          className: 'toaster error-toast',
          icon: <MdError />
        },
      }}
    />
    </BrowserRouter>
  )
}

export default App