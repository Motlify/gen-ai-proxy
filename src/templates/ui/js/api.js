// src/templates/ui/js/api.js

import { showAuthSection } from './auth.js';

// Helper for authenticated fetch requests
export async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem('jwt_token');
    if (token) {
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${token}`,
        };
    }
    const response = await fetch(url, options);
    if (response.status === 401) {
        // If unauthorized, clear token and show login
        localStorage.removeItem('jwt_token');
        if (window.location.pathname === '/dashboard.html') {
            window.location.href = '/'; // Redirect from dashboard to main login
        } else {
            showAuthSection(); // Stay on main page, show auth
        }
        alert('Session expired or unauthorized. Please log in again.');
    }
    return response;
}

// Handle API Key creation
export const handleCreateApiKey = async (event) => {
    event.preventDefault();
    const name = document.getElementById('api_key_name').value;
    const token = localStorage.getItem('jwt_token'); // Use jwt_token

    try {
        const response = await fetch('/api/api-keys', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name })
        });

        if (response.ok) {
            const data = await response.json();
            alert(`API Key created: ${data.api_key}`);
            return { success: true };
        } else {
            const errorData = await response.json();
            alert(`Error creating API Key: ${errorData.error}`);
            return { success: false, error: errorData.error };
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while creating the API Key.');
        return { success: false, error: error.message };
    }
};

// Handle Connection creation
export const handleCreateConnection = async (event) => {
    event.preventDefault();
    const name = document.getElementById('connection_name').value;
    const provider_id = document.getElementById('connection_provider_id').value;
    const api_key = document.getElementById('connection_api_key').value;
    const token = localStorage.getItem('jwt_token');

    try {
        const response = await fetch('/api/connections', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name, provider_id, api_key })
        });

        if (response.ok) {
            alert('Connection created successfully!');
            return { success: true };
        } else {
            const errorData = await response.json();
            alert(`Error creating Connection: ${errorData.error}`);
            return { success: false, error: errorData.error };
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while creating the Connection.');
        return { success: false, error: error.message };
    }
};

// Handle Provider creation
export const handleCreateProvider = async (event) => {
    event.preventDefault();
    const name = document.getElementById('provider_name').value;
    const type = document.getElementById('provider_type').value;
    const base_url = document.getElementById('provider_base_url').value;
    const token = localStorage.getItem('jwt_token');

    try {
        const response = await fetch('/api/providers', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name, type, base_url })
        });

        if (response.ok) {
            alert('Provider created successfully!');
            return { success: true };
        } else {
            const errorData = await response.json();
            alert(`Error creating Provider: ${errorData.error}`);
            return { success: false, error: errorData.error };
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while creating the Provider.');
        return { success: false, error: error.message };
    }
};

// Handle Model creation
export const handleCreateModel = async (event) => {
    event.preventDefault();
    const connection_id = document.getElementById('model_connection_id').value;
    const proxy_model_id = document.getElementById('model_proxy_model_id').value;
    const provider_model_id = document.getElementById('model_provider_model_id').value;
    const price_input = parseFloat(document.getElementById('model_price_input').value);
    const price_output = parseFloat(document.getElementById('model_price_output').value);
    const thinking = document.getElementById('model_thinking').checked;
    const tools_usage = document.getElementById('model_tools_usage').checked;
    const type = document.getElementById('model_type').value;
    const token = localStorage.getItem('jwt_token');

    const modelData = {
        connection_id,
        proxy_model_id,
        provider_model_id,
        price_input: isNaN(price_input) ? 0 : price_input,
        price_output: isNaN(price_output) ? 0 : price_output,
        thinking,
        tools_usage,
        type
    };

    try {
        const response = await fetch(`/api/models`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(modelData)
        });

        if (response.ok) {
            alert('Model created successfully!');
            return { success: true };
        } else {
            const errorData = await response.json();
            alert(`Error creating Model: ${errorData.error}`);
            return { success: false, error: errorData.error };
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while creating the Model.');
        return { success: false, error: error.message };
    }
};

// Generic delete function
export async function handleDelete(id, type) {
    if (confirm(`Are you sure you want to delete this ${type}?`)) {
        try {
            let url = '';
            switch (type) {
                case 'api-key':
                    url = `/api/api-keys/${id}`;
                    break;
                case 'model':
                    url = `/api/models/${id}`;
                    break;
                case 'connection':
                    url = `/api/connections/${id}`;
                    break;
                case 'provider':
                    url = `/api/providers/${id}`;
                    break;
            }
            const response = await authenticatedFetch(url, { method: 'DELETE' });
            if (response.ok) {
                alert(`${type} deleted successfully!`);
                return { success: true };
            } else {
                const errorData = await response.json();
                alert(`Failed to delete ${type}: ${errorData.error || response.statusText}`);
                return { success: false, error: errorData.error || response.statusText };
            }
        } catch (error) {
            console.error(`Error deleting ${type}:`, error);
            alert(`An error occurred while deleting the ${type}.`);
            return { success: false, error: error.message };
        }
    }
    return { success: false, cancelled: true };
}

// Fetch all providers for dropdowns
export async function fetchProvidersForSelect() {
    try {
        const response = await authenticatedFetch('/api/providers');
        if (response.ok) {
            const providers = await response.json();
            return { success: true, data: providers };
        } else {
            const errorData = await response.json();
            return { success: false, error: errorData.error || 'Failed to fetch providers' };
        }
    } catch (error) {
        console.error('Error fetching providers for select:', error);
        return { success: false, error: error.message };
    }
}

// Fetch all connections for dropdowns
export async function fetchConnectionsForSelect() {
    try {
        const response = await authenticatedFetch('/api/connections');
        if (response.ok) {
            const connections = await response.json();
            return { success: true, data: connections.connections };
        } else {
            const errorData = await response.json();
            return { success: false, error: errorData.error || 'Failed to fetch connections' };
        }
    } catch (error) {
        console.error('Error fetching connections for select:', error);
        return { success: false, error: error.message };
    }
}

export async function fetchConversationLogsApi() {
    try {
        const response = await authenticatedFetch('/api/conversation_logs');
        if (response.ok) {
            const logs = await response.json();
            return { success: true, data: logs };
        } else {
            const errorData = await response.json();
            return { success: false, error: errorData.error || 'Failed to fetch conversation logs' };
        }
    } catch (error) {
        console.error('Error fetching conversation logs:', error);
        return { success: false, error: error.message };
    }
}
