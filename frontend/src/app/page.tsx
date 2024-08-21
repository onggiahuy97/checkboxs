"use client";

import React, { useEffect, useState } from 'react';

export default function Home() {
    const [grid, setGrid] = useState<boolean[][]>(Array.from({ length: 15 }, () => Array(15).fill(false)));
    const wsRef = React.useRef<WebSocket | null>(null);

    const fetchCurrentState = () => {
       // Fetch the current state of the grid from the server with a GET request to /grid
        fetch('http://localhost:8080/grid')
        .then((response) => response.json())
        .then((data) => {
            setGrid(data);
        })
    }

    const connectWebSocket = () => {
        const ws = new WebSocket('ws://localhost:8080/ws');

        ws.onopen = () => {
            fetchCurrentState()
            console.log('websocket connected');
        };

        ws.onmessage = (event) => {
            if (event.data === 'updated') {
                fetch('http://localhost:8080/grid')
                .then((response) => response.json())
                .then((data) => {
                    setGrid(data);
                })
            }
        };

        ws.onclose = () => {
            console.log('WebSocket closed. Reconnect will be attempted in 1 second.')
            // Reconect after 1 second
            setTimeout(() => {
                connectWebSocket();
            }, 1000);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            ws.close();
        };

        wsRef.current = ws;
    };

    useEffect(() => {
        connectWebSocket();

        // Clean up WebSocket connection on component unmount
        return () => {
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, []);

    const handleCheckboxChange = (row: number, col: number) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            // Send the row and column index to the server
            wsRef.current.send(JSON.stringify({ row, col }));
        } else {
            console.log('WebSocket is not open. Cannot send message.');
        }
    };

    return (
        <div className="flex justify-center items-center  h-screen">
            <div
                className="grid gap-2"
                style={{ gridTemplateColumns: `repeat(15, minmax(0, 1fr))` }}
            >
                {grid.map((row, rowIndex) =>
                    row.map((checked, colIndex) => (
                        <div key={`${rowIndex}-${colIndex}`} className="flex items-center">
                            <input
                                type="checkbox"
                                className="h-8 w-8 mr-2"
                                checked={checked}
                                onChange={() => handleCheckboxChange(rowIndex, colIndex)}
                            />
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}

