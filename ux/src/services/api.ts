import axios from 'axios';

const API_BASE_URL = '/api/v1';

// 创建axios实例
const api = axios.create({
    baseURL: API_BASE_URL,
    timeout: 30000, // 增加超时时间到30秒
    headers: {
        'Content-Type': 'application/json',
    },
});

// 请求拦截器
api.interceptors.request.use(
    (config) => {
        // 调试日志
        console.log('发送请求:', config.method?.toUpperCase(), config.url);
        if (config.data) {
            try {
                console.log('请求数据:', typeof config.data === 'string' ? config.data : JSON.stringify(config.data, null, 2));
            } catch (e) {
                console.log('请求数据(无法序列化):', config.data);
            }
        }

        // 可以在这里添加认证信息等
        return config;
    },
    (error) => {
        console.error('请求拦截器错误:', error);
        return Promise.reject(error);
    }
);

// 响应拦截器
api.interceptors.response.use(
    (response) => {
        // 调试日志
        console.log('API响应状态:', response.status, response.statusText);
        try {
            console.log('API响应数据:', JSON.stringify(response.data, null, 2));
        } catch (e) {
            console.log('API响应数据(无法序列化):', response.data);
        }

        // 直接返回响应数据，不做处理
        // 这样在服务层可以根据实际情况处理不同的响应结构
        return response.data;
    },
    (error) => {
        // 处理错误响应
        console.error('API响应错误:', error);

        if (error.response) {
            // 服务器返回错误
            console.error('API错误状态:', error.response.status);
            try {
                console.error('API错误数据:', JSON.stringify(error.response.data, null, 2));
            } catch (e) {
                console.error('API错误数据(无法序列化):', error.response.data);
            }
        } else if (error.request) {
            // 请求发送但没有收到响应
            console.error('网络错误 - 没有收到响应:', error.request);
        } else {
            // 请求设置时发生错误
            console.error('请求错误:', error.message);
        }

        return Promise.reject(error);
    }
);

// 测试API连接
api.get('/ping')
    .then(response => {
        console.log('API连接测试成功:', response);
    })
    .catch(error => {
        console.error('API连接测试失败:', error);
    });

export default api; 