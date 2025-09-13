import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import App from './App.tsx'
import ErrorBoundary from '@/components/error-boundary'
import './index.css'
import 'dayjs/locale/zh-cn'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ErrorBoundary>
      <BrowserRouter>
        <ConfigProvider 
          locale={zhCN}
          theme={{
            token: {
              colorPrimary: '#1890ff',
              borderRadius: 6,
            },
            components: {
              Layout: {
                bodyBg: '#f5f5f5',
                headerBg: '#fff',
                siderBg: '#001529',
              },
              Menu: {
                darkItemBg: '#001529',
                darkSubMenuItemBg: '#000c17',
              },
            }
          }}
        >
          <App />
        </ConfigProvider>
      </BrowserRouter>
    </ErrorBoundary>
  </React.StrictMode>,
)