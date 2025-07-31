import { useCallback } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { createClient } from '@connectrpc/connect'
import { createConnectTransport } from "@connectrpc/connect-web";
import { OrderService } from './proto/order/v1/order_service_pb'
import { timestampNow } from '@bufbuild/protobuf/wkt'

const transport = createConnectTransport({
  baseUrl: "http://localhost:7444",
});

const client = createClient(OrderService, transport);

function App() {
  const onClick = useCallback(() => {
    client.createOrder({
      id: "ord_012",
      createdAt: timestampNow(),
      items: [
        {
          position: 10,
          productId: 'prod_1',
          quantity: 1n,
        },
      ],
    })
  }, []);

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={onClick}>
          Create Order
        </button>
      </div>
    </>
  )
}

export default App
