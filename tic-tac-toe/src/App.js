import React, { useEffect, useState } from 'react';
import './styles.css';

function Square({ value, onClick }) {
  return (
    <button className="square" onClick={onClick}>
      {value}
    </button>
  );
}

function Board({ room, isXNext }) {
  const [squares, setSquares] = useState(Array(9).fill(null));
  const [gameStatus, setGameStatus] = useState("Waiting for another player...");
  const [socket, setSocket] = useState(null);
  // const [xIsNext, setXIsNext] = useState(isXNext);
  const [won, setWon] = useState(false);
  const [turn, setTurn] = useState(isXNext);

  useEffect(() => {
    if (!room) return;

    const ws = new WebSocket(`wss://omkar.bhaskaraa45.me/tictactoe/ws?room=${room}&join=${!isXNext ? "1" : "0"}`);

    ws.onopen = () => {
      console.log('WebSocket Connected');
      setSocket(ws);
      setGameStatus("Connected");
      ws.send(JSON.stringify({ type: 'ready' }));
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log(message)
      switch (message.type) {
        case 'makeMove':
          // setSquares(message.index);
          markBox(message.index, true)
          // setGameStatus(message.state.status);
          break;
        case 'restart':
          setSquares(Array(9).fill(null));
          setTurn(false);
          setWon(false);
          setGameStatus("Connected");
          break;
        default:
          console.log('Unknown message type:', message.type);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket Disconnected');
      setGameStatus("Disconnected from server");
      setSocket(null);
    };

    ws.onerror = (error) => {
      console.error('WebSocket Error:', error);
      setGameStatus("Error in connection");
    };

    return () => {
      ws.close();
    };
  }, [room]);  // Re-run effect when room changes

  const handleClick = (i) => {
    if (!turn) {
      return
    }
    if (squares[i] != null) {
      return
    }
    console.log("Clicked index:", i);
    if (socket) {
      socket.send(JSON.stringify({ type: 'makeMove', index: i }));
      markBox(i);
    } else {
      console.log("Socket not available or connection not established");
    }
  };

  const markBox = (i, isOpp) => {
    setSquares((prevSquares) => {
      if (calculateWinner(prevSquares) || prevSquares[i]) {
        return prevSquares;
      }
      const newSquares = prevSquares.slice();
      newSquares[i] = isOpp ? (isXNext ? 'O' : 'X') : (isXNext ? 'X' : 'O');
      return newSquares;
    });

    isOpp ? setTurn(true) : setTurn(false)
  };

  const renderSquare = (i) => {
    return <Square value={squares[i]} onClick={() => handleClick(i)} />;
  };

  const winner = calculateWinner(squares);
  const isFull = squares.every(square => square !== null);

  if (winner && !won) {
    setGameStatus(`Winner: ${winner}`)
    setWon(true)
  }

  if (!winner && isFull && !won) {
    setGameStatus(`It's a DRAW`)
    setWon(true)
  }
  // winner ?
  //   setGameStatus(`Winner: ${winner}`) :
  //   setGameStatus(`Next player: ${xIsNext ? 'X' : 'O'}`);

  const restart = () => {
    if (!won) {
      return 
    }
    setSquares(Array(9).fill(null))
    if (socket) {
      socket.send(JSON.stringify({ type: 'restart' }));
      setTurn(true)
      setWon(false)
      setGameStatus("Connected")
    } else {
      console.log("Socket not available or connection not established");
    }
  }

  return (
    <div>
      <div className="status">{gameStatus}</div>
      <div className='turn'>{turn ? "Your turn" : "Opponent's turn"}</div>
      <div className="board-row">
        {renderSquare(0)}
        {renderSquare(1)}
        {renderSquare(2)}
      </div>
      <div className="board-row">
        {renderSquare(3)}
        {renderSquare(4)}
        {renderSquare(5)}
      </div>
      <div className="board-row">
        {renderSquare(6)}
        {renderSquare(7)}
        {renderSquare(8)}
      </div>
      <div>Room ID: {room}</div>
      <div className={won ? "restart" : "restart disable"}>
        <button onClick={restart}>Restart</button>
      </div>
    </div>
  );
}

function calculateWinner(squares) {
  const lines = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ];
  for (let i = 0; i < lines.length; i++) {
    const [a, b, c] = lines[i];
    if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
      return squares[a];
    }
  }
  return null;
}


function App() {
  const [room, setRoom] = useState('');
  const [joined, setJoined] = useState(false);
  const [created, setCreated] = useState(false);

  const handleRoomChange = (event) => {
    setRoom(event.target.value);
  };

  const joinRoom = () => {
    if (room) setJoined(true);
  };

  const createRoom = () => {
    const uniqueRoomId = Math.floor(Math.random() * 10000); // Simple random room number
    setRoom(uniqueRoomId.toString());
    setJoined(true);
    setCreated(true);
  };

  return (
    <div className="game">
      {!joined && (
        <div>
          <div className='create-room'>
            <button onClick={createRoom}>Create Room</button>
          </div>
          <input type="text" value={room} onChange={handleRoomChange} placeholder="Enter room number" />
          <button onClick={joinRoom}>Join Room</button>
        </div>
      )}
      {joined && (
        <div className="game-board">
          <Board room={room} isXNext={created} />

        </div>
      )}
    </div>
  );
}

export default App;
